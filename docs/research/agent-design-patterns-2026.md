# Agent 设计模式调研报告 (2025-2026)

> **调研时间**: 2026-02-08 | **目的**: DivineSense 三鹦鹉统一架构设计
>
> **调研范围**: Google、Anthropic、OpenAI 等大厂最新 Agent 研究和设计模式

---

## 执行摘要

| 厂商 | 核心模式 | 关键发现 | 论文/发布时间 |
|:-----|:---------|:---------|:-------------|
| **Anthropic** | Orchestrator-Worker | 90.2% 性能提升，8条Prompt工程原则 | 2025-06 |
| **Google** | 四种架构评估 | "更多Agent更好"是迷思，对齐原则 | 2025-12 |
| **OpenAI** | Swarm Handoffs | 轻量级协调框架，已被Agents SDK取代 | 2025 |

**核心结论**：
- ✅ 任务特性决定架构选择（Google对齐原则）
- ✅ 中心化协调器适合需要综合的场景（Anthropic）
- ❌ 顺序任务使用多Agent会降低性能（-27%）
- ✅ Prompt工程比架构选择更重要

---

## 1. Anthropic: Orchestrator-Worker 模式

**来源**: [How we built our multi-agent research system](https://www.anthropic.com/news/multi-agent-research) (2025-06-13)

### 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                  LeadResearcher (Claude Opus 4)             │
│  - 规划任务                                                 │
│  - 分配子任务                                               │
│  - 综合结果                                                 │
└──────────────┬────────────┬────────────┬────────────────────┘
               │            │            │
       ┌───────▼────┐ ┌────▼─────┐ ┌───▼──────┐ ┌─────────┐
       │ Subagent 1 │ │Subagent 2│ │Subagent 3│ │   ...   │
       │(Sonnet 4) │ │(Sonnet 4)│ │(Sonnet 4)│ │         │
       └───────────┘ └──────────┘ └──────────┘ └─────────┘
              并行执行 → 结果汇总
```

### 8条Prompt工程原则

| # | 原则 | 说明 |
|:--|:-----|:-----|
| 1 | **Think like your agents** | 提供子代理的完整上下文，不要假设它们知道什么 |
| 2 | **Teach orchestrator how to delegate** | 明确何时使用哪个子代理 |
| 3 | **Prevent unnecessary conversations** | 子代理间不通信，避免级联延迟 |
| 4 | **Simplify communication** | 标准化输出格式（如JSON） |
| 5 | **Keep it concrete** | 使用具体示例而非抽象描述 |
| 6 | **Take time to plan** | 规划步骤比执行更重要 |
| 7 | **Teach your agents how to fail** | 定义错误处理流程 |
| 8 | **Track your agents** | 监控每个子代理的性能 |

### 性能数据

- **单Agent (Claude Opus 4)**: 基准性能
- **多Agent系统**: 90.2% 性能提升
- **关键因素**: 子代理使用更便宜的模型（Claude Sonnet 4），降低成本

### 代码模式

```python
# Anthropic 的实际实现
lead_researcher = Agent(
    model="claude-opus-4",
    instructions="""你是主编，协调3个子代理并行研究。
    任务分解原则：
    1. 将复杂任务拆分为独立子任务
    2. 分配给最合适的子代理
    3. 综合子代理结果，生成最终答案
    """
)

subagent = Agent(
    model="claude-sonnet-4",
    instructions=f"""你负责具体任务。
    上下文：{full_context}
    输出格式：JSON
    失败处理：如果无法完成，返回 {{'error': 'reason'}}
    """
)
```

---

## 2. Google: 对齐原则 (Alignment Principle)

**来源**: [Towards a Science of Scaling Agent Systems](https://research.google/blog/towards-a-science-of-scaling-agent-systems/) (2025-12)

### 180个配置评估

Google 系统性评估了 180 种不同的 Agent 配置，得出以下定量原则：

### 任务类型与架构匹配

| 任务类型 | 多Agent效果 | 推荐架构 | 性能变化 |
|:---------|:-----------|:---------|:---------|
| **可并行** | +28% 性能 | Independent / Decentralized | ✅ 显著提升 |
| **顺序依赖** | -27% 性能 | Single-Agent | ❌ 显著下降 |
| **混合** | +5% ~ -15% | Hybrid | ⚠️ 需谨慎设计 |

### 四种架构模式

```
1. Single-Agent     ──→ 一个Agent完成所有任务
                        适用: 顺序依赖任务

2. Independent      ──→ 多Agent独立工作，无协调
                        适用: 高度并行任务

3. Centralized      ──→ 中心Orchestrator协调子Agent
                        适用: 需要结果综合

4. Decentralized    ──→ Agent间直接通信
                        适用: 分布式决策

5. Hybrid           ──→ 混合以上模式
                        适用: 复杂场景
```

### 误差放大公式

```
Error_Amplification = 1 - (1 - base_error)^n

示例：
- 当 n=5, base_error=5%:
    Error_Amplification = 1 - 0.95^5 = 22.6%

- 当 n=10, base_error=5%:
    Error_Amplification = 1 - 0.95^10 = 40.1%
```

**启示**: 更多Agent ≠ 更好性能，需权衡误差放大风险。

### 核心结论

| 迷思 | 真相 |
|:-----|:-----|
| "更多Agent = 更好性能" | ❌ 任务特性决定架构 |
| "多Agent总是更快" | ❌ 顺序任务性能下降27% |
| "复杂场景需要多Agent" | ⚠️ 需要评估并行度 |

---

## 3. OpenAI Swarm: Handoff 模式

**来源**: [OpenAI Swarm GitHub](https://github.com/openai/swarm) (2025)

### 核心原语

```python
# Agent 定义
agent_a = Agent(
    name="Agent A",
    instructions="你是专家A，负责...",
    functions=[func1, func2]
)

# Handoff 定义（状态转移）
def transfer_to_agent_b():
    return agent_b

# 执行
response = swarm.run(agent_a, "用户消息")
```

### 设计原则

| 原则 | 说明 |
|:-----|:-----|
| **无状态** | Agent不保留对话历史 |
| **轻量级** | 仅专注协调和切换 |
| **函数驱动** | Handoff通过函数返回实现 |

### 注意事项

> ⚠️ **Swarm 已被 OpenAI Agents SDK 取代**
>
> Swarm 是教育性框架，生产环境推荐使用 [OpenAI Agents SDK](https://platform.openai.com/docs/agents)

---

## 4. DivineSense 三鹦鹉统一架构建议

### 当前状态分析

| Agent | 执行模式 | 工具集 | 特点 |
|:-----|:---------|:------|:-----|
| **MemoParrot** | ReAct循环 | memo_search | 单一检索 |
| **ScheduleParrotV2** | Native工具调用 | schedule_* CRUD | 顺序依赖，薄包装 |
| **AmazingParrot** | 两阶段 (Plan→Retrieve→Synthesize) | memo + schedule | 并行检索两者 |

**问题**:
- ~70% 代码重复
- ScheduleParrotV2 使用完全不同的框架
- 行为模式差异大（ReAct vs Native vs Two-Phase）

### 推荐方案: Centralized Orchestrator + Specialized Subagents

```
┌─────────────────────────────────────────────────────────────┐
│                  UniversalParrot (Base)                     │
│  - 统一的 ReAct 循环引擎                                     │
│  - 统一的 LRU 缓存层                                         │
│  - 统一的 SessionStats 收集                                 │
│  - 统一的 Tool Calling 框架                                 │
└─────────────────────────────────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
   ┌────▼─────┐    ┌─────▼──────┐    ┌────▼─────┐
   │MemoParrot│    │ScheduleParrot│   │AmazingParrot│
   │          │    │    V2       │    │          │
   └────┬─────┘    └─────┬──────┘    └────┬─────┘
        │                │                 │
   memo_search     schedule_*        并行调用两者
                      工具集
```

### 统一后的差异点

| 方面 | 实现方式 | 示例 |
|:-----|:---------|:-----|
| **Prompt** | 每只鹦鹉独立的 system_prompt | `GetMemoSystemPrompt()` vs `GetScheduleSystemPrompt()` |
| **Tool Chain** | 配置化的工具列表 | `[memoSearchTool]` vs `[scheduleAddTool, scheduleQueryTool, ...]` |
| **行为模式** | 参数化的执行策略 | `mode: "sequential"` vs `mode: "parallel"` |

---

## 5. 交互式确认问题

### Q1: 架构模式选择

根据 Google 研究，任务特性决定最佳架构。三鹦鹉场景：
- **MemoParrot**: 单一检索任务（可并行）
- **ScheduleParrotV2**: 顺序依赖的 CRUD 操作
- **AmazingParrot**: 需要并行检索两者

**选项**：
1. **Centralized Orchestrator** — Anthropic 模式，统一调度
2. **Skills Pattern** — 单 Agent 动态加载技能
3. **Hybrid** — Memo/Amazing 并行，Schedule 独立

**推荐**: 选项 1 (Centralized Orchestrator)，因为 Amazing 需要协调两者。

### Q2: ScheduleParrotV2 特殊处理

当前 `ScheduleParrotV2` 是 `SchedulerAgentV2` 的薄包装，使用完全不同的框架。

**选项**：
1. **保持独立框架** — Schedule 任务顺序性强，不适合 ReAct
2. **统一到 ReAct** — 将原生工具调用改为 ReAct 循环
3. **双模式支持** — UniversalParrot 支持 ReAct 和 Native 两种模式

**推荐**: 选项 2，但保留 Native 快速路径（简单 CRUD 直达）。

### Q3: Prompt 管理策略

**选项**：
1. **代码内嵌** — 当前方式，每个 Agent 文件中定义
2. **PromptRegistry 集中** — 已存在，扩展覆盖三鹦鹉
3. **文件外部化** — `prompts/` 目录，便于热更新

**推荐**: 选项 2 (PromptRegistry)，已有基础设施。

### Q4: 工具链差异化管理

**选项**：
1. **硬编码工具集** — 每个 Agent 构造时传入固定工具列表
2. **能力声明式** — Agent 声明能力，框架自动匹配工具
3. **动态加载** — 运行时根据任务类型加载工具

**推荐**: 选项 1，简单直接，符合"减法 > 加法"原则。

### Q5: 向后兼容策略

**选项**：
1. **一次性大重构** — 破坏性变更，更新所有调用方
2. **适配器模式** — 保留旧接口，内部委托给新架构
3. **渐进式迁移** — 新功能用新架构，旧代码逐步迁移

**推荐**: 选项 2 (适配器模式)，降低风险。

---

## 6. 参考文献

1. Anthropic. (2025). [How we built our multi-agent research system](https://www.anthropic.com/news/multi-agent-research)
2. Google Research. (2025). [Towards a Science of Scaling Agent Systems](https://research.google/blog/towards-a-science-of-scaling-agent-systems/)
3. OpenAI. (2025). [Swarm: Lightweight multi-agent orchestration](https://github.com/openai/swarm)
4. LangChain. (2025). [LangChain Multi-Agent Architectures](https://langchain-ai.github.io/langgraph/concepts/multi_agent/)

---

*文档版本: v1.0 | 创建日期: 2026-02-08 | 维护者: DivineSense AI Team*
