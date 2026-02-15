# Orchestrator 架构升级与实施方案

> **日期**: 2026-02-15
> **状态**: Ready for Implementation
> **架构蓝图**: [Orchestrator-Workers Architecture](./orchestrator-workers.md)
> **目标**: 整合 DAG 调度能力、上下文感知与工程韧性，构建支持复杂任务编排的智能 Agent 架构。

---

## 1. 背景与目标

当前 `Orchestrator` 的 `Executor` 组件仅支持简单的并发/串行模式，无法满足复杂任务编排需求。我们需要将其升级为支持 **DAG (有向无环图)** 调度的智能引擎。

**战略对齐 (Strategic Alignment)**:
虽然架构蓝图 (`orchestrator-workers.md`) 规划优先实现 Handoff 机制，但考虑到当前意图执行的确定性需求更为迫切，本方案选择 **优先实现 DAG 调度**。在 Handoff 缺位的情况下，我们将采用 **Fail-Fast (快速失败)** 策略作为过渡方案。

**系统现状 (System Status)**:
- **Decomposer**: 已支持生成含 `dependencies` 字段的任务计划，但下游尚未消费。
- **Executor**: 处于 MVP 阶段，仅支持全局 `Parallel` 开关。
- **Handoff/CapabilityMap**: 尚未实现 (Not Implemented)。

**核心目标**:
1.  **DAG 调度**: 基于 Kahn 算法实现动态任务分发，支持最大化并行。
2.  **上下文注入**: 实现任务间数据流转，支持 `{{task.result}}` 动态替换。
3.  **工程韧性**: 引入 Retry、Exponential Backoff 和 Cascade Skip 机制。
4.  **Agent 增强**: 优化 Schedule/Memo Agent 的交互协议与鲁棒性。

---

## 2. 核心架构设计

### 2.1 DAG 调度引擎 (Dynamic Dispatch via Kahn's Algorithm)

弃用简单的层级调度，采用 **Kahn 算法** 实现基于依赖关系的最大化并行调度。

*   **数据结构**:
    *   `taskMap`: map[string]*Task (任务索引)
    *   `graph`: map[string][]string (邻接表: Upstream -> Downstreams)
    *   `inDegree`: map[string]int (入度表: 记录每个任务的未完成依赖数)
    *   `readyQueue`: Channel (线程安全的就绪任务队列)
    *   `Worker Pool`: 大小为 `MaxParallelTasks` 的 goroutine 池。

*   **调度流程**:
    1.  **初始化**:
        *   计算所有任务 `inDegree`。
        *   将所有 `inDegree == 0` 的任务推入 `readyQueue`。
    2.  **Worker Loop**:
        *   从 `readyQueue` 抢占任务。
        *   **Wait**: 若 Queue 空但 `activeTasks > 0`，则 Pending 等待。
        *   **Execute**: 执行任务业务逻辑。
        *   **Graceful Shutdown**: 显式等待所有活跃 Worker 完成 (`wg.Wait()`)。
    3.  **状态传播**:
        *   **Success**: 查找所有下游任务 `D`，`inDegree[D]--`。若 `inDegree[D] == 0`，将 `D` 推入 `readyQueue`。
        *   **Failure**: 触发 **Retry** (指数退避)；若最终失败，触发 **Cascade Skip**，递归标记所有下游任务为 `Skipped`。

### 2.2 变量替换与上下文注入 (Advanced Context Injection)

在执行任务前处理 `Input`，支持从上游任务或全局上下文提取数据。

*   **语法**: `{{Source.Field | Modifier}}`
    *   `{{task_id.result}}`: 引用上游任务结果。
    *   `{{task_id.result | summary}}`: 引用结果，若超长则自动摘要（需调用 LLM）。
    *   `{{global.time}}`: 引用当前时间上下文。
*   **实现细节**:
    *   **Regex**: 使用 `\{\{\s*([a-zA-Z0-9_\-]+)\.result\s*\}\}` 匹配。
    *   **Token 安全**:
        *   **Trace**: 在注入前后记录 `TokenUsage`。
        *   **Truncation/Summary**: 检测替换后文本长度，若超过阈值（如 2000 tokens），自动触发 Summarizer 压缩。

### 2.3 韧性设计 (Resilience Patterns)

*   **Retry Policy**:
    *   对 `TaskStatusFailed` 进行拦截。
    *   应用 **Exponential Backoff**: initial=1s, multiplier=2, max=10s。
    *   仅重试 "Transient Errors" (网络超时、LLM 5xx)，业务逻辑错误不重试。
*   **Panic Recovery**: Worker 必须 `defer recover()`，防止单任务崩溃导致 Orchestrator 退出。

### 2.4 结构感知聚合 (Structure-Aware Aggregation)

Aggregator 应根据来源 Agent 类型格式化最终输出：
*   **Schedule 来源**: 渲染为 Markdown 表格。
*   **Memo 来源**: 渲染为带引用的列表。

---

## 3. 实施路线图 (Implementation Roadmap)

### 阶段一：快速赢面 (Agent 增强与清理)

本阶段专注于高影响力、低风险的变更，主要涉及配置更新、提示词优化及可观测性建设。

#### 1.1 Schedule Agent 增强 (`config/parrots/schedule.yaml`)

**目标**: 防止过早调用工具，提高对模糊输入的鲁棒性，减少幻觉。

**行动**: 更新 `system_prompt`，加入严格的 "思考-行动" (Thought-Action) 协议和澄清触发器。

**建议变更**:
```yaml
system_prompt: |
  ## Identity
  ... (保留现有 Identity)

  ## Execution Protocol (Strict Order)
  1. <Analyze>: 分析用户意图。时间/时长是否明确？
     - IF 没有时间: 调用 `find_free_time` 或 `ask_user`。
     - IF 修改日程: 必须先调用 `schedule_query` 找到目标事件。
  2. <Validation>: 检查逻辑冲突 (例如：凌晨 3 点开会)。
  3. <Execution>: 调用工具。
  4. <Reflection>: 评估工具输出。
     - IF 冲突: 礼貌地提出替代方案。不要只说 "失败"。

  ## Clarification Triggers (澄清触发器)
  - IF 用户说 "安排会议" (无时间): Ask "请问您希望安排在具体哪天？或者通过 find_free_time 帮您查找合适的时间？"
  - IF 用户说 "和他们开会" (无具体人): Ask "请问是和哪个团队或具体哪位同事？"
```

#### 1.2 Memo Agent 增强 (`config/parrots/memo.yaml`)

**目标**: 区分"列表浏览"与"内容搜索"意图，通过查询扩展提高召回率。

**行动**: 在 `system_prompt` 中加入详细的意图分类和查询扩展指令。

**建议变更**:
```yaml
system_prompt: |
  ## Identity
  ... (保留现有 Identity)

  ## Advanced Search Protocol
  1. **Intent Classification (意图识别)**:
     - **Listing (列表/浏览)**: 用户想要按时间浏览笔记 (e.g., "近期笔记", "今天的记录")。
       - ACTION: 将模糊时间转化为具体的时间词。
       - Example: "近期" -> query="最近7天" (系统底层支持: 今天, 昨天, 本周, 最近N天)。
       - **禁止**: 直接搜索 "近期" 这种模糊关键词。
     - **Searching (搜索)**: 用户想要查找特定内容话题 (e.g., "PostgreSQL 配置")。
       - ACTION: 进入查询扩展流程。

  2. **Query Expansion (查询扩展 - 仅搜索模式)**:
     - 不要只使用用户的原始输入。
     - 生成 2-3 个相关的关键词或同义词。
     - Example: User "DB error" -> Query "DB error database crash exception postgres"

  3. **Answer Synthesis (结果合成)**:
     - IF 找到多条笔记: 先总结共同点。
     - IF 找到具体答案: 直接引用笔记内容回答用户问题。
     - 必须始终使用 `[UID]` 或标题标注来源。
```

#### 1.2.1 Memo Agent 能力与策略分析（补充背景）

*   **无需拆分**: Memo Agent **不需要**拆分为多个工具 (如 `sql_search`, `vector_search`)。其底层的 `AdaptiveRetriever` (自适应检索器) 已具备智能路由能力。
*   **架构关系**: Orchesrator 负责 L1 级 Global Router (跨域)，Memo Retriever 负责 L2 级 Local Router (领域内深搜)。

#### 1.3 可观测性与清理
*   **结构化日志**: 确保 `Decomposer`, `Executor` 输出带 `trace_id` 的 JSON 日志。
*   **清理**: 移除 `scheduler_v2.go` 等废弃代码。

### 阶段二：Orchestrator 核心重构 (DAG & Executor)

#### 2.0 任务 0: Decomposer 专家感知 (前置依赖)
*   **背景**: 由于 CapabilityMap 尚未实现，Decomposer 无法动态获取专家列表。
*   **行动**: 在 `config/orchestrator/decomposer.yaml` 的 System Prompt 中明确列出当前可用的 Agents (`memo`, `schedule`) 及其能力描述，确保 LLM 生成准确的 `agent_id`。

#### 2.1 任务 1: 实现 DAG 引擎核心 (`ai/agents/orchestrator/dag_scheduler.go`)
*   定义 `DAGScheduler` 结构体，封装 `inDegree` 和 `readyQueue`。
*   实现 `scheduleLoop()`: 只有当任务依赖满足时才分发。
*   实现 `markDone(taskID)` 和 `markFailed(taskID)`: 处理依赖更新和 Skipping 逻辑。

#### 2.2 任务 2: 增强 Context Injector (`ai/agents/orchestrator/context_injector.go`)
*   实现 `ResolveInput` 解析 `{{...}}` 语法。
*   **时间上下文**: 在 `Decomposer.Decompose` 中调用 `universal.BuildTimeContext(time.Now())`，并在 Prompt 中注入 `{{.TimeContext}}` 以支持相对时间解析（如"下周五"）。

#### 2.3 任务 3: 升级 Executor (`ai/agents/orchestrator/executor.go`)
*   集成 `DAGScheduler`。
*   增加 `Retry` loop。
*   处理 `StatusSkipped` 事件上报。

### 阶段三：测试与验收

#### 3.1 单元测试 (`ai/agents/orchestrator/executor_dag_test.go`)
需验证以下 Case：
1.  **线性依赖**: A -> B -> C，验证执行顺序与结果传递。
2.  **菱形依赖**: A -> [B, C] -> D，验证 B/C 并行，D 等待 B/C 完成。
3.  **变量替换**: 验证 `Input: "分析 {{t1.result}}"` 被正确替换。
4.  **错误传播**: 验证 A 失败后，依赖 A 的 B 状态变为 `Skipped`。
5.  **Retry 机制**: Mock 失败 2 次成功 1 次的任务，验证重试逻辑。
6.  **Decomposer DAG 解析**: 验证 JSON 解析能正确读取 `dependencies` 字段。

#### 3.2 手动验收 (Acceptance Criteria)
1.  **复杂日程**: "帮我安排下周二下午的团队同步。" -> 验证 Decomposer 准确识别 "下周二"，Executor 正确执行。
2.  **跨域协作**: "找到昨天关于 DB Bug 的笔记，并安排明天的评审会。" -> 验证 Task 2 (Schedule) 依赖 Task 1 (Memo) 的结果。
3.  **模糊指令**: "安排个会。" -> 验证 Agent 主动询问细节。

---

## 4. 风险与注意事项

1.  **任务 ID 稳定性**: LLM 生成的 ID (如 `task_1`) 必须在 `dependencies` 中严格匹配。需加强 `parseTaskPlan` 的校验。
2.  **循环依赖**: 初始化 DAG 时需进行拓扑排序检测环，如有环直接报错。
3.  **Token 爆炸**: 上游结果过长时，需依赖 "Summarizer" 机制或日志告警 (单次注入 > 100KB WARN)。
