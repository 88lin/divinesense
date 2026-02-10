# DivineSense 开发者指南索引

> **版本**: v0.97.0 | **更新时间**: 2026-02-10

本文档目录为 DivineSense 开发者提供全面的技术指南和架构文档。

---

## 📖 目录

### 核心架构文档

| 文档 | 描述 | 面向读者 |
|:-----|:-----|:---------|
| **[架构文档](ARCHITECTURE.md)** | 项目整体架构、技术栈、五位鹦鹉代理系统、智能路由 | 所有开发者 |
| **[后端与数据库指南](BACKEND_DB.md)** | Go 后端开发、数据库设计、API 设计模式 | 后端开发者 |
| **[前端开发指南](FRONTEND.md)** | React 组件、布局系统、i18n、状态管理 | 前端开发者 |

### 专题文档

| 文档 | 描述 | 面向读者 |
|:-----|:-----|:---------|
| **[AI Chat 界面架构](AI_CHAT_INTERFACE.md)** | AI 聊天 UI/UX 设计、组件层级、交互流程 | 前端开发者 |
| **[CC Runner 架构](CC_RUNNER_ARCH.md)** | Geek Mode 核心架构、Claude Code CLI 集成 | 后端开发者 |
| **[Quickstart Agent](QUICKSTART_AGENT.md)** | 快速入门 AI 代理开发 | 新手开发者 |

### 系统模块文档

| 文档 | 描述 | 面向读者 |
|:-----|:-----|:---------|
| **[Filter 系统](FILTER_SYSTEM.md)** | 敏感信息过滤、安全策略 | 后端开发者 |
| **[Stats 系统](STATS_SYSTEM.md)** | 成本追踪、指标收集、告警 | 后端开发者 |
| **[Tracing 系统](TRACING_SYSTEM.md)** | OpenTelemetry 分布式追踪、性能分析 | 后端开发者 |
| **[Preload 系统](PRELOAD_SYSTEM.md)** | 智能缓存预加载、性能优化 | 后端开发者 |

---

## 🚀 快速导航

### 我想了解...

**项目整体架构**
→ 从 [架构文档](ARCHITECTURE.md) 开始，了解技术栈、目录结构和五位鹦鹉代理系统

**开发环境搭建**
→ 参考 [贡献指南](../CONTRIBUTING.md) 的"开发环境搭建"章节

**后端开发**
→ 阅读 [后端与数据库指南](BACKEND_DB.md)，了解 API 设计模式、命名约定和数据库架构

**前端开发**
→ 阅读 [前端开发指南](FRONTEND.md)，了解布局架构、组件模式和 i18n 规范

**AI 代理开发**
→ 从 [Quickstart Agent](QUICKSTART_AGENT.md) 开始，然后参考 [架构文档](ARCHITECTURE.md) 中的 Parrot 代理架构章节

**Geek Mode / Claude Code CLI 集成**
→ 参考 [CC Runner 架构](CC_RUNNER_ARCH.md)

**成本追踪和性能优化**
→ 阅读 [Stats 系统](STATS_SYSTEM.md) 和 [Tracing 系统](TRACING_SYSTEM.md)

---

## 📁 文档组织结构

```
docs/
├── dev-guides/           # 开发者指南（当前目录）
│   ├── README.md         # 本索引文件 ✨ 新建
│   │
│   ├── 核心架构文档
│   ├── ARCHITECTURE.md   # 架构文档
│   ├── BACKEND_DB.md     # 后端与数据库
│   ├── FRONTEND.md       # 前端开发
│   │
│   ├── 专题文档
│   ├── AI_CHAT_INTERFACE.md    # AI Chat UI 架构
│   ├── CC_RUNNER_ARCH.md      # CC Runner 异步架构
│   ├── QUICKSTART_AGENT.md    # Agent 开发快速入门
│   │
│   ├── 系统模块文档
│   ├── FILTER_SYSTEM.md       # 敏感信息过滤器
│   ├── STATS_SYSTEM.md        # 统计与告警系统
│   ├── TRACING_SYSTEM.md      # 分布式链路追踪
│   ├── PRELOAD_SYSTEM.md      # 预测缓存预加载
│   │
│   ├── 开发参考文档
│   ├── UNEXPOSED_FEATURES.md  # 未暴露功能清单
│   ├── AGENT_TESTING.md        # AI 代理测试指南
│   ├── COMMON_TASKS.md         # 常见开发任务步骤
│   ├── PROJECT_PATHS.md        # 项目路径参考
│   ├── SQLITE_VEC_USAGE.md     # SQLite 向量搜索使用
│   ├── SQLITE_VEC_USAGE_GUIDE.md  # SQLite vec 详细指南
│   └── TEST_COVERAGE_PLAN.md   # 测试覆盖计划
│
├── deployment/           # 部署指南
│   └── BINARY_DEPLOYMENT.md
│
├── user-guides/          # 用户指南
│   └── CHAT_APPS.md
│
├── research/             # 研究文档
│   └── DEBUG_LESSONS.md
│
└── specs/                # 技术规格
    ├── block-design/
    │   └── unified-block-model.md
    └── cc_runner_async_arch.md
```

### 文档保鲜状态

| 文档 | 状态 | 最后检查版本 |
|:-----|:-----|:-------------|
| README.md | ✅ 新建 | v0.97.0 |
| ARCHITECTURE.md | ✅ 已更新 | v0.97.0 |
| BACKEND_DB.md | ✅ 已更新 | v0.97.0 |
| FRONTEND.md | ✅ 已更新 | v0.97.0 |
| AI_CHAT_INTERFACE.md | ✅ 已验证 | v0.97.0 |
| CC_RUNNER_ARCH.md | ✅ 已完成 | v0.97.0 |
| QUICKSTART_AGENT.md | ✅ 已验证 | v0.97.0 |
| FILTER_SYSTEM.md | ✅ 已完成 | v0.97.0 |
| STATS_SYSTEM.md | ✅ 已完成 | v0.97.0 |
| TRACING_SYSTEM.md | ✅ 已完成 | v0.97.0 |
| PRELOAD_SYSTEM.md | ✅ 已完成 | v0.97.0 |
| UNEXPOSED_FEATURES.md | ✅ 已验证 | v0.97.0 |
| AGENT_TESTING.md | ✅ 已验证 | v0.97.0 |
| COMMON_TASKS.md | ✅ 已验证 | v0.97.0 |

---

## 🔗 相关资源

### 规范文档

| 文档 | 位置 | 描述 |
|:-----|:-----|:-----|
| **Git 工作流** | `../../.claude/rules/git-workflow.md` | 分支管理、PR 规范 |
| **代码风格** | `../../.claude/rules/code-style.md` | Go/React 代码规范 |
| **国际化规范** | `../../.claude/rules/i18n.md` | i18n 最佳实践 |

### 项目根文档

| 文档 | 描述 |
|:-----|:-----|
| [README.md](../../README.md) | 项目概述、快速开始 |
| [CONTRIBUTING.md](../../CONTRIBUTING.md) | 贡献指南、开发规范 |
| [CHANGELOG.md](../../CHANGELOG.md) | 版本更新记录 |

---

## 📝 文档维护

### 保鲜状态

所有文档应保持"保鲜状态"：
- 标注最后检查版本和日期
- 重大变更时同步更新文档
- 发现过时内容及时修正

### 更新规范

1. **版本同步**：文档中引用的版本号应与 git tag 一致（单一版本来源）
   ```bash
   # 获取当前项目版本
   git tag -l | sort -V | tail -1 | sed 's/^v//'
   ```
2. **交叉引用**：使用相对路径维护文档间的链接
3. **示例代码**：所有代码示例应可实际运行

---

*最后更新：2026-02-10*
