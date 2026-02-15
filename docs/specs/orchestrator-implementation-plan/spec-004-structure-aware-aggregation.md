# SPEC-004: 结构感知聚合

> 优先级: P1 | 阶段: 阶段二 | 状态: 设计中

## 概述

Aggregator 根据来源 Agent 类型格式化最终输出，实现结构感知的响应渲染。

## 详细设计

### Agent 类型与渲染策略

| Agent 类型 | 渲染格式 | 示例 |
|:-----------|:---------|:-----|
| Schedule | Markdown 表格 | 包含时间、地点、参与者列 |
| Memo | 带引用的列表 | 包含 `[UID]` 或标题标注来源 |
| Geek | 原始输出 | 直接返回 CLI 输出 |

### 数据结构

```go
// Aggregator 聚合器
type Aggregator struct {
    agentRenderers map[string]AgentRenderer // agentType -> renderer
}

// AgentRenderer Agent 特定的渲染器接口
type AgentRenderer interface {
    Render(results []TaskResult) string
}

// ScheduleRenderer 日程渲染器
type ScheduleRenderer struct{}

func (sr *ScheduleRenderer) Render(results []TaskResult) string {
    // 渲染为 Markdown 表格
    var buf strings.Builder
    buf.WriteString("| 时间 | 地点 | 主题 | 参与者 |\n")
    buf.WriteString("|------|------|------|--------|\n")

    for _, r := range results {
        event := parseScheduleEvent(r.Result)
        buf.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
            event.Time, event.Location, event.Title, event.Participants))
    }
    return buf.String()
}

// MemoRenderer 笔记渲染器
type MemoRenderer struct{}

func (mr *MemoRenderer) Render(results []TaskResult) string {
    var buf strings.Builder
    for _, r := range results {
        note := parseMemoResult(r.Result)
        buf.WriteString(fmt.Sprintf("- [%s] %s\n  > %s\n",
            note.UID, note.Title, note.Summary))
    }
    return buf.String()
}
```

## 验收标准

- [ ] Schedule 来源渲染为 Markdown 表格
- [ ] Memo 来源渲染为带引用的列表
- [ ] 多类型混合时按 Agent 分组渲染

## 实现提示

1. **注册机制**: 在 Aggregator 初始化时注册各 Agent 的 Renderer
2. **Fallback**: 未知类型使用默认渲染 (直接输出)

## 依赖

- 前置: 无
- 后置: 无
