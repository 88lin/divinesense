# P1-04: jsonSchema 去重

> **阶段**: Phase 1 — 清理桥接层 & DRY 治理  
> **原则**: DRY  
> **风险**: 🟢 低  
> **预计工作量**: 0.5 天  
> **前置依赖**: 无

## 背景

`jsonSchema` 结构体在两处独立定义，代码完全相同（字段、JSON 标签、`MarshalJSON` 方法均一致）：

- `ai/title_generator.go` L211-223
- `ai/agents/llm_intent_classifier.go` L286-298

## 目标

将 `jsonSchema` 提取到公共位置，消除复制粘贴重复。

## 涉及文件

| 操作   | 文件                                 |
| :----- | :----------------------------------- |
| NEW    | `ai/core/llm/schema.go`              |
| MODIFY | `ai/title_generator.go`              |
| MODIFY | `ai/agents/llm_intent_classifier.go` |

## 改造内容

1. **创建** `ai/core/llm/schema.go`：
   - 导出 `JSONSchema` 结构体（首字母大写）
   - 包含 `MarshalJSON` 方法
2. **修改** `ai/title_generator.go`：删除本地 `jsonSchema`，改用 `llm.JSONSchema`
3. **修改** `ai/agents/llm_intent_classifier.go`：删除本地 `jsonSchema`，改用 `llm.JSONSchema`

## 验收条件

- [ ] `ai/core/llm/schema.go` 存在且导出 `JSONSchema` 类型
- [ ] `ai/title_generator.go` 不再包含 `jsonSchema` 结构体定义
- [ ] `ai/agents/llm_intent_classifier.go` 不再包含 `jsonSchema` 结构体定义
- [ ] 全量搜索 `type jsonSchema struct` 仅返回 0 处结果
- [ ] `go build ./...` 编译通过
- [ ] `go test ./ai/... -count=1` 全部通过
