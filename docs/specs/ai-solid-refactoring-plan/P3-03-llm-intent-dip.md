# P3-03: LLMIntentClassifier DIP 重构

> **阶段**: Phase 3 — DIP / ISP 治理  
> **原则**: DIP (依赖倒置)  
> **风险**: 🟢 低  
> **预计工作量**: 0.5 天  
> **前置依赖**: P1-04 (jsonSchema 去重完成后)

## 背景

`ai/agents/llm_intent_classifier.go` 直接实例化 `openai.Client`，绕过了 `core/llm.Service` 抽象层。与 P1-02 (title_generator) 属于同类问题。

## 目标

将 `LLMIntentClassifier` 的 LLM 依赖从具体 `openai.Client` 改为抽象 `LLMService` 接口注入。

## 涉及文件

| 操作   | 文件                                     |
| :----- | :--------------------------------------- |
| MODIFY | `ai/agents/llm_intent_classifier.go`     |
| MODIFY | 所有调用 `NewLLMIntentClassifier` 的文件 |

## 改造内容

```diff
 type LLMIntentClassifier struct {
-    client   *openai.Client
-    model    string
+    llm      ai.LLMService
 }

-func NewLLMIntentClassifier(cfg LLMIntentConfig) *LLMIntentClassifier {
-    client := openai.NewClientWithConfig(...)
-    return &LLMIntentClassifier{client: client, model: cfg.Model}
+func NewLLMIntentClassifier(llm ai.LLMService) *LLMIntentClassifier {
+    return &LLMIntentClassifier{llm: llm}
 }

 func (c *LLMIntentClassifier) ClassifyWithDetails(...) (*IntentResult, error) {
-    resp, err := c.client.CreateChatCompletion(ctx, req)
+    resp, _, err := c.llm.Chat(ctx, messages)
 }
```

## 验收条件

- [ ] `ai/agents/llm_intent_classifier.go` 不再 import `github.com/sashabaranov/go-openai`
- [ ] `LLMIntentClassifier` 通过构造函数接收 `LLMService` 接口
- [ ] `LLMIntentConfig` 结构体已废弃或移除
- [ ] `go build ./...` 编译通过
- [ ] `go test ./ai/agents/... -count=1` 全部通过
- [ ] 可通过 mock `LLMService` 对 `LLMIntentClassifier` 进行单元测试
