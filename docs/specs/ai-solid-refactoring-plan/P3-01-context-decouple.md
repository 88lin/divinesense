# P3-01: ConversationContext 领域解耦

> **阶段**: Phase 3 — DIP / ISP 治理  
> **原则**: ISP (接口隔离) + SRP (单一职责)  
> **风险**: 🟡 中  
> **预计工作量**: 2 天  
> **前置依赖**: P2-01 (scheduler.go 拆分完成)

## 背景

`ai/agents/context.go` 中的 `ConversationContext` 混入了日程领域特有的类型（`ScheduleDraft`、`WorkingState`、`WorkflowStep`、`ExtractRefinement`），并直接 import `services/schedule` 和 `store` 包。这意味着 Memo Agent 等无日程功能的代理也被迫依赖日程领域类型。

## 目标

将 `ConversationContext` 中的领域特有状态抽离为独立类型，通过通用扩展机制 (`Extensions map[string]any`) 实现领域隔离。

## 涉及文件

| 操作   | 文件                                                    |
| :----- | :------------------------------------------------------ |
| MODIFY | `ai/agents/context.go`                                  |
| NEW    | `ai/agents/tools/schedule/context.go`（日程领域上下文） |
| MODIFY | 所有读写 `ConversationContext.WorkingState` 的文件      |

## 改造内容

### Step 1: 修改 ConversationContext

```diff
 type ConversationContext struct {
     SessionID    string
     UserID       int32
     Turns        []ConversationTurn
     LastRoute    ChatRouteType
-    WorkingState *WorkingState
+    Extensions   map[string]any
 }

-type WorkingState struct { ... }
-type ScheduleDraft struct { ... }
-type WorkflowStep string
-type ExtractRefinement struct { ... }
```

### Step 2: 新增领域上下文

```go
// agents/tools/schedule/context.go [NEW]
package schedule

type WorkingState struct {
    ProposedSchedule *ScheduleDraft
    CurrentStep      WorkflowStep
    Conflicts        []*ConflictInfo
}

// 提供类型安全的 getter/setter
func GetWorkingState(ctx *agent.ConversationContext) *WorkingState { ... }
func SetWorkingState(ctx *agent.ConversationContext, ws *WorkingState) { ... }
```

### Step 3: 更新调用方

所有直接读写 `ctx.WorkingState` 的代码改为 `schedule.GetWorkingState(ctx)` / `schedule.SetWorkingState(ctx, ws)`。

## 验收条件

- [ ] `ai/agents/context.go` 不再包含 `WorkingState`、`ScheduleDraft`、`WorkflowStep`、`ExtractRefinement` 类型定义
- [ ] `ai/agents/context.go` 不再 import `services/schedule` 或 `store`
- [ ] `ConversationContext` 包含 `Extensions map[string]any` 字段
- [ ] 日程领域类型存在于 `agents/tools/schedule/context.go`
- [ ] 所有调用方通过类型安全的辅助函数访问日程上下文
- [ ] `go build ./...` 编译通过
- [ ] `go test ./ai/... -count=1` 全部通过
