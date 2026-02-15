# SPEC-003: 韧性设计 (Retry/Exponential Backoff)

> 优先级: P0 | 阶段: 阶段二 | 状态: 设计中

## 概述

为 Orchestrator 引入工程韧性机制，包括 Retry 策略、指数退避和级联跳过 (Cascade Skip)，确保复杂任务编排的可靠性。

## 详细设计

### Retry Policy

**重试触发条件**:
- 网络超时
- LLM 返回 5xx 错误
- 临时性服务不可用

**不重试条件**:
- 业务逻辑错误 (如参数验证失败)
- 认证失败 (401/403)
- 客户端错误 (4xx)

### Exponential Backoff 配置

```go
// RetryConfig 重试配置
type RetryConfig struct {
    MaxRetries      int           // 最大重试次数 (默认 3)
    InitialDelay    time.Duration // 初始延迟 (1s)
    MaxDelay        time.Duration // 最大延迟 (10s)
    Multiplier      float64       // 退避倍数 (2.0)
    Jitter          bool          // 是否添加随机抖动
}

// 计算下次延迟
func (rc *RetryConfig) NextDelay(attempt int) time.Duration {
    delay := rc.InitialDelay * time.Duration(math.Pow(rc.Multiplier, float64(attempt)))
    if delay > rc.MaxDelay {
        delay = rc.MaxDelay
    }
    if rc.Jitter {
        delay = delay / 2 + time.Duration(rand.Int63n(int64(delay/2)))
    }
    return delay
}
```

### Panic Recovery

```go
// SafeExecute 安全执行包装
func SafeExecute(fn func() error) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
            // 记录日志，上报 metrics
            logger.Error("task panic", "recover", r)
        }
    }()
    return fn()
}
```

### Cascade Skip 逻辑

```
Task A 失败
    │
    ├──▶ Task B (依赖 A) → 状态 = Skipped, 原因 = "上游任务失败"
    │
    ├──▶ Task C (依赖 A, B) → 状态 = Skipped, 原因 = "上游任务失败"
    │
    └──▶ Task D (独立) → 状态 = Unchanged (继续执行)
```

```go
// cascadeSkip 递归跳过下游任务
func (d *DAGScheduler) cascadeSkip(failedTaskID string) {
    for _, downstream := range d.graph[failedTaskID] {
        task := d.taskMap[downstream]
        if task.Status != TaskStatusSkipped {
            task.Status = TaskStatusSkipped
            task.Error = "上游任务 " + failedTaskID + " 失败"
            // 递归传播
            d.cascadeSkip(downstream)
        }
    }
}
```

## 验收标准

- [ ] 瞬时错误重试 3 次，每次延迟翻倍 (1s → 2s → 4s)
- [ ] 业务错误不重试，直接失败
- [ ] 任务 panic 被捕获，不导致 Orchestrator 崩溃
- [ ] A 失败后，依赖 A 的所有下游任务状态变为 Skipped
- [ ] 独立任务不受上游失败影响

## 实现提示

1. **错误分类**: 定义 `IsTransientError(err error)` 判断是否可重试
2. **Metrics**: 上报重试次数、成功率等指标
3. **日志**: 重试时记录 `warn` 级别日志，包含尝试次数

## 依赖

- 前置: SPEC-001 (DAG 调度)
- 后置: SPEC-009 (Executor 升级)
