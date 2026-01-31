# Reference - DivineSense 项目参考文档

> Idea Researcher 的项目上下文、技术栈、代理系统等参考信息。

---

## DivineSense 项目上下文

### 动态上下文发现

> **重要**：项目上下文会随时间变化，必须动态发现而非硬编码。

```bash
# 1. 获取仓库（动态）
git remote get-url origin | sed 's/.*github.com[:/]\(.*\)\.git/\1/'

# 2. 发现可用的 AI 代理（动态）
find plugin/ai/agent -name "*_parrot.go" 2>/dev/null | sed 's/.*\///' | sed 's/_parrot.go//'

# 3. 发现前端页面（动态）
find web/src/pages -name "*.tsx" -exec basename {} .tsx \; 2>/dev/null

# 4. 发现调研报告（动态）
ls docs/research/*.md 2>/dev/null | xargs -I {} basename {} .md

# 5. 读取项目配置（动态）
cat go.mod | grep "^module" | awk '{print $2}'
cat web/package.json | grep '"version"' | head -1
```

### 核心原则

| 原则 | 说明 |
|:-----|:-----|
| **隐私优先** | 所有数据自托管，无遥测 |
| **AI 增强** | 五代理系统 |
| **技术杠杆** | 自动化 > 手动，智能 > 机械 |
| **渐进式** | MVP 先行，迭代完善 |

### 技术栈

| 层级 | 技术 |
|:-----|:-----|
| **后端** | Go 1.25, Echo, Connect RPC |
| **前端** | React 18, Vite 7, TypeScript, Tailwind 4 |
| **数据库** | PostgreSQL (生产，完整 AI) / SQLite (开发，无 AI) |
| **AI** | DeepSeek V3, bge-m3, bge-reranker-v2-m3 |

### 现有功能模块

| 模块 | 路径 | 说明 |
|:-----|:-----|:-----|
| **笔记** | `web/src/pages/Home.tsx` | Markdown 编辑、语义搜索 |
| **日程** | `web/src/pages/Schedule.tsx` | 自然语言创建、冲突检测 |
| **AI 聊天** | `web/src/pages/AIChat.tsx` | 智能路由、多代理 |
| **探索** | `web/src/pages/Explore.tsx` | 搜索和发现内容 |

### 设计约束

| 约束 | 说明 |
|:-----|:-----|
| **i18n 强制** | 所有 UI 文本必须双语（en/zh-Hans） |
| **SQLite 无 AI** | 开发数据库不支持 AI 功能 |
| **原子提交** | 每个 commit 只做一件事 |
| **PR 审查** | 所有变更通过 PR 合并 |

---

## 五代理系统

### 代理概览

> ⚠️ **动态发现**：使用 `find plugin/ai/agent -name "*_parrot.go"` 获取最新列表。

```
ChatRouter
├── EvolutionParrot (进化) - 源代码修改 + PR
├── GeekParrot (极客) - Claude Code CLI 集成
├── MemoParrot (灰灰) - 笔记搜索
├── ScheduleParrotV2 (金刚) - 日程管理
└── AmazingParrot (惊奇) - 综合助理
```

### 代理能力矩阵

| 代理 | 中文名 | 可用工具 | 设计约束 |
|:-----|:-------|:---------|:---------|
| **EvolutionParrot** | 进化 | 源代码修改+PR | 仅管理员可用，需 PR 审查 |
| **GeekParrot** | 极客 | Claude Code CLI | 用户沙箱隔离，无源码访问权限 |
| **MemoParrot** | 灰灰 | memo_search | 需要考虑 embedding 维度(1024)、检索阈值 |
| **ScheduleParrotV2** | 金刚 | add/query/update/find_free_time | 需要考虑冲突检测、时区、周期事件 |
| **AmazingParrot** | 惊奇 | 所有工具 | 需要设计并发策略和错误处理 |

### 代理路由优先级

```
EvolutionMode (最高) → GeekMode → 常规代理路由
```

| 模式 | 目标用户 | 工作目录 | 产出物 |
|:-----|:---------|:---------|:-------|
| **EvolutionMode** | 仅管理员 | 源代码根目录 | GitHub PR |
| **GeekMode** | 所有用户 | 用户沙箱 | 代码产物（用户下载） |
| **常规模式** | 所有用户 | — | 搜索/创建结果 |

---

## Agent 协同规则

### 代理决策树

```
功能涉及笔记语义搜索？
├── 是 → MemoParrot (memo_search)
│         └── 确认：需要新 embedding？需要新检索策略？
└── 否 → 继续判断

功能涉及日程管理？
├── 是 → ScheduleParrotV2
│         └── 确认：需要新工具？
└── 否 → 继续判断

功能涉及代码生成/执行？
├── 是 → GeekParrot/EvolutionParrot
│         ├── GeekMode: 用户沙箱
│         └── EvolutionMode: 源代码（需管理员）
└── 否 → 继续判断

功能跨多个领域？
└── 是 → AmazingParrot（并发或顺序编排）
```

### 协同检查清单

- [ ] 是否需要调用现有 Parrot 的工具？
- [ ] 是否需要为 Parrot 新增工具？
- [ ] 是否需要修改 Parrot 的提示词？
- [ ] 新功能是否会干扰 Parrot 的现有能力？

### 代理间通信模式

| 模式 | 描述 | 场景 | 示例 |
|:-----|:-----|:-----|:-----|
| **顺序** | A 完成后调用 B | 有依赖 | 先搜索笔记，再创建日程 |
| **并发** | A 和 B 同时执行 | 无依赖 | 同时搜索笔记和查询日程 |
| **条件** | 根据结果选择 | 动态决策 | 搜索失败时用 LLM 重新理解 |
| **聚合** | 多个结果汇总 | 综合分析 | 笔记+日程生成周报 |

---

## 工具使用策略

### 工具映射

| 功能 | MCP 工具 | CLI 等效 |
|:-----|:---------|:--------|
| 读取文件 | `Read` | — |
| 搜索代码 | `Grep` | `grep -r` |
| 发现文件 | `Glob` | `find` |
| 网页搜索 | `WebSearch` | — |
| 深度阅读 | `mcp__web-reader__webReader` | `curl` |
| 搜索 Issue | `mcp__plugin_github_github__search_issues` | `gh issue list` |
| 创建 Issue | `mcp__plugin_github_github__issue_write` | `gh issue create` |
| 询问用户 | `AskUserQuestion` | — |

### 使用流程

```
用户 Idea
    → Read/Grep: 检查现有实现
    → search_issues: 检查重复
    → WebSearch/webReader: 竞品分析
    → AskUserQuestion: 澄清需求
    → 生成方案
    → AskUserQuestion: 确认方案
    → issue_write: 创建 Issue
```

---

## 数据库架构

### 核心表

| 表名 | 用途 | 关键列 |
|:-----|:-----|:-----|
| `memo` | 笔记内容 | id, user_id, content |
| `memo_embedding` | 向量嵌入 | memo_id, embedding(vector(1024)) |
| `schedule` | 日程 | id, user_id, start_time |
| `conversation_context` | 会话持久化 | session_id, context_data(JSONB) |

### 数据库策略

| 数据库 | 用途 | AI 支持 |
|:-------|:-----|:-------|
| **PostgreSQL** | 生产 | ✅ 完整 |
| **SQLite** | 开发 | ❌ 不支持 |

---

## 前端架构

### 页面组件

| 路径 | 组件 | 布局 |
|:-----|:-----|:-----|
| `/` | `Home.tsx` | MainLayout |
| `/explore` | `Explore.tsx` | MainLayout |
| `/chat` | `AIChat.tsx` | AIChatLayout |
| `/schedule` | `Schedule.tsx` | ScheduleLayout |

### Tailwind CSS 4 陷阱

> **切勿使用语义化 `max-w-sm/md/lg/xl`** —— 解析为约 16px。

```tsx
// ❌ 错误
<DialogContent className="max-w-md">

// ✅ 正确
<DialogContent className="max-w-[28rem]">
```

---

## Git 工作流

### 分支命名

```
feat/<issue-id>-简短描述     # 功能
fix/<issue-id>-简短描述      # 修复
refactor/<issue-id>-简短描述 # 重构
evolution/<issue-id>-任务名   # 进化
```

### Commit 格式

```
<type>(<scope>): <description>

类型: feat / fix / refactor / perf / docs / test / chore
示例: feat(ai): add intent router
```

---

## Issue 标签

### 功能类型
`enhancement` | `feature-ai` | `feature-web` | `feature-backend`

### 优先级
`priority-high` | `priority-medium` | `priority-low`

### 复杂度
`complexity-small` | `complexity-medium` | `complexity-large` | `complexity-xlarge`

### 代理标签
`parrot-memo` | `parrot-schedule` | `parrot-amazing` | `parrot-geek` | `parrot-evolution`

---

*文档版本：v3.1 | 最后更新：2025-01-31*
