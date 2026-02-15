# SPEC-007: 可观测性建设

> 优先级: P1 | 阶段: 阶段一 | 状态: 待实现

## 概述

为 Orchestrator 引入结构化日志和可观测性支持，确保任务执行链路可追踪。

## 详细设计

### 日志要求

所有 Orchestrator 组件必须输出带 `trace_id` 的 JSON 格式日志。

#### Decomposer 日志

```go
// Decompose 开始
slog.Info("decomposer: start decompose",
    "trace_id", traceID,
    "user_input", userInput,
    "timestamp", time.Now().Unix(),
)

// Decompose 完成
slog.Info("decomposer: decompose complete",
    "trace_id", traceID,
    "task_count", len(plan.Tasks),
    "duration_ms", duration.Milliseconds(),
)
```

#### Executor 日志

```go
// 任务开始执行
slog.Info("executor: task start",
    "trace_id", traceID,
    "task_id", task.ID,
    "agent_id", task.Agent,
    "dependencies", task.DependsOn,
)

// 任务完成
slog.Info("executor: task complete",
    "trace_id", traceID,
    "task_id", task.ID,
    "status", result.Status,
    "duration_ms", duration.Milliseconds(),
)

// 任务失败
slog.Warn("executor: task failed",
    "trace_id", traceID,
    "task_id", task.ID,
    "error", result.Error,
    "retry_count", task.RetryCount,
)

// DAG 调度事件
slog.Info("executor: dag schedule",
    "trace_id", traceID,
    "ready_tasks", len(readyQueue),
    "active_tasks", len(activeTasks),
)
```

#### Handoff 日志

```go
// Handoff 开始
slog.Info("handoff: start",
    "trace_id", traceID,
    "task_id", task.ID,
    "from_agent", task.Agent,
    "to_agent", result.NewExpert,
    "depth", result.Depth,
)

// Handoff 失败
slog.Warn("handoff: failed",
    "trace_id", traceID,
    "task_id", task.ID,
    "reason", result.Reason,
    "fallback_message", result.FallbackMessage,
)
```

### Trace ID 传递

```go
// TaskContext 在任务间传递 trace_id
type TaskContext struct {
    TraceID    string
    UserID     int32
    BlockID    int64
    ParentTaskID string
}

// 从 BlockContext 初始化
func NewTaskContext(blockCtx *BlockContext) *TaskContext {
    return &TaskContext{
        TraceID: generateTraceID(),
        UserID:  blockCtx.UserID,
        BlockID: blockCtx.BlockID,
    }
}
```

### 日志字段规范

| 字段 | 类型 | 说明 | 必填 |
|:-----|:-----|:-----|:-----|
| trace_id | string | 链路追踪 ID | ✅ |
| timestamp | int64 | Unix 时间戳 (毫秒) | ✅ |
| level | string | 日志级别 (debug/info/warn/error) | ✅ |
| component | string | 组件名 (decomposer/executor/handoff) | ✅ |
| task_id | string | 任务 ID | 任务相关时 |
| duration_ms | int64 | 执行耗时 (毫秒) | 任务完成时 |

## 验收标准

- [ ] Decomposer 输出带 trace_id 的 JSON 日志
- [ ] Executor 输出带 trace_id 的 JSON 日志
- [ ] Handoff 输出带 trace_id 的 JSON 日志
- [ ] trace_id 在任务链中传递
- [ ] 日志支持 JSON 格式输出

## 实现提示

1. **文件位置**: `ai/agents/orchestrator/` 现有文件
2. **配置**: 通过 `config/logging.yaml` 控制日志级别
3. **格式**: 使用 `log/slog` 的 JSON handler

## 依赖

- 前置: 无
- 后置: 无
