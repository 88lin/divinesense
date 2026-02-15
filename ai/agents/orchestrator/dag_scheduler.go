package orchestrator

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// DAGScheduler handles dependency-based task execution.
// It implements Kahn's Algorithm for dynamic task dispatching.
type DAGScheduler struct {
	// Task management
	tasks       map[string]*Task
	graph       map[string][]string // upstream -> downstreams
	inDegree    map[string]int      // task -> remaining dependencies count
	taskIndices map[string]int      // taskID -> original index

	// Execution state
	readyQueue chan string // Tasks ready to run (inDegree = 0)

	// Synchronization
	mu            sync.Mutex
	activeWorkers int

	// Component access
	executor   *Executor
	injector   *ContextInjector
	dispatcher *EventDispatcher

	// Observability
	traceID string
}

func NewDAGScheduler(executor *Executor, tasks []*Task, traceID string, dispatcher *EventDispatcher) (*DAGScheduler, error) {
	s := &DAGScheduler{
		tasks:       make(map[string]*Task),
		graph:       make(map[string][]string),
		inDegree:    make(map[string]int),
		taskIndices: make(map[string]int),
		readyQueue:  make(chan string, len(tasks)), // Buffer large enough for all tasks
		executor:    executor,
		injector:    NewContextInjector(),
		traceID:     traceID,
		dispatcher:  dispatcher,
	}

	// Initialize state
	for i, task := range tasks {
		s.tasks[task.ID] = task
		s.taskIndices[task.ID] = i
		s.inDegree[task.ID] = 0 // Init to 0
	}

	// Build dependency graph
	for _, task := range tasks {
		for _, depID := range task.Dependencies {
			if _, exists := s.tasks[depID]; !exists {
				return nil, fmt.Errorf("task %s depends on unknown task %s", task.ID, depID)
			}

			s.graph[depID] = append(s.graph[depID], task.ID)
			s.inDegree[task.ID]++
		}
	}

	// Seed ready queue
	for _, task := range tasks {
		if s.inDegree[task.ID] == 0 {
			s.readyQueue <- task.ID
		}
	}

	return s, nil
}

// Run starts the scheduling loop.
func (s *DAGScheduler) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	// For MVP, we use a dispatcher loop that feeds a semaphore-controlled worker pool.
	sem := make(chan struct{}, s.executor.config.MaxParallelTasks)

	// Log DAG schedule start
	slog.Info("executor: dag schedule",
		"trace_id", s.traceID,
		"ready_tasks", len(s.tasks),
		"active_tasks", 0,
	)

	// Orchestration loop
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			return err
		case taskID := <-s.readyQueue:
			s.mu.Lock()
			s.activeWorkers++
			s.mu.Unlock()

			wg.Add(1)

			go func(tid string) {
				defer wg.Done()
				sem <- struct{}{}        // Acquire token (blocks if full)
				defer func() { <-sem }() // Release token

				// Wrap execution in a function to handle panics and state updates safely
				func() {
					defer func() {
						// Handle panic
						if r := recover(); r != nil {
							slog.Error("executor: panic in task execution",
								"trace_id", s.traceID,
								"task_id", tid,
								"panic", r)

							// Mark task as failed and cascade
							s.mu.Lock()
							task := s.tasks[tid]
							if !task.Status.IsTerminal() {
								task.Status = TaskStatusFailed
								task.Error = fmt.Sprintf("Panic: %v", r)
								s.cascadeSkip(tid)
							}
							s.mu.Unlock()
						}

						// Always decrement active workers
						s.mu.Lock()
						s.activeWorkers--

						// Log DAG schedule state after task completion
						readyCount := len(s.readyQueue)
						slog.Info("executor: dag schedule",
							"trace_id", s.traceID,
							"ready_tasks", readyCount,
							"active_tasks", s.activeWorkers,
						)
						s.mu.Unlock()
					}()

					// Execute task
					err := s.executeTask(ctx, tid)
					if err != nil {
						// Log error but don't stop the world unless strict mode
						// This allows partial failure in DAG execution
						slog.Warn("DAG task execution failed",
							"trace_id", s.traceID,
							"task_id", tid,
							"error", err.Error())
					}
				}()
			}(taskID)

		default:
			// Check for completion or deadlock
			s.mu.Lock()
			active := s.activeWorkers
			completed := 0
			for _, t := range s.tasks {
				if t.Status.IsTerminal() {
					completed++
				}
			}
			s.mu.Unlock()

			if completed == len(s.tasks) {
				wg.Wait() // Ensure all workers fully exit
				return nil
			}

			if active == 0 {
				// No active workers and ready queue is empty (implied by default case)
				// This means we are stuck -> Cycle detected
				return fmt.Errorf("cycle detected or deadlock: %d/%d tasks completed", completed, len(s.tasks))
			}

			// Avoid hot loop
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// executeTask handles the actual execution logic
func (s *DAGScheduler) executeTask(ctx context.Context, taskID string) error {
	task := s.tasks[taskID]

	// 1. Context Injection
	resolvedInput, err := s.injector.ResolveInput(task.Input, s.tasks)
	if err != nil {
		// Injection failed (e.g. ref not found)
		s.mu.Lock()
		task.Status = TaskStatusFailed
		task.Error = fmt.Sprintf("Context Injection Error: %v", err)
		s.cascadeSkip(taskID)
		s.mu.Unlock()
		return err
	}
	task.Input = resolvedInput

	// 2. Execute
	// Use the event dispatcher for safe event emission
	taskIndex := s.taskIndices[taskID]
	err = s.executor.executeTask(ctx, task, taskIndex, s.dispatcher, s.traceID)

	s.mu.Lock()
	defer s.mu.Unlock()

	if err != nil {
		// Task Failed Logic (Executor already set Status/Error)
		// Cascade Skip
		s.cascadeSkip(taskID)
	} else {
		// Task Success Logic (Executor already set Status/Result)

		// Unblock downstream
		for _, downstreamID := range s.graph[taskID] {
			s.inDegree[downstreamID]--
			if s.inDegree[downstreamID] == 0 {
				s.readyQueue <- downstreamID
			}
		}
	}

	return nil
}

func (s *DAGScheduler) cascadeSkip(failedTaskID string) {
	// BFS or DFS to mark all reachable nodes as Skipped
	queue := []string{failedTaskID}
	visited := make(map[string]bool)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if visited[curr] {
			continue
		}
		visited[curr] = true

		for _, downstream := range s.graph[curr] {
			downTask := s.tasks[downstream]
			if downTask.Status == TaskStatusPending {
				downTask.Status = TaskStatusSkipped
				downTask.Error = fmt.Sprintf("Skipped due to upstream failure in %s", curr)
				queue = append(queue, downstream)

				// Also treat as "Done" for scheduler accounting logic?
				// Yes, otherwise we deadlock.
			}
		}
	}
}
