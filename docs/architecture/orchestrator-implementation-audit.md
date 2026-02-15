# Orchestrator 实施审计报告 (Implementation Audit Report)

> **日期**: 2026-02-15
> **审计对象**: `ai/agents/orchestrator` (DAGScheduler, Executor, ContextInjector)
> **状态**: 🔴 CRITICAL (存在严重并发风险与功能缺失)

---

## 1. 执行摘要 (Executive Summary)

本次审计针对 `Orchestrator` 模块的核心实现进行了静态代码分析与架构一致性审查。审计发现，虽然 **DAG 调度逻辑 (Kahn算法)** 和 **基础执行流** 已跑通，但存在 **Critical 级并发安全隐患** 和 **High 级功能缺失**，目前状态 **不可直接上线**。

核心风险在于 `DAGScheduler` 与 `Executor` 之间共享 `Task` 状态的竞态条件，以及 Context Injection 实现的脆弱性。此外，部分设计违背了 DRY 和 SOLID 原则，增加了维护成本。

---

## 2. 架构一致性审计 (Architecture Compliance)

对比 [Orchestrator Implementation Plan](./orchestrator-implementation-plan.md)，发现以下差异：

| 模块              | 规划要求                           | 当前实现                                   | 差异等级   |
| :---------------- | :--------------------------------- | :----------------------------------------- | :--------- |
| **Resilience**    | 支持 Exponential Backoff 重试      | **完全缺失** (Missing)                     | 🔴 Critical |
| **Input**         | 支持 Context Injection (`{{...}}`) | 实现脆弱 (Regex Replace)，不支持 JSON 转义 | 🟠 High     |
| **Panic**         | Worker 必须 `defer recover`        | ✅ 已在 `DAGScheduler` 中实现               | 🟢 Pass     |
| **Observability** | 结构化日志与 Trace ID              | ✅ 已集成 `slog` 与 `trace_id`              | 🟢 Pass     |
| **Termination**   | 级联跳过 (Cascade Skip)            | ✅ 已实现 BFS 传播逻辑                      | 🟢 Pass     |

**主要发现**:
1.  **重试机制缺位**: 方案明确要求处理 Transient Errors，但目前 Executor 遇到错误直接触发 Handoff 或失败，降低了系统鲁棒性。
2.  **上下文注入风险**: `ContextInjector` 直接字符串替换，若上游结果包含特殊字符（如双引号、换行），将破坏下游 Agent 的 JSON Input 结构，导致解析失败。

---

## 3. 代码质量审计 (Code Quality & Best Practices)

### 3.1 并发安全 (Concurrency Safety) - 🔴 CRITICAL

*   **数据竞争 (Data Race)**:
    *   **位置**: `DAGScheduler.Run()` vs `Executor.executeTaskWithHandoff()`
    *   **现象**: `DAGScheduler` 在主循环中持有 `s.mu` 锁读取 `task.Status`。然而，`Executor` 在 goroutine 中修改 `task.Status`、`task.Result` 时 **未持有任何锁**，也未通过 Channel 通信。
    *   **后果**: 极高概率在多任务并发时导致 Panic (concurrent map read/write) 或状态不一致（调度器读到脏数据导致过早退出或死锁）。

*   **锁粒度混淆**:
    *   **位置**: `ContextInjector.ResolveInput`
    *   **现象**: 使用 `ci.mu.RLock()` 保护传入的局部变量 `tasks` map 的读取。
    *   **问题**: 该锁无法保护 `tasks` map 的内容（Task 结构体字段）。如果 `Execution` 线程正在写入 `task.Result`，此处的读取将发生竞争。`ContextInjector` 不应持有锁，锁应由 `Task` 自身或 `TaskRepository` 管理。

### 3.2 SOLID 原则分析

*   **单一职责原则 (SRP) - 违背**:
    *   `Executor` 承担了过多职责：任务执行、Handoff 策略、事件格式化 (JSON Marshal)、结果收集。
    *   **建议**: 将事件格式化提取到 `EventFormatter`，将 Handoff 决策逻辑彻底剥离。

*   **开闭原则 (OCP) - 违背**:
    *   `ContextInjector` 使用硬编码 Regex (`taskResultRegex`)。若需支持 `{{global.time}}` 或其他变量源，必须修改核心逻辑。
    *   **建议**: 引入 `VariableProvider` 接口，支持可插拔的变量源。

*   **依赖倒置原则 (DIP) - 部分违背**:
    *   `DAGScheduler` 直接依赖具体的Struct `Executor`，而非接口。虽然在内部包中可接受，但这限制了单测 Mock 的能力（目前 `executor_dag_test.go` 可能不得不集成测试）。

### 3.3 DRY (Don't Repeat Yourself) 分析

*   **Event Marshaling 重复**:
    *   `sendPlanEvent`, `sendTaskStartEvent`, `sendTaskEndEvent` 中存在重复的 JSON Marshal 和 Error Handling 代码。
    *   **建议**: 统一封装 `sendEvent(type, payload)` 方法。

*   **Task 状态变更分散**:
    *   Task 的状态变更散落在 `DAGScheduler` 的 panic handler、`Executor` 的成功/失败路径、以及 `handoffHandler` 中。这导致状态流转难以追踪。

---

## 4. 详细问题清单 (Detailed Findings)

### 4.1 严重缺陷 (Critical Defects)

1.  **Race Condition**: `Task` 结构体字段 (`Status`, `Result`, `Error`) 在 `DAGScheduler` (Reader) 和 `Executor` (Writer) 之间缺乏同步机制。
2.  **JSON Injection Vulnerability**: `task.Result` 注入到 JSON 格式的 `task.Input` 时未进行转义。
3.  **Missing Retry Loop**: 网络波动等瞬时错误会导致任务直接失败。

### 4.2 改进建议 (Improvement Areas)

1.  **Context Injection**:
    *   应使用 `text/template` 或自定义 Parser 替代简单的 Regex Replace。
    *   必须识别上下文：若在 JSON 字符串中，应对注入内容进行 `json.Marshal` 转义。

2.  **Handoff Logic**:
    *   目前的递归深度检查 (`depth`) 依赖调用方传递，缺乏全局最大深度强制检查。建议在 `Executor` 入口处强制 `if depth > MaxHandoffDepth { return Error }`。

3.  **Observability**:
    *   目前的日志虽然有 `trace_id`，但缺乏 `span_id` (Task 级 Trace)。建议为每个 Task生成子 Trace ID。

---

## 5. 修复计划建议 (Action Plan)

### 步骤 1: 修复并发安全 (Priority: P0)
*   **方案**:引入 `ThreadSafeTask` 包装器或在 `Task` 结构体中增加 `sync.RWMutex`。所有状态读写必须通过 Setter/Getter 加锁。
*   或者：采用 `Actor` 模型，只有 `DAGScheduler` 主循环有权修改状态，Executor 通过 Channel 发送 `TaskResult` 消息给 Scheduler。

### 步骤 2: 增强 Context Injector (Priority: P1)
*   **方案**: 重构 `ResolveInput`，支持 "Smart Replacement"（检测 JSON 上下文并自动转义）。

### 步骤 3: 实现 Retry 机制 (Priority: P1)
*   **方案**: 在 `Executor.executeTask` 内部增加 `Retrier` 装饰器，实现 Exponential Backoff。

### 步骤 4: 重构 Event 发送 (Priority: P2)
*   **方案**: 提取 `EventFactory`，消除重复的 JSON Marshal 代码。
