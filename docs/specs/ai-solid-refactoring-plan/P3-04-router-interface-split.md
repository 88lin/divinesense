# P3-04: RouterService 接口拆分

> **阶段**: Phase 3 — DIP / ISP 治理  
> **原则**: ISP (接口隔离)  
> **风险**: 🟡 中  
> **预计工作量**: 1 天  
> **前置依赖**: 无

## 背景

当前 `routing.RouterService` 接口同时包含意图分类、模型选择和统计三类职责：

```go
type RouterService interface {
    ClassifyIntent(ctx, input) (Intent, float32, bool, error)
    SelectModel(ctx, task) (ModelConfig, error)
    RecordFeedback(ctx, feedback) error
    GetRouterStats(ctx, userID, timeRange) (*RouterStats, error)
}
```

不需要模型选择的消费者也被迫依赖 `SelectModel` 方法，违反 ISP。

## 目标

将 `RouterService` 拆分为更小的接口，使消费者只依赖所需的方法子集。

## 涉及文件

| 操作   | 文件                                                               |
| :----- | :----------------------------------------------------------------- |
| MODIFY | `ai/routing/interface.go`                                          |
| MODIFY | `ai/routing/service.go`（确保 `Service` 仍然满足所有拆分后的接口） |
| MODIFY | 所有使用 `RouterService` 的消费方（缩小依赖的接口类型）            |

## 改造内容

```diff
 // routing/interface.go

+// IntentClassifier 仅负责意图分类
+type IntentClassifier interface {
+    ClassifyIntent(ctx context.Context, input string) (Intent, float32, bool, error)
+}

+// ModelSelector 独立的模型选择接口
+type ModelSelector interface {
+    SelectModel(ctx context.Context, task TaskType) (ModelConfig, error)
+}

+// FeedbackCollector 反馈收集接口
+type FeedbackCollector interface {
+    RecordFeedback(ctx context.Context, feedback *RouterFeedback) error
+    GetRouterStats(ctx context.Context, userID int32, timeRange time.Duration) (*RouterStats, error)
+}

 // RouterService 保持向后兼容的聚合接口
 type RouterService interface {
+    IntentClassifier
+    ModelSelector
+    FeedbackCollector
-    ClassifyIntent(...)
-    SelectModel(...)
-    RecordFeedback(...)
-    GetRouterStats(...)
 }
```

## 验收条件

- [ ] `routing/interface.go` 中存在 `IntentClassifier`、`ModelSelector`、`FeedbackCollector` 三个独立接口
- [ ] `RouterService` 通过组合（embedding）包含上述三个接口
- [ ] `routing.Service` 结构体同时实现所有三个子接口
- [ ] 至少有 1 个消费方已改用更窄的接口类型（如 `IntentClassifier`）
- [ ] `go build ./...` 编译通过
- [ ] `go test ./ai/routing/... -count=1` 全部通过
- [ ] `go vet ./ai/...` 无警告
