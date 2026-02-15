# 端到端分析：Orchestrator DAG 实施方案

## 1. 场景模拟 (Scenario Simulation)
**用户请求**: "总结过去 2 小时的邮件，并检查我的日程是否有冲突。"

**任务拆解 (假设)**:
- **T1** (`email`): "列出 `now() - 2h` 期间收到的邮件"
- **T2** (`schedule`): "列出今天的日历事件"
- **T3** (`memo`): "基于 {{T1.result}} 和 {{T2.result}} 进行总结"
  - *依赖*: [T1, T2]

## 2. 执行追踪 (Execution Tracing)

### 步骤 1: 初始化 (Initialization)
- `Executor` 接收包含 3 个任务的 `TaskPlan`。
- `DAGScheduler` 构建图:
  - `readyQueue` (就绪队列): [T1, T2]
  - `inDegree` (入度): {T1:0, T2:0, T3:2}

### 步骤 2: 并行执行 (T1, T2)
- Worker 分别领取 T1 和 T2。
- **T1 执行**:
  - Agent 返回邮件列表 JSON (例如 50 封邮件)。
  - **结果大小风险**: 50 封邮件可能达到 50KB 文本。`resultCollector` 允许通过 (上限 10MB)。
  - T1 完成。`Result` 已存储。
  - T1 通知 T3 的部分依赖完成。`inDegree[T3]` -> 1。
- **T2 执行**:
  - Agent 返回 "今日无日程"。
  - T2 完成。
  - T2 通知 T3。`inDegree[T3]` -> 0。T3 被推入 `readyQueue`。

### 步骤 3: 依赖执行 (T3)
- Worker 领取 T3。
- **上下文注入 (Context Injection)**:
  - `ResolveInput` 读取 `T3.Input`。
  - 将 `{{T1.result}}` 替换为 50KB 的邮件 JSON。
  - 将 `{{T2.result}}` 替换为 "今日无日程"。
  - **上下文窗口风险**: T3 的 Input 现在约 50KB。
    - 如果 `memo` Agent 使用小模型 (如 8k context)，可能会截断或失败。
    - *缓解措施 (未来)*: 引入 Token 计数器或中间总结步骤。
- **执行**:
  - `memo` Agent 处理该长 Context。
  - 返回总结结果。

### 步骤 4: 完成 (Completion)
- 所有任务为终态。
- `Executor` 聚合结果并返回。

## 3. 风险评估 (Risk Assessment)

### 🚨 关键风险 (Pre-Production)

1.  **"软失败" 传播 (Soft Failure Propagation)**
    - **场景**: T1 优雅失败 (LLM 回复 "我无法访问邮件")。
    - **行为**: T1 状态为 `Completed`，Result 是 "我无法访问邮件"。
    - **影响**: T3 将此文本作为数据输入。总结变成了 "基于 '我无法访问邮件'..."
    - **分析**: 这是 "Garbage In, Garbage Out"。在没有 Handoff/结构化错误码的情况下，NLP Agent 难以避免此类问题。
    - **结论**: MVP 阶段可接受，但需优化 Decomposer 的 Prompt 以检查数据有效性。

2.  **上下文爆炸 (Context Explosion)**
    - **场景**: 上游任务返回海量数据 (例如 `search` 返回完整 HTML)。
    - **影响**: 下游注入导致 Context Window 溢出。
    - **结论**: 需监控 `TokenUsage`。未来需引入 "Context Truncation" (上下文截断) 特性。

3.  **变量语法僵化 (Variable Syntax Rigidity)**
    - **场景**: Decomposer 生成 `{{T1.output}}` 而非 `{{T1.result}}`。
    - **影响**: 注入失败 (正则不匹配)。Input 保留原始占位符。
    - **结果**: 下游 Agent 看到的是占位符文本。
    - **结论**: Decomposer Prompt 必须严格规定使用 `{{task_id.result}}` 语法。

## 4. 代码审计 (Code Auditing)

- **Map 访问**: `tasks` Map 在 `Run` 期间结构不可变。`task.Result` 为指针值，是可变的。
- **内存屏障**: `readyQueue` 通道的发送/接收操作构成了 happen-before 关系。T1 写入 Result -> T1 完成 -> 发送 T3 就绪 -> T3 启动 -> T3 读取 T1 Result。**线程安全**。
- **忙等待**: 已通过 `time.Sleep` 和 `activeWorkers` 计数器修复。**安全**。

## 5. 结论 (Conclusion)
即使在复杂的并发场景下，该实施方案在**逻辑上是完备的**且**线程安全**。主要风险来自于 **LLM 的行为不确定性** (上下文限制、幻觉、软失败)，而非 DAG 调度逻辑本身。

**建议**: 推进到集成测试阶段。
