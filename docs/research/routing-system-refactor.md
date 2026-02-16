# 路由系统重构调研报告

> 调研时间: 2026-02-16 | 版本: v1.0 | 状态: 已创建 Issue #241

## 背景

在 E2E 测试中发现路由问题：
- 用户输入 "查看今天的笔记" 被错误路由到 SCHEDULE（置信度 0.85）
- Handoff 失败，提示 "笔记管理专家 is not a valid expert"

## 问题根因

### 问题 1: Fast Path 绕过权重机制

**位置**: `ai/routing/rule_matcher.go:465-468`

```go
// FAST PATH: Time pattern + query pattern → schedule query
if m.hasTimePattern(input) && queryPatternRegex.MatchString(lower) && !m.hasMemoKeyword(input) {
    return IntentScheduleQuery, 0.85, true
}
```

- "查看今天的笔记" 匹配 "今天"（时间词）+ "查看"（查询词）
- 但 `hasMemoKeyword()` 返回 false（因为 "笔记" 不在 capability_triggers 中）
- 导致直接返回 schedule，权重计算被完全绕过

### 问题 2: Handoff Expert 验证失败

**位置**: `ai/agents/tools/report_inability.go:12-15`

```go
var ValidExpertAgents = map[string]bool{
    "memo":     true,
    "schedule": true,
}
```

- LLM 返回中文 title（如 "笔记搜索专家"）
- 白名单只接受英文 name（如 "memo"）
- 验证失败，Handoff 无法执行

## 解决方案

### 决策

| 问题 | 决策 |
|------|------|
| 配置兼容性 | 强制式 - 所有配置文件必须包含 routing 字段 |
| 关键字匹配 | 词边界匹配 - 使用正则 `\bkeyword\b` |
| 无匹配行为 | 转 Orchestrator - 复杂请求由 Orchestrator 处理 |
| 学习系统 | 长期记忆存储 - 用户偏好存数据库，不修改配置文件 |
| 热更新 | 信号触发 - 收到 SIGHUP 时重载配置 |

### 技术架构

```
config/parrots/*.yaml (强制 routing 字段)
        ↓
AgentRegistry (运行时配置加载器)
        ↓
RouterService (配置驱动的路由)
        ├─ Layer 1: Cache (LRU)
        ├─ Layer 2: Keyword Scoring (配置驱动)
        ├─ Layer 3: HITL (冲突解决)
        └─ Layer 4: Sticky + Learning
```

### 核心数据结构

```go
type AgentConfig struct {
    Name            string          `yaml:"name"`
    DisplayName     string          `yaml:"display_name"`
    RoutingKeywords []RoutingKeyword `yaml:"routing.keywords"`
    WeightConfig    *AgentWeightConfig `yaml:"routing.weight_config"`
    MutexGroups    []MutexGroup    `yaml:"routing.mutex_groups"`
}

type RoutingKeyword struct {
    Keyword  string `yaml:"keyword"`  // "笔记", "note", "memo"
    Category string `yaml:"category"` // "memo", "schedule"
    Weight   int    `yaml:"weight"`   // 1-10, 权重越高越优先
}

type AgentWeightConfig struct {
    BaseScore    int `yaml:"base_score"`
    KeywordBonus int `yaml:"keyword_bonus"`   // 每匹配一个关键字 +2
    TimewordBonus int `yaml:"timeword_bonus"` // 每匹配一个时间词 +1
}
```

## 验收标准

- [ ] 零硬编码：所有路由关键字从配置文件读取
- [ ] 新增 Agent：只需添加 YAML 配置，无需修改代码
- [ ] 权重可配置：关键字权重在配置文件中设置
- [ ] HITL：冲突时请求用户确认
- [ ] 粘性路由：基于历史会话记住用户偏好
- [ ] 性能：路由决策 < 1ms（缓存命中）
- [ ] `make check-all` 通过

## 相关 Issue

- #230: 查询"查看今天的笔记"被错误路由到SCHEDULE且未触发Handoff
- #241: 路由系统全面重构 - 零硬编码全配置化

## 参考资料

- Lobe Chat: 配置化 Agent 注册机制
- ChatBotKit: 多 Provider Fallback
- LangChain: Tool + Agent 描述驱动
