# P3-02: error_class.go 依赖反转

> **阶段**: Phase 3 — DIP / ISP 治理  
> **原则**: DIP (依赖倒置)  
> **风险**: 🟡 中  
> **预计工作量**: 1 天  
> **前置依赖**: 无

## 背景

`ai/agents/error_class.go` 直接导入了上层包：

- `server/service/schedule`（`schedule.ErrScheduleConflict`）
- `store/db/postgres`（`postgresstore.ConflictConstraintError`）

这违反了分层架构原则：**AI 层不应依赖 Server/Store 层**。

## 目标

通过引入接口抽象消除 `error_class.go` 对 `server/` 和 `store/` 包的导入。

## 涉及文件

| 操作   | 文件                                                |
| :----- | :-------------------------------------------------- |
| MODIFY | `ai/agents/error_class.go`                          |
| MODIFY | `server/service/schedule/` 中的错误类型（实现接口） |
| MODIFY | `store/db/postgres/` 中的错误类型（实现接口）       |

## 改造内容

### Step 1: 在 agents 包中定义接口

```go
// ai/agents/error_class.go
type ConflictError interface {
    error
    IsConflict() bool
}
```

### Step 2: 修改错误分类逻辑

```diff
 func ClassifyError(err error) *ClassifiedError {
-    if errors.Is(err, schedule.ErrScheduleConflict) { ... }
-    var conflictErr *postgresstore.ConflictConstraintError
-    if errors.As(err, &conflictErr) { ... }
+    var conflictErr ConflictError
+    if errors.As(err, &conflictErr) && conflictErr.IsConflict() {
+        return &ClassifiedError{
+            Class:      ErrorClassConflict,
+            Original:   err,
+            ActionHint: "find_free_time",
+        }
+    }
 }
```

### Step 3: 让上层类型实现接口

在 `server/service/schedule` 和 `store/db/postgres` 中为现有的冲突错误类型添加 `IsConflict() bool` 方法。

## 验收条件

- [ ] `ai/agents/error_class.go` 不再 import `server/service/schedule`
- [ ] `ai/agents/error_class.go` 不再 import `store/db/postgres`
- [ ] `ConflictError` 接口已定义并在 `ClassifyError` 中使用
- [ ] `schedule.ErrScheduleConflict` 和 `postgresstore.ConflictConstraintError` 实现了 `ConflictError` 接口（或 wrapped error 已实现）
- [ ] `go build ./...` 编译通过
- [ ] `go test ./ai/agents/... -count=1` 全部通过
- [ ] `go test ./server/... -count=1` 全部通过（确保上层适配无误）
