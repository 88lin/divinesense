# SPEC-002: 变量替换与上下文注入

> 优先级: P0 | 阶段: 阶段二 | 状态: 设计中

## 概述

实现任务输入的动态变量替换，支持从上游任务结果或全局上下文提取数据。语法: `{{Source.Field | Modifier}}`

## 详细设计

### 核心逻辑

**支持的语法**:
| 语法 | 说明 | 示例 |
|:-----|:-----|:-----|
| `{{task_id.result}}` | 引用上游任务结果 | `{{t1.result}}` |
| `{{task_id.result \| summary}}` | 引用并自动摘要 | `{{t1.result \| summary}}` |
| `{{global.time}}` | 引用当前时间上下文 | `{{global.time}}` |
| `{{global.date}}` | 引用当前日期 | `{{global.date}}` |

**处理流程**:
1. 解析输入字符串，匹配 `{{...}}` 模式
2. 识别 Source 类型 (task / global)
3. 从上下文存储中获取对应值
4. 如果值超过阈值，触发 Summary 修饰器
5. 替换原位置并返回

### 数据结构

```go
// ContextInjector 上下文注入器
type ContextInjector struct {
    taskResults map[string]string           // taskID -> result
    globalCtx    map[string]interface{}     // 全局上下文
    maxTokens    int                        // 最大 token 阈值
    llmClient    LLMClient                  // 用于 Summary
}

// VariableRef 变量引用
type VariableRef struct {
    Source    string // "task" or "global"
    Field     string // "result", "time", "date"
    Modifier  string // "summary", empty
    Raw       string // 原始字符串
}
```

### Token 安全策略

```go
const (
    MaxInputTokens   = 2000  // 单次注入阈值
    WarnThresholdKB  = 100    // 超过 100KB 记录警告
)

// 替换前检查
func (ci *ContextInjector) safeReplace(input string, taskResults map[string]string) (string, error) {
    // 1. 解析所有变量引用
    refs := ci.parseVariables(input)

    // 2. 估算替换后 token 数
    estimated := ci.estimateTokens(input, taskResults)

    // 3. 超过阈值则触发 Summary
    if estimated > MaxInputTokens {
        return ci.replaceWithSummary(input, taskResults)
    }

    return ci.replace(input, taskResults)
}
```

## 验收标准

- [ ] `Input: "分析 {{t1.result}}"` 正确替换为上游任务结果
- [ ] 超长结果自动触发 Summary (调用 LLM 压缩)
- [ ] `{{global.time}}` 返回格式化的时间字符串
- [ ] 引用不存在的 task 返回错误
- [ ] 替换前后记录 TokenUsage

## 实现提示

1. **Regex**: 使用 `\{\{\s*([a-zA-Z0-9_\-]+)\.(result|time|date)\s*(\|\s*summary)?\}\}`
2. **Summary**: 复用现有的 LLM Client，构造摘要 Prompt
3. **Trace**: 在注入前后记录 `TokenUsage`，用于可观测性

## 依赖

- 前置: SPEC-001 (DAG 调度)
- 后置: SPEC-009 (Executor 升级)
