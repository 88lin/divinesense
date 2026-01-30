# CLAUDE.md

> Claude Code 辅助开发的主要上下文文档。详细指南请参见下方的文档索引。

## 产品愿景

**DivineSense (神识)**：AI 代理驱动的个人「第二大脑」—— 通过智能代理自动化任务、过滤高价值信息、以技术杠杆提升生产力。

---

## 快速开始

**默认端口**：`make start` → localhost:25173 (前端) / 28081 (后端) / 25432 (PostgreSQL)

| 命令             | 作用                                  |
| :--------------- | :------------------------------------ |
| `make start`     | 启动全栈 (PostgreSQL + 后端 + 前端)   |
| `make stop`      | 停止所有服务                          |
| `make test`      | 运行后端测试                          |
| `make build-all` | 构建二进制 + 静态资源                 |
| `make check-all` | 运行所有预提交检查 (构建、测试、i18n) |

**技术栈**：Go 1.25 + React 18 (Vite/Tailwind 4) + PostgreSQL (生产) / SQLite (开发)

---

## 核心规则

### 1. 国际化 (i18n)
- **所有 UI 文本必须使用 `t("key")`** —— 禁止硬编码字符串
- 翻译 key 必须同时存在于 `en.json` 和 `zh-Hans.json`
- 验证命令：`make check-i18n`

### 2. 数据库策略
- **PostgreSQL**：生产环境，完整 AI 支持 (pgvector)
- **SQLite**：开发环境，**不支持 AI 功能**
- 测试 AI 相关功能时务必使用 PostgreSQL

### 3. 代码风格

**Go**：
- 文件命名：`snake_case.go`
- 日志：使用 `log/slog` 结构化日志
- 遵循标准 Go 项目布局

**React/TypeScript**：
- 组件：PascalCase 命名 (`UserProfile.tsx`)
- Hooks：`use` 前缀 (`useUserData()`)
- 样式：使用 Tailwind CSS 类名 (参见下方)

**AI 路由**：
- 后端 `ChatRouter` 处理意图分类
- 位置：`plugin/ai/agent/chat_router.go`
- 规则匹配 (0ms) → LLM 降级 (~400ms)
- 路由到：MEMO / SCHEDULE / AMAZING 代理

### 4. Tailwind CSS 4 — 关键

> **切勿使用语义化 `max-w-sm/md/lg/xl`** —— 在 Tailwind v4 中它们解析为约 16px。
>
> **请使用显式值**：`max-w-[24rem]`, `max-w-[28rem]` 等

详见 `docs/dev-guides/FRONTEND.md` 中的 Tailwind v4 陷阱说明。

### 5. Git 约定 (严格执行)

**所有开发与进化模式均需遵守以下流程：**

1.  **强制分支开发**：禁止直接在 `main` 分支修改。
    *   新功能：`feat/功能名称`
    *   修复：`fix/问题描述`
    *   进化：`evolution/task-id`
2.  **强制 Pull Request (PR)**：所有变更必须通过 PR 合并。
    *   **禁止** `git push origin main`
    *   **禁止** 本地直接 `git merge` 到 main
3.  **提交前检查**：必须在本地通过 Lint 和 Test。
    *   运行 `make check-all` (包含 fmt, vet, test, i18n-check)
    *   只有检查通过的代码才允许 commit

遵循约定式提交：

| 类型       | 范围     | 示例                             |
| :--------- | :------- | :------------------------------- |
| `feat`     | 功能区域 | `feat(ai): 添加意图路由器`       |
| `fix`      | Bug 区域 | `fix(db): 修复竞态条件`          |
| `refactor` | 代码区域 | `refactor(frontend): 提取 hooks` |
| `perf`     | N/A      | `perf(query): 优化向量搜索`      |
| `docs`     | N/A      | `docs(readme): 更新快速开始`     |

**格式**：`<type>(<scope>): <description>`

**始终包含**：
```
Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
```

### 6. 测试

- 提交前运行 `make test`
- AI 测试：`make test-ai`
- 数据库变更后验证迁移

---

## 文档索引

| 领域     | 文件                                   | 参考时机                           |
| :------- | :------------------------------------- | :--------------------------------- |
| **后端** | `docs/dev-guides/BACKEND_DB.md`        | API 设计、数据库、Docker、环境变量 |
| **前端** | `docs/dev-guides/FRONTEND.md`          | 布局结构、Tailwind 陷阱、组件模式  |
| **架构** | `docs/dev-guides/ARCHITECTURE.md`      | 项目结构、AI 代理、数据流          |
| **部署** | `docs/deployment/BINARY_DEPLOYMENT.md` | 二进制部署、Geek Mode 配置         |

---

## 关键项目路径

| 领域       | 路径                                                             | 用途                                  |
| :--------- | :--------------------------------------------------------------- | :------------------------------------ |
| API 处理器 | `server/router/api/v1/`                                          | REST/Connect RPC 端点                 |
| AI 代理    | `plugin/ai/agent/`                                               | Parrot 代理 (MEMO、SCHEDULE、AMAZING) |
| AI 服务    | `plugin/ai/{memory,router,vector,aitime,cache,metrics,session}/` | AI 基础设施                           |
| 查询引擎   | `server/queryengine/`                                            | 混合 RAG 检索 (BM25 + 向量)           |
| 前端页面   | `web/src/pages/`                                                 | 页面组件                              |
| 布局       | `web/src/layouts/`                                               | 共享布局组件                          |
| 数据库模型 | `store/db/postgres/`                                             | PostgreSQL 模型                       |
| 数据库迁移 | `store/migration/postgres/`                                      | 数据库迁移                            |
| 发布脚本   | `scripts/release/`                                               | 构建和打包发布二进制                  |
| 部署       | `deploy/aliyun/`                                                 | 阿里云部署 (Docker/二进制模式)        |

---

## 常见任务

### 添加新 API 端点
1. 在 `server/router/api/v1/` 创建处理器
2. 在 `server/router/api/v1/routes.go` 添加路由
3. 如使用 Connect RPC，更新 proto 文件
4. 运行 `make check-build` 验证

### 添加新前端页面
1. 在 `web/src/pages/` 创建组件
2. 在 `web/src/router/` 添加路由
3. 向 `en.json` 和 `zh-Hans.json` 添加 i18n key
4. 运行 `make check-i18n` 验证

### 修改数据库 Schema
1. 在 `store/migration/postgres/` 创建迁移
2. 更新 `store/db/postgres/` 中的模型
3. 测试：`make db-reset`（仅限开发环境！）
4. 运行 `make test` 验证

### 添加 AI 功能
1. 确定代理类型 (MEMO/SCHEDULE/AMAZING)
2. 更新 `plugin/ai/agent/` 中的代理
3. 在 `chat_router.go` 添加路由规则
4. 使用 PostgreSQL 测试（AI 功能必需）

### 构建发布二进制
1. 更新 `internal/version/version.go` 中的版本号
2. 运行 `make release-all VERSION=v1.0.0`
3. 在 `dist/` 查找二进制，在 `releases/` 查找安装包

---

## 提交前检查清单

提交前运行：

```bash
make check-all
```

这会验证：
- 构建通过 (`go build ./...`)
- 测试通过 (`go test ./...`)
- i18n key 完整

---

## 环境变量

主要 `.env` 变量（详见 `.env.example`）：

| 变量                                | 用途                | 默认值        |
| :---------------------------------- | :------------------ | :------------ |
| `DIVINESENSE_DRIVER`                | 数据库驱动          | `postgres`    |
| `DIVINESENSE_DSN`                   | 数据库连接字符串    | —             |
| `DIVINESENSE_AI_ENABLED`            | 启用 AI 功能        | `false`       |
| `DIVINESENSE_AI_EMBEDDING_PROVIDER` | 向量化 API 提供商   | `siliconflow` |
| `DIVINESENSE_AI_LLM_PROVIDER`       | LLM 提供商          | `deepseek`    |
| `SILICONFLOW_API_KEY`               | SiliconFlow API key | —             |
| `DEEPSEEK_API_KEY`                  | DeepSeek API key    | —             |
| `OPENAI_API_KEY`                    | OpenAI API key      | —             |

---

## 故障排查

| 问题              | 解决方案                                                      |
| :---------------- | :------------------------------------------------------------ |
| AI 不可用         | 确保 PostgreSQL 运行且 `DIVINESENSE_AI_ENABLED=true`          |
| Tailwind 样式异常 | 使用显式值 (`max-w-[24rem]`) 而非语义化 (`max-w-md`)          |
| i18n 检查失败     | 向 `web/src/locales/en.json` 和 `zh-Hans.json` 添加缺失的 key |
| 构建失败          | 运行 `make deps` 更新 Go 模块                                 |
| 测试失败          | 确保 PostgreSQL 运行在 25432 端口                             |

---

## 元认知系统

> **元认知 (Meta-Cognition)** — DivineSense 对自身知识状态、检索质量、代理决策的监控与反思机制，用于持续优化系统性能。

### 1. AI 代理自我进化机制

**现状基础**：
- **三层路由**：规则匹配 (0ms) → 历史匹配 (~10ms) → LLM 分类 (~400ms)
- **历史学习**：成功的路由决策保存到 `EpisodicMemory`，供未来模式识别
- **A/B 测试框架**：支持提示词版本实验与自动推荐 (`plugin/ai/agent/prompts.go`)

**进化方向**：
- **主动学习**：收集用户显式反馈（👍/👎）用于路由规则调优
- **个性化适配**：基于用户交互模式优化代理选择权重
- **冷启动优化**：预置常见查询模式，减少新用户对 LLM 的依赖

### 2. 系统元认知监控

**可观测性栈** (`server/internal/observability/`)：
- **结构化日志**：`log/slog` 请求追踪 (request_id, user_id, duration_ms)
- **指标聚合**：内存 Aggregator → PostgreSQL Persister → Service API
- **质量评估**：检索结果分等 (High/Medium/Low Quality) 基于：
  - 前 2 名分数差距 > 0.20 → High
  - 第 1 名分数 > 0.90 → High
  - 第 1 名分数 > 0.70 → Medium

**监控指标**：
| 维度 | 指标                        | 用途         |
| :--- | :-------------------------- | :----------- |
| 路由 | 各层命中率、置信度          | 优化路由规则 |
| 检索 | BM25/向量权重、RRF 融合效果 | 调整检索策略 |
| 代理 | 工具成功率、P95 延迟        | 识别性能瓶颈 |
| 质量 | 结果相关性、用户反馈        | 评估整体效能 |

**改进空间**：
- 实时监控面板（前端可视化）
- 基于阈值的告警机制
- Prometheus 指标导出

### 3. 文档自我进化 (Claude Code)

**原则**：CLAUDE.md 应随项目演进自动更新，而非静态文档。

**触发更新时机**：
- 新增 AI 代理类型 → 更新 `AI 路由` 章节
- 数据库 Schema 变更 → 更新 `关键项目路径`
- 新增 Make 命令 → 更新 `快速开始` 表格
- 架构重构 → 更新 `文档索引`

**维护协议**：
1. 代码变更时，同步检查相关文档章节
2. 文档过时标记 `[需更新]`，在下次提交时修复
3. 遵循原子提交：代码变更与文档更新分属不同 commit

### 4. 外部借鉴与吸收

**参考项目**：
- **OpenClaw** (`openclaw/openclaw`) — 个人 AI 助手（原 ClawdBot）
  - 自进化能力：自主编写新技能代码
  - 多平台集成：WhatsApp, Telegram, Discord, Signal
  - 本地优先：支持 Ollama 本地 LLM
  - 主动交互：可主动发起消息

**吸收策略**：
- **定期审视**：每季度评估顶级开源项目的新特性
- **适配原则**：仅吸收符合 DivineSense 架构的精华
- **保持简洁**：避免引入不必要的复杂度

**参考清单**：
| 项目     | 关注点     | 借鉴方向         |
| :------- | :--------- | :--------------- |
| OpenClaw | 自进化技能 | 代理自主能力扩展 |
| Leon AI  | 模块化技能 | 技能系统架构     |
| Electric | 实时同步   | 多端数据一致性   |

---

## 产品功能概览

### 笔记管理
- Markdown 编辑器（KaTeX, Mermaid, GFM）
- 智能标签系统（AI 推荐）
- 语义搜索（BM25 + 向量）
- 笔记关联与重复检测

### 日程管理
- 自然语言创建
- 智能冲突检测
- 多视图日历（月/周/日/列表）
- 周期事件

### AI 智能代理
- **灰灰** (MemoParrot)：笔记检索
- **金刚** (ScheduleParrot)：日程管理
- **惊奇** (AmazingParrot)：综合助理

### 其他功能
- 知识图谱可视化
- 每日回顾（间隔重复）
- Geek Mode（Claude Code CLI 集成）
