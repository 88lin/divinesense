# 开发工作流

> DivineSense 项目的标准开发流程和命令

---

## 多任务管理

**何时创建 TODO LIST**：3+ 个优化点、任务 > 1 小时

```
TaskCreate → TaskList → TaskUpdate(in_progress) → TaskUpdate(completed)
```

---

## 开发命令

### 服务控制
```bash
make start              # 启动所有服务（PostgreSQL + 后端 + 前端）
make stop               # 停止所有服务
make restart            # 重启后端服务（需手动确认）
make run                # 仅启动后端
make web                # 仅启动前端
```

### 数据库
```bash
make db-shell           # 进入 PostgreSQL shell
make db-reset           # 重置数据库（破坏性）
make db-vector          # 验证 pgvector 扩展
```

### 检查和测试
```bash
make check-all          # 完整检查（格式 + vet + 测试 + i18n）
make check-i18n         # i18n 完整性检查
make ci-check           # 模拟 CI 检查
make test               # 运行所有测试
```

### 构建
```bash
make build              # 构建后端二进制
make build-web          # 构建前端静态资源
make build-all          # 同时构建前后端
```

---

## 服务重启规范

**⚠️ 禁止直接执行启停命令**

修改后端代码后，通知用户手动执行 `make restart`

---

## Git 提交流程

```
make check-all → feat/fix 分支 → PR → 合并
```

详细规范：@.claude/rules/git-workflow.md

### 分支命名
- 功能：`feat/<issue-id>-description`
- 修复：`fix/<issue-id>-description`
- 重构：`refactor/<issue-id>-description`

### Commit 格式
```
<type>(<scope>): <description>

Refs #<issue-id>

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
```

---

## 前端开发

### 从 web/ 目录运行
```bash
pnpm dev               # 开发服务器
pnpm build             # 生产构建
pnpm lint              # Lint 检查
pnpm lint:fix          # 自动修复
```

### 从项目根目录
```bash
make web               # 启动前端
make build-web         # 构建前端
```

---

## 后端开发

```bash
make run               # 启动后端
make lint              # golangci-lint
make vet               # go vet
make test-ai           # AI 相关测试
```

---

## Proto 变更流程

1. 修改 `.proto` 文件
2. 运行 `make generate` 重新生成代码
3. 更新前后端绑定
4. 提交变更

---

*文档路径：docs/dev-guides/WORKFLOW.md*
