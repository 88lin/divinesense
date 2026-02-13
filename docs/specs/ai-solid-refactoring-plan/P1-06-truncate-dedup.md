# P1-06: truncate 工具函数统一

> **阶段**: Phase 1 — 清理桥接层 & DRY 治理  
> **原则**: DRY  
> **风险**: 🟢 低  
> **预计工作量**: 1 天  
> **前置依赖**: P1-03 (cc_runner.go 清理后 TruncateString 别名已移除)

## 背景

字符串截断函数在 **6 处**独立定义（含 1 处别名），逻辑等价：

| 函数                        | 位置                                   | 类型                   |
| :-------------------------- | :------------------------------------- | :--------------------- |
| `TruncateString(s, maxLen)` | `agents/runner/types.go` L8            | 源定义（导出）         |
| `TruncateString(s, maxLen)` | `agents/cc_runner.go` L250             | 别名（P1-03 清理）     |
| `truncateForLog(s, maxLen)` | `agents/llm_intent_classifier.go` L231 | 私有重复               |
| `truncate(s, maxLen)`       | `routing/utils.go` L7                  | 私有重复               |
| `truncateLog(s, maxLen)`    | `tags/layer3_llm.go` L154              | 私有重复               |
| `Truncate(content, maxLen)` | `duplicate/similarity.go` L105         | 导出重复（rune-aware） |

> [!NOTE]
> `context/priority.go` 中的 `truncateToTokens(content, maxTokens)` 按 token 截断，语义不同，不在合并范围内。

## 目标

统一为一个公共工具函数，消除 5 处重复（保留 1 处定义）。

## 涉及文件

| 操作          | 文件                                                                                                 |
| :------------ | :--------------------------------------------------------------------------------------------------- |
| MODIFY 或 NEW | 建议新建 `ai/internal/strutil/truncate.go` 或使用 `agents/runner/types.go` 中已有的 `TruncateString` |
| MODIFY        | `ai/agents/llm_intent_classifier.go`（删除 `truncateForLog`）                                        |
| MODIFY        | `ai/routing/utils.go`（删除 `truncate`）                                                             |
| MODIFY        | `ai/tags/layer3_llm.go`（删除 `truncateLog`）                                                        |
| MODIFY        | `ai/duplicate/similarity.go`（删除 `Truncate`，注意此处是 rune-aware）                               |

## 改造内容

1. **审计**差异：`duplicate/similarity.go` 的 `Truncate` 使用 `[]rune` 截断（Unicode 安全），其他使用 `len(s)` 截断。统一实现应采用 rune-aware 版本
2. **确认**保留位置并统一签名
3. **逐文件替换**所有调用方
4. **删除**已不再使用的本地函数定义

## 验收条件

- [ ] 全量搜索 `func truncate\|func truncateForLog\|func truncateLog\|func Truncate` 在 `ai/` 下仅返回 1 处结果（统一定义处）
- [ ] 统一实现采用 rune-aware 截断（Unicode 安全）
- [ ] `go build ./...` 编译通过
- [ ] `go test ./ai/... -count=1` 全部通过

## 注意事项

- `duplicate/similarity.go` 的 `Truncate` 是导出函数，若有外部包引用需同步更新
- 统一后的函数应使用 `[]rune` 确保中文等多字节字符正确截断
