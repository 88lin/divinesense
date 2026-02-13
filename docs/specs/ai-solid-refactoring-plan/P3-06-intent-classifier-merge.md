# P3-06: IntentClassifier 双轨合并

> **阶段**: Phase 3 — DIP / ISP 治理  
> **原则**: DRY + SRP  
> **风险**: 🟡 中  
> **预计工作量**: 1.5 天  
> **前置依赖**: P3-03 (LLMIntentClassifier DIP 重构完成), P4-01 (IntentRegistry 可选前置)

## 背景

`agents/` 包内存在两套并行的意图分类器：

| 实现                         | 文件                              | 方法              |
| :--------------------------- | :-------------------------------- | :---------------- |
| 规则版 `IntentClassifier`    | `agents/intent_classifier.go`     | 关键词 + 正则匹配 |
| LLM 版 `LLMIntentClassifier` | `agents/llm_intent_classifier.go` | LLM 调用          |

两者的意图枚举 (`TaskIntent`) 完全一致，`ShouldUsePlanExecute` 等方法逻辑重复。而 `routing/` 包已实现了分层路由（cache → rule → LLM fallback），`agents/` 包内的旧实现形成冗余。

## 目标

将 `agents/` 包内的双轨 IntentClassifier 合并为单一实现，复用 `routing/` 的分层机制，消除 ~400 行重复逻辑。

## 涉及文件

| 操作                | 文件                                                            |
| :------------------ | :-------------------------------------------------------------- |
| DELETE 或 DEPRECATE | `ai/agents/intent_classifier.go`                                |
| MODIFY              | `ai/agents/llm_intent_classifier.go`（简化为 routing 的适配层） |
| MODIFY              | 所有调用 `IntentClassifier.Classify` 的消费方                   |

## 改造内容

### 方案 A（推荐）: 统一到 routing 分层架构

1. **删除** `agents/intent_classifier.go` 中的规则匹配逻辑
2. **将** `LLMIntentClassifier` 简化为 `routing.IntentClassifier` 接口的适配器
3. 消费方改为依赖 `routing.IntentClassifier` 接口
4. `TaskIntent` 枚举保留在 `agents/` 包中（被 routing 引用）

### 方案 B（保守）: 规则版降级为 LLM 版的前置快速路径

1. **将** `IntentClassifier` 嵌入 `LLMIntentClassifier` 作为快速路径
2. `LLMIntentClassifier.Classify` 先调用规则匹配，confidence 不足时才调 LLM
3. 删除重复的 `ShouldUsePlanExecute` 等公共方法

## 验收条件

- [ ] `agents/` 包中的意图分类入口仅一个（不再有两套并行的分类器）
- [ ] `TaskIntent` 枚举在代码库中仅定义一处
- [ ] `ShouldUsePlanExecute` 方法仅存在一处
- [ ] `go build ./...` 编译通过
- [ ] `go test ./ai/agents/... -count=1` 全部通过
- [ ] `go test ./ai/routing/... -count=1` 全部通过

## 注意事项

- 需先确认 `agents/intent_classifier.go` 的哪些消费方在非 routing 路径中独立调用
- 若存在不经过 routing 的直接调用，需保留一个轻量级适配层
