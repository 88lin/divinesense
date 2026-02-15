# SPEC-008: Executor 升级

> 优先级: P0 | 阶段: 阶段二 | 状态: 设计中

## 概述

将现有的简单 Executor 升级为集成 DAGScheduler、ContextInjector 和韧性机制的智能执行器。

## 详细设计

### 架构整合

```
┌─────────────────────────────────────────────────────────────┐
│                        Executor                              │
│  ┌─────────────────┐    ┌──────────────────┐                │
│  │   DAGScheduler │    │  ContextInjector │                │
│  │  - Kahn 算法   │    │  - 变量替换      │                │
│  │  - 依赖管理    │    │  - Token 安全   │                │
│  └────────┬────────┘    └────────┬─────────┘                │
│           │                       │                          │
│           └───────────┬───────────┘                          │
│                       ▼                                      │
│              ┌────────────────┐                              │
│              │  Retry Engine │                              │
│              │  - Backoff    │                              │
│              │  - Panic Recov │                              │
│              └────────────────┘                              │
└─────────────────────────────────────────────────────────────┘
```

### 核心执行流程

```go
// Execute 执行任务计划
func (e *Executor) Execute(ctx context.Context, plan *TaskPlan) (*ExecutionResult, error) {
    // 1. 初始化 DAG Scheduler
    scheduler := NewDAGScheduler(plan.Tasks, e.maxParallel)

    // 2. 启动调度循环
    results, err := scheduler.Schedule(ctx)
    if err != nil {
        return nil, err
    }

    // 3. 聚合结果
    return e.aggregate(results), nil
}

// Schedule 调度执行
func (ds *DAGScheduler) Schedule(ctx context.Context) (map[string]*TaskResult, error) {
    results := make(map[string]*TaskResult)

    for {
        select {
        case taskID := <-ds.readyQueue:
            // 执行任务
            result := ds.executeTask(ctx, taskID)
            results[taskID] = result

            // 处理结果传播
            if result.Status == TaskStatusSuccess {
                ds.markDone(taskID)
            } else {
                ds.markFailed(taskID)
            }

        case <-ctx.Done():
            // 优雅退出
            ds.gracefulShutdown()
            return results, ctx.Err()
        }
    }
}

// executeTask 单任务执行 (含重试)
func (ds *DAGScheduler) executeTask(ctx context.Context, taskID string) *TaskResult {
    task := ds.taskMap[taskID]

    // 1. 上下文注入
    input := ds.injector.Resolve(task.Input, ds.taskResults)

    // 2. Retry Loop
    var lastErr error
    for attempt := 0; attempt <= ds.retryConfig.MaxRetries; attempt++ {
        if attempt > 0 {
            delay := ds.retryConfig.NextDelay(attempt - 1)
            time.Sleep(delay)
        }

        // 执行
        result, err := ds.runAgent(ctx, task.AgentID, input)
        if err == nil {
            return result
        }

        if !IsTransientError(err) {
            // 非瞬时错误，不重试
            return &TaskResult{Status: TaskStatusFailed, Error: err.Error()}
        }

        lastErr = err
        ds.logger.Warn("task failed, retrying", "task", taskID, "attempt", attempt)
    }

    // 重试耗尽
    return &TaskResult{Status: TaskStatusFailed, Error: lastErr.Error()}
}
```

### 事件上报

```go
// TaskEvent 任务状态事件
type TaskEvent struct {
    TaskID   string    `json:"task_id"`
    Status   string    `json:"status"`
    Result   string    `json:"result,omitempty"`
    Error    string    `json:"error,omitempty"`
    Duration int64     `json:"duration_ms"`
}

// 事件处理
func (e *Executor) handleTaskEvent(event TaskEvent) {
    switch event.Status {
    case TaskStatusSkipped:
        // 上报 cascade skip 事件
        e.metrics.RecordCascadeSkip(event.TaskID)
    case TaskStatusFailed:
        // 上报失败事件
        e.metrics.RecordFailure(event.TaskID, event.Error)
    }
}
```

## 验收标准

- [ ] Executor 集成 DAGScheduler 实现并行调度
- [ ] 任务输入正确注入上下文变量
- [ ] 失败任务触发重试逻辑
- [ ] 上游失败时下游任务被标记为 Skipped
- [ ] Graceful Shutdown: 取消时等待活跃任务完成

## 实现提示

1. **文件位置**: `ai/agents/orchestrator/executor.go`
2. **依赖**: 需要先实现 SPEC-001, SPEC-002, SPEC-003
3. **测试**: 创建完整流程的集成测试

## 依赖

- 前置: SPEC-001, SPEC-002, SPEC-003
- 后置: 无
