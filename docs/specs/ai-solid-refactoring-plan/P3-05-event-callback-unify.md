# P3-05: EventCallback 类型统一

> **阶段**: Phase 3 — DIP / ISP 治理  
> **原则**: DRY + LSP (里氏替换)  
> **风险**: 🟡 中  
> **预计工作量**: 1 天  
> **前置依赖**: 无

## 背景

事件回调相关类型在代码库中存在 **3 个独立定义**，签名不一致：

| 类型               | 位置                                | 签名                      | 角色                         |
| :----------------- | :---------------------------------- | :------------------------ | :--------------------------- |
| `EventCallback`    | `agents/runner/runner.go` L62       | `func(string, any) error` | 核心回调（可传播错误）       |
| `SafeCallbackFunc` | `agents/runner/runner.go` L68       | `func(string, any)`       | 包装回调（吞掉错误）         |
| `EventCallback`    | `agents/orchestrator/types.go` L176 | `func(string, string)`    | 简化回调（string-only data） |

> [!IMPORTANT]
> `SafeCallbackFunc` 定义在 `runner/runner.go`（不是 `base_parrot.go`），签名为 `func(string, any)` 无返回值。  
> `base_parrot.go` 仅在 `SendPhaseChange`/`SendProgress` 中使用它。  
> `cc_runner.go` L148 是 `runner.SafeCallbackFunc` 的别名。

### 设计关系

```
EventCallback (有 error 返回)
    ↓ SafeCallback() 包装
SafeCallbackFunc (无 error 返回，内部 log error)
```

`SafeCallback` 函数 (`runner/runner.go` L74) 将 `EventCallback` 包装为 `SafeCallbackFunc`，吞掉 error 并记录日志。

## 目标

定义统一的 `EventCallback` 类型，所有包引用同一定义。保留 `SafeCallbackFunc` 作为便捷包装。

## 涉及文件

| 操作   | 文件                                                    |
| :----- | :------------------------------------------------------ |
| NEW    | `ai/agents/events/callback.go`                          |
| MODIFY | `ai/agents/runner/runner.go`（改为引用 events 包）      |
| MODIFY | `ai/agents/orchestrator/types.go`（改为引用 events 包） |
| MODIFY | 所有使用旧回调类型的调用方                              |

## 改造内容

### Step 1: 定义统一类型

```go
// ai/agents/events/callback.go [NEW]
package events

// Callback 是统一的事件回调类型
type Callback func(eventType string, eventData any) error

// SafeCallback 是不传播错误的回调包装，用于非关键事件
type SafeCallback func(eventType string, eventData any)

// NoopCallback 不做任何事的回调
var NoopCallback Callback = func(string, any) error { return nil }

// WrapSafe 将 Callback 包装为 SafeCallback（内部记录 error）
func WrapSafe(cb Callback) SafeCallback { ... }
```

### Step 2: 各包改为类型别名或直接引用

```diff
 // agents/runner/runner.go
-type EventCallback func(eventType string, eventData any) error
-type SafeCallbackFunc func(eventType string, eventData any)
+type EventCallback = events.Callback
+type SafeCallbackFunc = events.SafeCallback

 // agents/orchestrator/types.go
-type EventCallback func(eventType string, eventData string)
+type EventCallback = events.Callback
 // 注意：orchestrator 原签名 data 为 string，需将内部传入改为 any
```

### Step 3: 适配 orchestrator 的签名差异

orchestrator 原来使用 `(string, string)` 签名，需要将内部传入字符串 data 的地方直接传入 `any`，接收方从 `any` 中 type-assert 为 `string`。

## 验收条件

- [ ] `ai/agents/events/callback.go` 存在且导出 `Callback` 和 `SafeCallback` 类型
- [ ] `agents/runner/runner.go` 中 `EventCallback` 为 `events.Callback` 的别名或已删除
- [ ] `agents/runner/runner.go` 中 `SafeCallbackFunc` 为 `events.SafeCallback` 的别名或已删除
- [ ] `agents/orchestrator/types.go` 中 `EventCallback` 为 `events.Callback` 的别名或已删除
- [ ] 全量搜索 `type EventCallback func` 和 `type SafeCallbackFunc func` 返回 0 处结果
- [ ] `go build ./...` 编译通过
- [ ] `go test ./ai/... -count=1` 全部通过

## 注意事项

- orchestrator 的签名差异（`string` vs `any`、无 error 返回）需特别关注适配
- `SafeCallbackFunc` 是故意"吞掉 error"的设计，统一后应保留此包装层
