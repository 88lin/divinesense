# 🤝 贡献指南 - Contributing to DivineSense

> **欢迎来到 DivineSense 社区！** 🎉
> 感谢你对 DivineSense 项目的关注！本文档将帮助你以零摩擦的方式开始贡献。

---

## 📖 项目简介

DivineSense (神识) 是一款 AI 驱动的个人第二大脑，采用 **Orchestrator-Workers 多代理架构**，通过智能代理自动化任务、过滤高价值信息、以技术杠杆提升生产力。

### 核心特性

- **Orchestrator-Workers 架构**：LLM 驱动的任务分解与协调
- **专家代理 (Expert Agents)**：MemoParrot (灰灰)、ScheduleParrot (时巧) 等领域专家
- **外部执行器**：GeekParrot (Claude Code CLI)、EvolutionParrot (自我进化)
- **智能路由**：Cache → Rule → History → LLM，响应延迟 0-400ms
- **混合检索**：BM25 + 向量搜索 + RRF 融合
- **Chat Apps 集成**：支持 Telegram、钉钉、WhatsApp
- **单二进制部署**：极简部署，无需 Node.js/Nginx

### 技术栈

| 领域 | 技术选型 |
|:-----|:---------|
| **后端** | Go 1.25, Echo, Connect RPC, pgvector |
| **前端** | React 18, Vite 7, TypeScript, Tailwind CSS 4, Radix UI |
| **数据库** | PostgreSQL 16+（生产），SQLite（开发） |
| **AI** | 智谱 GLM / DeepSeek（对话），SiliconFlow（嵌入/分类/重排） |

---

## 🚀 开发环境搭建

### 前置要求

确保你的开发环境已安装以下工具：

- **Go**: >= 1.25
- **Node.js**: >= 20（推荐使用 `fnm` 或 `nvm` 管理）
- **pnpm**: >= 9（`npm install -g pnpm`）
- **Docker**: 用于运行本地数据库
- **Make**: 构建工具（Windows 用户请使用 WSL2 或 Git Bash）
- **Git**: 版本控制
- **GitHub CLI**: `gh` 命令行工具（推荐）

### 快速启动

```bash
# 1. 克隆项目
git clone https://github.com/hrygo/divinesense.git
cd divinesense

# 2. 安装所有依赖（Backend + Frontend）
make deps-all

# 3. 安装 Git Hooks（必需）
make install-hooks

# 4. 启动基础设施（PostgreSQL Docker）
make docker-up

# 5. 启动开发服务（后端 + 前端）
make start
```

访问 http://localhost:25173

### 验证安装

```bash
# 检查服务状态
make status

# 运行完整检查
make check-all
```

---

## 📂 项目结构

```
divinesense/
├── cmd/divinesense/     # 应用程序入口
├── server/              # HTTP/gRPC 服务器 & 路由
├── ai/                  # AI 核心模块
│   ├── agents/          # 代理系统
│   │   └── orchestrator/  # Orchestrator-Workers 架构
│   ├── routing/         # 智能路由
│   └── core/            # LLM/嵌入核心
├── web/                 # React 前端应用
├── store/               # 数据存储层
├── proto/               # Protobuf 定义
├── config/              # 配置文件（代理提示词等）
├── plugin/              # 插件系统
└── deploy/              # 部署脚本
```

### Orchestrator-Workers 架构

```
用户输入
    ↓
┌─────────────────────────────────────────┐
│            Orchestrator                 │  ← LLM 驱动任务分解
│  ┌─────────────────────────────────┐   │
│  │  Decomposer (任务分解)           │   │
│  │  Executor  (并行执行)            │   │
│  │  Aggregator (结果聚合)           │   │
│  └─────────────────────────────────┘   │
└────────────────┬────────────────────────┘
                 │
        ┌────────┴────────┐
        ↓                 ↓
┌───────────────┐ ┌───────────────┐
│ MemoParrot    │ │ ScheduleParrot│  ← Expert Agents (配置化)
└───────────────┘ └───────────────┘
```

**关键文件**：
| 文件 | 职责 |
|:-----|:-----|
| `ai/agents/orchestrator/orchestrator.go` | 核心编排器 |
| `ai/agents/orchestrator/decomposer.go` | 任务分解（DAG 依赖） |
| `ai/agents/orchestrator/executor.go` | 并行执行 |
| `ai/agents/orchestrator/aggregator.go` | 结果聚合 |
| `config/orchestrator/*.yaml` | 提示词配置 |

---

## 🛠 开发规范

### 代码风格

#### Go 后端

- **文件命名**：`snake_case.go`
- **日志**：使用 `log/slog` 结构化日志
- **错误处理**：始终检查并处理错误
- **注释**：导出函数必须有文档注释

#### React/TypeScript 前端

- **组件**：PascalCase 命名（`UserProfile.tsx`）
- **Hooks**：`use` 前缀（`useUserData()`）
- **样式**：使用 Tailwind CSS 类名
- **类型**：避免 `any`，使用具体类型

### Tailwind CSS 4 关键陷阱

> **⚠️ 切勿使用语义化 `max-w-sm/md/lg/xl`** —— 在 Tailwind v4 中它们解析为约 16px

**错误示例**：
```tsx
<DialogContent className="max-w-md">  // ❌ 坍缩成约 16px
```

**正确示例**：
```tsx
<DialogContent className="max-w-[28rem]">  // ✅ 448px
```

### 国际化 (i18n) 规范

**所有 UI 文本必须双语支持！**

1. **文件位置**：
   - 英文：`web/src/locales/en.json`
   - 简体中文：`web/src/locales/zh-Hans.json`

2. **添加新 Key 步骤**：
   ```tsx
   // 1. 在组件中使用
   const title = t("pages.home.title");

   // 2. 同时添加到两个翻译文件
   // en.json: { "pages": { "home": { "title": "Welcome" } } }
   // zh-Hans.json: { "pages": { "home": { "title": "欢迎" } } }

   // 3. 验证
   make check-i18n
   ```

---

## 🔄 Git 工作流

### 工作流概览

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  创建 Issue  │ → │  创建分支    │ → │  开发提交    │ → │  发起 PR     │ → │  审核合并    │
│  (gh issue) │    │  (checkout -b)│ │  (git commit)│  │  (gh pr create)│ │  (gh pr merge)│
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

### 1. 创建 Issue

```bash
# 创建新 Issue
gh issue create --title "[feat] 添加 AI 路由功能" --body "详细描述..."
```

**Issue 标题格式**：
- 功能：`[feat] 功能描述`
- 修复：`[fix] 问题描述`
- 重构：`[refactor] 重构描述`

### 2. 创建分支

**禁止直接在 `main` 分支修改**。为每个 Issue 创建独立分支。

```bash
# 确保本地 main 是最新的
git checkout main
git pull origin main

# 创建功能分支（引用 Issue 编号）
git checkout -b feat/123-add-ai-router
```

**分支命名规范**：

| 类型 | 格式 | 示例 |
|:-----|:-----|:-----|
| 功能 | `feat/<issue-id>-简短描述` | `feat/123-add-ai-router` |
| 修复 | `fix/<issue-id>-简短描述` | `fix/456-session-cleanup` |
| 重构 | `refactor/<issue-id>-简短描述` | `refactor/789-remove-dead-code` |

### 3. 开发与提交

#### 定期同步 main

```bash
# 每天开始工作前执行
git fetch origin
git rebase origin/main
```

#### 提交规范

我们遵循 **Conventional Commits** 规范：

| 类型 | 范围 | 示例 |
|:-----|:-----|:-----|
| `feat` | 功能区域 | `feat(ai): add intent router` |
| `fix` | Bug 区域 | `fix(db): resolve race condition` |
| `refactor` | 代码区域 | `refactor(frontend): extract hooks` |
| `docs` | 文档 | `docs(readme): update quick start` |
| `test` | 测试 | `test(ai): add agent test cases` |
| `chore` | 日常维护 | `chore(deps): upgrade dependencies` |

### 4. 发起 Pull Request

```bash
gh pr create --title "feat(ai): add intent router" --body "$(cat <<'EOF'
## 概述
添加 AI 意图路由器，支持用户查询自动分类

## 变更内容
- [ ] 实现 ChatRouter 四层路由
- [ ] 添加规则匹配引擎
- [ ] 集成历史上下文匹配

## 关联 Issue
Resolves #123

## 测试计划
- [ ] 本地测试通过
- [ ] 单元测试覆盖率 >80%
- [ ] `make check-all` 通过

## 检查清单
- [ ] 代码遵循项目规范
- [ ] 自我审查代码
- [ ] 注释说明复杂逻辑
- [ ] i18n 翻译已更新（如需要）

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

---

## 📦 提交 PR 流程

### 提交前检查清单

- [ ] **本地检查通过**：`make check-all`
- [ ] **代码自检**：无调试日志、逻辑清晰、错误处理完善
- [ ] **文档更新**：API 文档、README 或注释已更新
- [ ] **翻译更新**：UI 文本已添加到翻译文件

---

## 🛠 常用开发命令

### Makefile 命令速查

| 命令 | 描述 |
|:-----|:-----|
| `make help` | 显示所有可用命令 |
| `make deps-all` | 安装所有依赖 |
| `make docker-up` | 启动数据库容器 |
| `make start` | 同时启动前后端 |
| `make test` | 运行后端测试 |
| `make check-all` | 全量检查 |
| `make check-i18n` | 检查多语言一致性 |

### GitHub CLI 命令

```bash
# Issue 管理
gh issue list              # 列出所有 Issue
gh issue view 123          # 查看 Issue 详情
gh issue close 123         # 关闭 Issue

# PR 管理
gh pr list                 # 列出所有 PR
gh pr view 456             # 查看 PR 详情
gh pr merge 456            # 合并 PR
```

---

## 🆘 获取帮助

### 遇到问题？

1. **先搜索**：
   - 查看 [Issues](https://github.com/hrygo/divinesense/issues)
   - 搜索现有讨论

2. **创建 Issue**：
   ```bash
   gh issue create --interactive
   ```

3. **社区讨论**：
   - [Discussions](https://github.com/hrygo/divinesense/discussions)

### 开发资源

- **项目文档**：
  - [架构设计](docs/dev-guides/ARCHITECTURE.md)
  - [后端指南](docs/dev-guides/BACKEND_DB.md)
  - [前端指南](docs/dev-guides/FRONTEND.md)
  - [部署指南](docs/deployment/BINARY_DEPLOYMENT.md)

- **规范文档**：
  - [Git 工作流](.claude/rules/git-workflow.md)
  - [代码风格](.claude/rules/code-style.md)
  - [国际化规范](.claude/rules/i18n.md)

---

## 📚 相关文档

快速查找相关开发文档：

- **项目概述**：[README.md](README.md)
- **架构设计**：[docs/dev-guides/ARCHITECTURE.md](docs/dev-guides/ARCHITECTURE.md)
- **后端开发**：[docs/dev-guides/BACKEND_DB.md](docs/dev-guides/BACKEND_DB.md)
- **前端开发**：[docs/dev-guides/FRONTEND.md](docs/dev-guides/FRONTEND.md)
- **调试经验**：[docs/research/DEBUG_LESSONS.md](docs/research/DEBUG_LESSONS.md)

---

Happy Coding! 🚀

*最后更新：2026-02-12 (v0.99.0)*
