# SPEC-007: Decomposer 专家感知

> 优先级: P0 | 阶段: 阶段二 | 状态: 设计中

## 概述

在 Decomposer 的 System Prompt 中明确列出当前可用的 Agents 及其能力描述，确保 LLM 生成准确的 `agent_id`。

## 详细设计

### 背景

由于 CapabilityMap 尚未实现，Decomposer 无法动态获取专家列表。需要在 Prompt 中硬编码当前可用的 Agents。

### Agent 能力定义

```yaml
# config/orchestrator/decomposer.yaml

system_prompt: |
  ## Available Expert Agents

  你可以调用以下专家代理完成用户任务：

  1. **memo (灰灰)** - 笔记搜索与管理
     - 能力: 语义搜索笔记、按时间浏览、笔记摘要
     - 适用: 查找信息、回忆记录、整理知识

  2. **schedule (时巧)** - 日程管理
     - 能力: 创建日程、查询空闲时间、修改/删除日程
     - 适用: 安排会议、设置提醒、管理时间

  3. **geek (极客)** - 外部工具执行
     - 能力: 执行 Claude Code CLI 命令、代码操作
     - 适用: 代码修改、文件操作、技术任务

  ## Task Output Format

  每个任务必须指定：
  - `id`: 唯一标识，如 "task_1", "task_2"
  - `agent_id`: 上述 agent 之一
  - `input`: 任务的输入内容
  - `dependencies`: 依赖的任务 ID 列表 (可选)
```

### 时间上下文注入

```go
// Decompose 时注入时间上下文
func (d *Decomposer) Decompose(ctx context.Context, userInput string) (*TaskPlan, error) {
    timeCtx := universal.BuildTimeContext(time.Now())

    prompt := d.promptTemplate.Execute(struct {
        UserInput   string
        TimeContext string
        Agents      string
    }{
        UserInput:   userInput,
        TimeContext: timeCtx,
        Agents:      d.agentDescriptions,
    })

    // 调用 LLM
    return d.llmClient.Generate(ctx, prompt)
}
```

## 验收标准

- [ ] Decomposer 生成的任务包含正确的 `agent_id`
- [ ] 支持识别 "搜索笔记" 使用 memo，"安排会议" 使用 schedule
- [ ] 相对时间解析正确 (如下周二 → 具体日期)
- [ ] 任务 ID 在 dependencies 中正确引用

## 实现提示

1. **文件位置**: `config/orchestrator/decomposer.yaml`
2. **同步**: 当新增 Agent 时，同步更新此处的 Agent 列表
3. **测试**: 验证不同任务类型生成的 agent_id 是否正确

## 依赖

- 前置: 无
- 后置: SPEC-001 (DAG 调度)
