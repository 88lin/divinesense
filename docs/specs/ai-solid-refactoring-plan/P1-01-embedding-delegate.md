# P1-01: embedding.go 委托重构

> **阶段**: Phase 1 — 清理桥接层 & DRY 治理  
> **原则**: DIP (依赖倒置)  
> **风险**: 🟢 低  
> **预计工作量**: 0.5 天

## 背景

`ai/embedding.go` 同时定义了 `EmbeddingService` 接口 **和** `embeddingService` 具体实现，直接依赖 `go-openai` 库。而 `ai/core/embedding` 已有完整的抽象实现。根包不应持有具体实现。

## 目标

将 `ai/embedding.go` 改为纯委托层，所有实现逻辑委托到 `ai/core/embedding` 包。

## 涉及文件

| 操作   | 文件                                                   |
| :----- | :----------------------------------------------------- |
| MODIFY | `ai/embedding.go`                                      |
| 验证   | 所有 `import "...ai"` 且使用 `EmbeddingService` 的文件 |

## 改造内容

1. **删除** `embeddingService` 具体结构体及其所有方法实现
2. **删除** 对 `go-openai` 的直接导入
3. **将** `EmbeddingService` 改为 `embedding.Service` 的类型别名（deprecated）
4. **将** `NewEmbeddingService` 改为构造函数委托，调用 `embedding.NewService`
5. **确认** `EmbeddingConfig` 与 `embedding.Config` 字段一致性，必要时对齐

## 验收条件

- [ ] `ai/embedding.go` 不再包含任何 `func (s *embeddingService)` 方法
- [ ] `ai/embedding.go` 不再 import `github.com/sashabaranov/go-openai`
- [ ] `go build ./...` 编译通过
- [ ] `go test ./ai/... -count=1` 全部通过
- [ ] 所有现有调用 `ai.NewEmbeddingService` 的代码无需修改即可正常工作

## 回滚方案

纯桥接重构，若出现问题直接 `git revert` 即可。
