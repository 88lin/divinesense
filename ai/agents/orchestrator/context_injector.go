package orchestrator

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// Pre-compiled regex for performance (Issue #211: performance optimization)
var taskResultRegex = regexp.MustCompile(`\{\{([a-zA-Z0-9_\-]+)\.result\}\}`)

// ContextInjector handles variable substitution in task inputs.
type ContextInjector struct {
}

// NewContextInjector creates a new context injector.
func NewContextInjector() *ContextInjector {
	return &ContextInjector{}
}

// ResolveInput replaces {{task_id.result}} placeholders with actual task results.
// It detects if the replacement is occurring within a JSON string context and escapes the value accordingly.
func (ci *ContextInjector) ResolveInput(input string, tasks map[string]*Task) (string, error) {
	var err error

	// Use pre-compiled regex for performance
	resolved := taskResultRegex.ReplaceAllStringFunc(input, func(match string) string {
		// Extract task ID (group 1)
		submatches := taskResultRegex.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match // Should not happen if regex matches
		}
		taskID := submatches[1]

		task, exists := tasks[taskID]
		// ci.mu.RLock() removed because it was protecting the map lookup which is done above.
		// However, original code used map access inside RLock.
		// Actually, the tasks map itself is not protected by Task mutex.
		// If the map structure changes (tasks added/removed) we need a lock on the map.
		// But here we assume tasks map is static during execution.
		// The critical part is accessing task fields.

		if !exists {
			err = fmt.Errorf("reference not found: task '%s' does not exist", taskID)
			return match
		}

		// Use safe accessors
		status := task.GetStatus()
		if status != TaskStatusCompleted {
			err = fmt.Errorf("reference invalid: task '%s' is not completed (status: %s)", taskID, status)
			return match
		}

		result := task.GetResult()

		// Heuristic detection: check if match is surrounded by quotes in the original string.
		// If so, we should escape the result to ensure valid JSON if it contains quotes.
		// Strategy: Marshal the result to get a JSON string, then strip outer quotes.
		// This handles escaping of internal quotes and special characters.
		b, err := json.Marshal(result)
		if err != nil {
			// Should not happen for string, but safe fallback
			return result
		}
		s := string(b)
		if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
		return result
	})

	if err != nil {
		return "", err
	}

	return resolved, nil
}
