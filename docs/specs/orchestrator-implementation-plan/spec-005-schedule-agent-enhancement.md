# SPEC-005: Schedule Agent 增强

> 优先级: P0 | 阶段: 阶段一 | 状态: 待实现

## 概述

增强 Schedule Agent 的交互协议，防止过早调用工具，提高对模糊输入的鲁棒性，减少幻觉。

## 详细设计

### 核心变更

更新 `config/parrots/schedule.yaml` 的 `system_prompt`，加入严格的 "思考-行动" (Thought-Action) 协议和澄清触发器。

### 协议要求

```
## Execution Protocol (Strict Order)
1. <Analyze>: 分析用户意图。时间/时长是否明确？
   - IF 没有时间: 调用 `find_free_time` 或 `ask_user`
   - IF 修改日程: 必须先调用 `schedule_query` 找到目标事件
2. <Validation>: 检查逻辑冲突 (例如：凌晨 3 点开会)
3. <Execution>: 调用工具
4. <Reflection>: 评估工具输出
   - IF 冲突: 礼貌地提出替代方案。不要只说 "失败"

## Clarification Triggers (澄清触发器)
- IF 用户说 "安排会议" (无时间):
  Ask "请问您希望安排在具体哪天？或者通过 find_free_time 帮您查找合适的时间？"
- IF 用户说 "和他们开会" (无具体人):
  Ask "请问是和哪个团队或具体哪位同事？"
```

### 配置示例

```yaml
system_prompt: |
  ## Identity
  你是时巧 (ScheduleParrot)，一个专业的日程管理助手...

  ## Execution Protocol (Strict Order)
  1. <Analyze>: 分析用户意图...
  ...

  ## Clarification Triggers
  - IF 用户说 "安排会议" (无时间): Ask "请问您希望安排在具体哪天..."
  ...
```

## 验收标准

- [ ] 用户说 "安排个会" 时，Agent 主动询问时间
- [ ] 用户说 "和他们开会" 时，Agent 主动询问具体人员
- [ ] 修改日程前，Agent 先调用查询工具找到目标事件
- [ ] 逻辑冲突时 (如凌晨安排会议)，Agent 提出替代方案

## 实现提示

1. **文件位置**: `config/parrots/schedule.yaml`
2. **测试**: 创建模糊指令的测试 Case
3. **监控**: 观察澄清率是否提升

## 依赖

- 前置: 无
- 后置: 无
