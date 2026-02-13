# P2-02: runner.go 拆分

> **阶段**: Phase 2 — SRP 治理  
> **原则**: SRP (单一职责)  
> **风险**: 🟡 中  
> **预计工作量**: 3 天  
> **前置依赖**: P1-03 (cc_runner.go 清理完成)

## 背景

`ai/agents/runner/runner.go` 为 35KB 超大文件（~1100 行），集中了 CCRunner 的会话管理、子进程控制、流式输出解析、System Prompt 构建、健康检查等多类职责。

## 目标

将 `runner.go` 按职责拆分为同包内的多个文件，每个文件 < 500 行。

## 目标结构

```
agents/runner/
├── runner.go            # CCRunner 核心生命周期（Init/Start/Stop/Execute）
├── process.go           # 子进程管理（启动/停止/重启）
├── stream.go            # 流式输出解析与事件分发
├── prompt.go            # System prompt 构建与模板渲染
├── conversation.go      # 会话消息历史管理
└── health.go            # 健康检查与自愈逻辑
```

## 涉及文件

| 操作   | 文件                                                |
| :----- | :-------------------------------------------------- |
| MODIFY | `agents/runner/runner.go`（缩减为核心生命周期管理） |
| NEW    | `agents/runner/process.go`                          |
| NEW    | `agents/runner/stream.go`                           |
| NEW    | `agents/runner/prompt.go`                           |
| NEW    | `agents/runner/conversation.go`                     |
| NEW    | `agents/runner/health.go`                           |

## 改造步骤

1. **分析** `runner.go` 中的方法，按职责域分类
2. **提取** 子进程管理相关方法到 `process.go`
3. **提取** 流式输出解析（StreamParser 等）到 `stream.go`
4. **提取** prompt 构建逻辑到 `prompt.go`
5. **提取** 会话消息管理到 `conversation.go`
6. **提取** 健康检查 / 自愈到 `health.go`
7. **保留** `runner.go` 仅包含 `CCRunner` 结构体定义和核心生命周期方法

## 验收条件

- [ ] `agents/runner/runner.go` 不超过 500 行
- [ ] `agents/runner/` 下新增 ≥ 4 个文件
- [ ] 所有新文件均属于 `package runner`（同包拆分，无 import 变更）
- [ ] `go build ./...` 编译通过
- [ ] `go test ./ai/agents/runner/... -count=1` 全部通过
- [ ] `go vet ./ai/agents/runner/...` 无警告

## 注意事项

- 同包拆分不影响外部 API，消费方无需修改
- 仅移动代码，**不改变业务逻辑**
- 私有方法和字段在包内仍可访问，无可见性问题
