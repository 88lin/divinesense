# SPEC-010: 手动验收测试用例

> 优先级: P1 | 阶段: 阶段三 | 状态: 待实现

## 概述

定义 Orchestrator 功能的手动验收测试用例，确保核心功能符合产品要求。

## 验收测试用例

### 场景 1: 复杂日程安排

**用户输入**:
```
"帮我安排下周二下午的团队同步。"
```

**预期行为**:

1. **Decomposer 解析**
   - 正确识别 "下周二" 为具体日期
   - 识别需要调用 schedule agent
   - 生成单一任务: `{ id: "task_1", agent_id: "schedule", input: "..." }`

2. **Executor 执行**
   - 将任务分发给 Schedule Agent
   - Schedule Agent 正确解析 "下周二下午"

3. **Agent 响应**
   - 如果时间明确：直接创建日程
   - 如果时间模糊：触发澄清流程

**验收标准**:
- [ ] Decomposer 正确解析 "下周二"
- [ ] 生成的任务 agent_id 为 "schedule"
- [ ] 日志包含 trace_id

---

### 场景 2: 跨域协作 (Memo + Schedule)

**用户输入**:
```
"找到昨天关于 DB Bug 的笔记，并安排明天的评审会。"
```

**预期行为**:

1. **Decomposer 解析**
   - 识别两个子任务:
     - Task 1: `{ id: "task_1", agent_id: "memo", input: "昨天关于 DB Bug 的笔记" }`
     - Task 2: `{ id: "task_2", agent_id: "schedule", input: "安排明天的评审会", dependencies: ["task_1"] }`

2. **DAG 调度**
   - Task 1 和 Task 2 无依赖关系，先并行执行
   - Task 2 依赖 Task 1，需等待 Task 1 完成

3. **上下文传递**
   - Task 2 执行时，Input 中的 `{{task_1.result}}` 被替换为 Task 1 的结果

**验收标准**:
- [ ] 生成包含依赖关系的任务计划
- [ ] Task 2 在 Task 1 完成后执行
- [ ] Task 2 Input 中正确替换了 Task 1 的结果

---

### 场景 3: 模糊指令处理

**用户输入**:
```
"安排个会。"
```

**预期行为**:

1. **Decomposer 解析**
   - 识别需要 schedule agent
   - 但缺少关键信息：时间、参与人、会议主题

2. **Schedule Agent 行为**
   - 检测到信息不足
   - 触发澄清流程，询问用户:
     - "请问您希望安排在具体哪天？"
     - "请问会议的主题是什么？"

**验收标准**:
- [ ] Agent 主动询问时间
- [ ] Agent 主动询问主题/参与人
- [ ] 不会尝试创建不完整的日程

---

### 场景 4: DAG 并行调度

**用户输入**:
```
"搜索最近的笔记，找出其中关于 AI 的，并安排一个讨论会。"
```

**预期行为**:

1. **Decomposer 解析**
   - Task 1: `{ id: "task_1", agent_id: "memo", input: "搜索最近的笔记" }`
   - Task 2: `{ id: "task_2", agent_id: "memo", input: "从结果中找出关于 AI 的", dependencies: ["task_1"] }`
   - Task 3: `{ id: "task_3", agent_id: "schedule", input: "安排讨论会", dependencies: ["task_2"] }`

2. **DAG 执行**
   - Task 1 先执行
   - Task 1 完成后，Task 2 和 Task 3 都就绪
   - Task 2 和 Task 3 并行执行 (无依赖关系)

**验收标准**:
- [ ] Task 1 -> Task 2 顺序执行
- [ ] Task 2 和 Task 3 并行执行
- [ ] 最终结果聚合了所有任务输出

---

### 场景 5: 错误恢复与重试

**用户输入**:
```
"帮我搜索一下笔记。"
```

**预期行为**:

1. **Agent 执行失败** (如网络超时)
   - Executor 触发重试 (Exponential Backoff)
   - 最多重试 3 次

2. **重试失败**
   - 任务状态标记为 Failed
   - 下游依赖任务标记为 Skipped
   - 返回友好的错误信息给用户

**验收标准**:
- [ ] 失败任务重试 3 次
- [ ] 重试间隔递增 (1s → 2s → 4s)
- [ ] 下游任务正确标记为 Skipped

---

### 场景 6: Handoff 机制

**用户输入**:
```
"帮我搜索笔记并安排会议。"
```

**预期行为**:

1. **Memo Agent 被错误调用**
   - 用户意图需要日程，但 Decomposer 错误分配给 Memo
   - Memo Agent 识别自己无法处理，触发 cannot_complete 事件

2. **Handoff 触发**
   - Handoff Handler 接收事件
   - 查找具备日程能力的替代 Agent (Schedule)
   - 将任务转交给 Schedule Agent

3. **最终结果**
   - 任务成功完成
   - 用户收到正确的日程安排

**验收标准**:
- [ ] Memo Agent 正确识别自己无法处理
- [ ] Handoff 成功转交给 Schedule Agent
- [ ] 用户最终获得正确的日程

---

## 测试执行指南

### 环境要求

| 环境 | 配置 |
|:-----|:-----|
| 数据库 | PostgreSQL (需要 AI 功能) |
| LLM | SiliconFlow + 智谱 Z.AI GLM |
| 服务 | 本地运行 `make start` |

### 执行步骤

1. 启动服务: `make start`
2. 打开前端: http://localhost:25173
3. 在聊天窗口输入测试用例
4. 观察:
   - 后端日志 (trace_id)
   - 前端响应
   - 数据库状态

### 验收检查清单

- [ ] 所有 6 个场景测试通过
- [ ] 日志中包含 trace_id
- [ ] 无 panic 或未处理错误
- [ ] 响应时间合理 (< 30s)

## 依赖

- 前置: SPEC-001 (DAG 调度), SPEC-003 (韧性), SPEC-005/006 (Agent 增强)
- 后置: 无
