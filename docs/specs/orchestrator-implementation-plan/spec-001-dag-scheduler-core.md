# SPEC-001: DAG 调度核心 (Kahn 算法)

> 优先级: P0 | 阶段: 阶段二 | 状态: 设计中

## 概述

实现基于 Kahn 算法的 DAG 调度引擎，支持任务依赖关系的最大化并行执行。废弃简单的层级调度，采用基于拓扑排序的动态分发机制。

## 详细设计

### 核心逻辑

**Kahn 算法核心**:
1. 计算所有任务的入度 (`inDegree`)
2. 将入度为 0 的任务加入就绪队列 (`readyQueue`)
3. Worker 从就绪队列抢占任务执行
4. 任务完成后，查找所有下游任务并递减其入度
5. 若下游任务入度变为 0，则推入就绪队列

**关键特性**:
- **最大化并行**: 无依赖任务可同时执行
- **死锁检测**: 初始化时检测环，存在环则报错
- **Graceful Shutdown**: 所有活跃 Worker 完成后才退出

### 数据结构

```go
// DAGScheduler DAG 调度器
type DAGScheduler struct {
    taskMap     map[string]*Task           // 任务索引
    graph       map[string][]string        // 邻接表: Upstream -> Downstreams
    inDegree    map[string]int             // 入度表
    readyQueue  chan string                 // 就绪任务队列
    mu          sync.RWMutex                // 保护共享状态

    // 配置
    MaxParallelTasks int                    // 最大并行任务数
    workerPool      *WorkerPool             // goroutine 池
}

// Task 任务定义
type Task struct {
    ID          string                 `json:"id"`
    AgentID     string                 `json:"agent_id"`
    Input       string                 `json:"input"`
    Status      TaskStatus             `json:"status"`
    Result      string                 `json:"result,omitempty"`
    Error       string                 `json:"error,omitempty"`
    RetryCount  int                    `json:"retry_count"`
    DependsOn   []string               `json:"dependencies,omitempty"` // 依赖的任务 ID
}
```

### 调度流程图

```
┌─────────────────────────────────────────────────────────────┐
│                      DAG 初始化                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ 1. 构建 graph (邻接表)                               │   │
│  │ 2. 计算 inDegree                                    │   │
│  │ 3. 检测循环依赖                                      │   │
│  │ 4. 将 inDegree==0 的任务推入 readyQueue             │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Worker Loop                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ SELECT:                                             │   │
│  │   CASE task := <- readyQueue:                      │   │
│  │     Execute(task)                                   │   │
│  │     → markDone(task):                              │   │
│  │       FOR each downstream:                          │   │
│  │         inDegree[downstream]--                      │   │
│  │         IF inDegree[downstream] == 0:              │   │
│  │           readyQueue <- downstream                 │   │
│  │   CASE <- ctx.Done():                              │   │
│  │     gracefulShutdown()                              │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 失败传播逻辑

```
┌─────────────────────────────────────────────────────────────┐
│                    任务失败处理流程                          │
│                                                             │
│  Task A 失败                                                │
│       │                                                     │
│       ▼                                                     │
│  ┌─────────────┐                                           │
│  │ Retry?      │──Yes──▶ 重试 (Exponential Backoff)        │
│  └──────┬──────┘                                           │
│         │ No                                                │
│         ▼                                                   │
│  ┌─────────────────┐                                        │
│  │ Cascade Skip    │                                        │
│  │ 递归标记所有下游 │                                        │
│  │ 任务为 Skipped  │                                        │
│  └─────────────────┘                                        │
└─────────────────────────────────────────────────────────────┘
```

## 验收标准

- [x] 线性依赖 A->B->C 顺序执行，B 等待 A 完成 (见 `TestDAG_LinearExecution`)
- [x] 菱形依赖 A->[B,C]->D，B 和 C 并行执行，D 等待 B 和 C 都完成 (见 `TestDAG_DiamondExecution`)
- [ ] 多组独立任务 (A,B) 和 (C,D) 完全并行 (需补充测试)
- [x] 存在环时初始化报错 "circular dependency detected" (见 `TestDAG_CircularDependency`)
- [x] 任务失败后下游任务状态变为 Skipped (见 `TestDAG_CascadeSkip`)

> **注意**: 单元测试代码已存在于 `ai/agents/orchestrator/executor_dag_test.go`

## 实现提示

1. **并发安全**: 使用 `sync.RWMutex` 保护 `inDegree` 和 `taskMap`
2. **Worker Pool**: 使用固定大小的 goroutine 池控制并发数
3. **优雅退出**: 使用 `context.Context` 支持取消，Worker 完成后才返回
4. **循环检测**: 初始化时对所有节点进行拓扑排序验证

## 依赖

- 前置: 无
- 后置: SPEC-002, SPEC-003, SPEC-009
