# P1-02: title_generator.go DIP 重构

> **阶段**: Phase 1 — 清理桥接层 & DRY 治理  
> **原则**: DIP (依赖倒置)  
> **风险**: 🟢 低  
> **预计工作量**: 0.5 天  
> **前置依赖**: 无

## 背景

`ai/title_generator.go` 直接实例化 `openai.Client`，绕过了 `core/llm.Service` 抽象层。这违反了 DIP，且使该组件无法在不引入 OpenAI SDK 的情况下被测试。

## 目标

将 `TitleGenerator` 的 LLM 依赖从具体 `openai.Client` 改为抽象 `LLMService` 接口注入。

## 涉及文件

| 操作   | 文件                                                |
| :----- | :-------------------------------------------------- |
| MODIFY | `ai/title_generator.go`                             |
| MODIFY | 所有调用 `NewTitleGenerator` 的文件（更新构造参数） |

## 改造内容

1. **移除** `TitleGenerator` 中的 `client *openai.Client` 和 `model string` 字段
2. **替换为** `llm LLMService` 字段（使用 `ai` 包内已定义的接口）
3. **修改** `NewTitleGenerator` 签名为 `NewTitleGenerator(llm LLMService) *TitleGenerator`
4. **修改** `Generate` 方法内部，由 `client.CreateChatCompletion` 改为 `llm.Chat`
5. **删除** 对 `github.com/sashabaranov/go-openai` 的直接 import
6. **更新** 所有调用方传入已有的 `LLMService` 实例

## 验收条件

- [ ] `ai/title_generator.go` 不再 import `github.com/sashabaranov/go-openai`
- [ ] `TitleGenerator` 通过构造函数接收 `LLMService` 接口
- [ ] `TitleGeneratorConfig` 结构体中的 `APIKey`/`BaseURL`/`Model` 字段已废弃或移除
- [ ] `go build ./...` 编译通过
- [ ] `go test ./ai/... -count=1` 全部通过
- [ ] 可通过 mock `LLMService` 对 `TitleGenerator` 进行单元测试
