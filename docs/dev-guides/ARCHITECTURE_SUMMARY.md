# 架构速览

> DivineSense 核心架构组件和概念快速参考

---

## 五位鹦鹉（内部代理）

| 鹦鹉                   | 领域                 | 配置文件                      |
| :--------------------- | :------------------- | :---------------------------- |
| MemoParrot (灰灰)      | 笔记搜索             | `config/parrots/memo.yaml`    |
| ScheduleParrot (时巧)  | 日程管理             | `config/parrots/schedule.yaml` |
| AmazingParrot (折衷)   | 综合助理             | `config/parrots/amazing.yaml` |
| GeekParrot (极客)      | Claude Code CLI 桥接 | 代码实现                      |
| EvolutionParrot (进化) | 自我进化             | 代码实现                      |

### 路由四层
```
用户输入 → Cache (LRU, 0ms) → Rule (0ms) → History (~10ms) → LLM (~400ms)
```

**注意**：`AUTO` 不是鹦鹉，是"请后端路由"的标记

---

## 关键概念

| 概念 | 实体             | 说明                 |
| :--- | :--------------- | :------------------- |
| 对话 | `AIConversation` | 包含多个 Block       |
| 块   | `AIBlock`        | 一个用户-AI 交互轮次 |
| 代理 | `ParrotAgent`    | 处理请求的 AI 实体   |
| 路由 | `ChatRouter`     | 决定使用哪只鹦鹉     |

### BlockMode vs AgentType

| BlockMode | 说明                     | AgentType | 说明           |
| :--------- | :----------------------- | :--------- | :------------- |
| `normal`   | 普通模式，由后端路由决定 | `AUTO`     | 路由标记       |
| `geek`     | 极客模式                 | `GEEK`     | Claude Code   |
| `evolution`| 进化模式                 | `EVOLUTION`| 自我进化       |

**重要**：Mode 和 Type 是独立维度

---

## UniversalParrot 架构

```
UniversalParrot (配置驱动)
├── ParrotFactory (工厂)
├── ParrotConfig (配置加载)
├── ExecutionStrategy (执行策略)
│   ├── DirectExecutor (原生调用)
│   ├── ReActExecutor (思考-行动循环)
│   └── PlanningExecutor (两阶段规划)
└── ToolRegistry (工具注册表)
```

---

## 执行策略对比

| 策略     | 适用场景         | 特点             |
| :------- | :--------------- | :--------------- |
| `direct` | 简单 CRUD        | 快速，单次调用   |
| `react`  | 多步推理         | 思考-行动循环    |
| `planning`| 复杂多工具协作   | 两阶段，可并行   |

---

*文档路径：docs/dev-guides/ARCHITECTURE_SUMMARY.md*
