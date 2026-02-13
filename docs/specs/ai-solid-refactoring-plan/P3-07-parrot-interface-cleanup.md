# P3-07: ParrotAgent 接口方法重叠治理

> **阶段**: Phase 3 — DIP / ISP 治理  
> **原则**: LSP (里氏替换) + ISP (接口隔离)  
> **风险**: 🟡 中  
> **预计工作量**: 1 天  
> **前置依赖**: P3-05 (EventCallback 统一完成)

## 背景

`ParrotAgent` 接口（`base_parrot.go` L326-333）中存在方法重叠：

```go
type ParrotAgent interface {
    Name() string
    Execute(ctx context.Context, userInput string, callback EventCallback) error
    ExecuteWithCallback(ctx context.Context, userInput string, history []string, callback EventCallback) error
    SelfDescribe() *ParrotSelfCognition
}
```

问题：
1. **方法重叠**: `Execute` 仅是 `ExecuteWithCallback(ctx, input, nil, callback)` 的简写形式
2. **所有实现均复制此模式**: `universal_parrot.go`、`geek/parrot.go`、`geek/evolution.go` 的 `Execute` 方法都是单行委托
3. **实现负担**: 每个 ParrotAgent 实现必须同时实现两个方法，增加 ~5 行样板代码

即：`Execute` 和 `ExecuteWithCallback` 的唯一区别是 `history []string` 参数，前者永远传 `nil`。

## 目标

消除方法重叠，统一为单一执行入口，将 `history` 作为可选参数。

## 涉及文件

| 操作   | 文件                                                           |
| :----- | :------------------------------------------------------------- |
| MODIFY | `ai/agents/base_parrot.go`（ParrotAgent 接口定义）             |
| MODIFY | `ai/agents/universal/universal_parrot.go`（删除 Execute 简写） |
| MODIFY | `ai/agents/geek/parrot.go`（删除 Execute 简写）                |
| MODIFY | `ai/agents/geek/evolution.go`（删除 Execute 简写）             |
| MODIFY | 所有调用 `Execute` 的消费方（改为传 nil history）              |

## 改造内容

### 方案（推荐）: 合并为单方法，history 默认 nil

```diff
 type ParrotAgent interface {
     Name() string
-    Execute(ctx context.Context, userInput string, callback EventCallback) error
-    ExecuteWithCallback(ctx context.Context, userInput string, history []string, callback EventCallback) error
+    Execute(ctx context.Context, userInput string, history []string, callback events.Callback) error
     SelfDescribe() *ParrotSelfCognition
 }
```

消费方原来调用 `Execute(ctx, input, callback)` 的改为 `Execute(ctx, input, nil, callback)`。

### 各实现简化

```diff
 // universal_parrot.go
-func (p *UniversalParrot) Execute(ctx context.Context, userInput string, callback EventCallback) error {
-    return p.ExecuteWithCallback(ctx, userInput, nil, callback)
-}
-func (p *UniversalParrot) ExecuteWithCallback(ctx, userInput, history, callback) error {
+func (p *UniversalParrot) Execute(ctx context.Context, userInput string, history []string, callback events.Callback) error {
     // ... 原 ExecuteWithCallback 实现
 }
```

## 验收条件

- [ ] `ParrotAgent` 接口仅包含一个执行方法 `Execute`
- [ ] 所有 ParrotAgent 实现仅需实现一个 `Execute` 方法
- [ ] 不需要 history 的调用方传 `nil`
- [ ] 全量搜索 `ExecuteWithCallback` 在接口定义和实现中返回 0 处结果
- [ ] `go build ./...` 编译通过
- [ ] `go test ./ai/agents/... -count=1` 全部通过

## 注意事项

- 所有实现中的 `Execute` 都是 `ExecuteWithCallback` 的一行委托，删除后不影响逻辑
- 需全量搜索 `.Execute(` 调用确认有多少消费方需新增 `nil` 参数
