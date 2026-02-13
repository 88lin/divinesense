# P2-01: scheduler.go 拆分

> **阶段**: Phase 2 — SRP 治理  
> **原则**: SRP (单一职责)  
> **风险**: 🟡 中  
> **预计工作量**: 3 天  
> **前置依赖**: P1-01 ~ P1-06 (Phase 1 完成)

## 背景

`ai/agents/tools/scheduler.go` 为 38KB 超大文件（~1200 行），集中了日程的解析、验证、创建、查询、冲突检测、批量操作全部逻辑。任何子功能的修改都必须触碰此文件，测试编写困难，违反 SRP。

## 目标

将 `scheduler.go` 按职责拆分为独立文件，放入 `agents/tools/schedule/` 子目录，每个文件 < 500 行。

## 目标结构

```
agents/tools/schedule/        # [NEW] 子包
├── tool.go                   # SchedulerTool 入口 + Run() 路由分派
├── parser.go                 # 时间/日程解析逻辑
├── creator.go                # 日程创建（单条 + 批量）
├── query.go                  # 日程查询
├── conflict.go               # 冲突检测与解决
├── validator.go              # 输入验证
└── formatter.go              # 结果格式化输出
```

## 涉及文件

| 操作   | 文件                                                  |
| :----- | :---------------------------------------------------- |
| NEW    | `agents/tools/schedule/` 下全部 7 个文件              |
| DELETE | `agents/tools/scheduler.go`（内容已迁移）             |
| MODIFY | 引用 `tools.SchedulerTool` 的文件（更新 import 路径） |

## 改造步骤

1. **创建** `agents/tools/schedule/` 子包
2. **提取** `tool.go`：保留 `SchedulerTool` 结构体、`Name()`、`Schema()`、`Run()` 方法。`Run()` 内部按操作类型分派到各子模块
3. **提取** `parser.go`：时间解析（`parseTimeExpression`、`parseRelativeTime` 等）
4. **提取** `creator.go`：`createSchedule`、`createBatchSchedule` 等
5. **提取** `query.go`：`querySchedule`、`listSchedules` 等
6. **提取** `conflict.go`：`checkConflict`、`resolveConflict` 等
7. **提取** `validator.go`：输入字段验证逻辑
8. **提取** `formatter.go`：结果格式化和用户友好消息构建
9. **保持** `ToolWithSchema` 接口合约不变，确保注册中心无感知
10. **删除** 原 `scheduler.go`

## 验收条件

- [ ] `agents/tools/scheduler.go` 已删除
- [ ] `agents/tools/schedule/` 下存在 ≥ 5 个文件
- [ ] 每个新文件不超过 500 行
- [ ] `ToolWithSchema` 接口合约未变更（工具注册无需修改）
- [ ] `go build ./...` 编译通过
- [ ] `go test ./ai/agents/tools/... -count=1` 全部通过
- [ ] `go vet ./ai/agents/tools/...` 无警告
- [ ] 各子模块之间无循环引用

## 注意事项

- 此任务仅做代码移动和组织，**不改变业务逻辑**
- 建议每个 PR 只移动一个子模块，以便 code review
