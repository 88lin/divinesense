# 重构方案：消除骨架层硬编码，符合 SOLID 原则

## 背景

当前代码存在多处硬编码，违反 SOLID 原则，新增专家 Agent 时需要修改多处代码。

## 问题分析

### 1. 硬编码点汇总

| 位置 | 硬编码内容 | 影响 |
|------|------------|------|
| `routing/rule_matcher.go:13-67` | `memo`/`schedule` 关键词 | 新专家无法自动路由 |
| `routing/rule_matcher.go:83-87` | 时间模式快速路径硬编码 | Issue #230 问题根因 |
| `universal/react_executor.go:249-257` | `schedule` 成功指示器 | 新专家无法触发 early stop |

### 2. 架构分析

**现有可复用组件**：
- `CapabilityMap`: 动态能力→专家映射（已实现）
- `IntentRegistry`: 可配置的意图注册表（部分使用）
- `ExpertRegistry`: 动态获取专家配置
- `ParrotSelfCognition`: 专家自描述配置

**依赖方向（当前）**：
```
RuleMatcher ← 硬编码 keywords
react_executor ← 硬编码 successIndicators
```

**期望依赖方向**：
```
RuleMatcher ← ExpertRegistry ← config/*.yaml
react_executor ← ExpertRegistry ← config/*.yaml
```

## 架构决策

| 决策项 | 选择 | 说明 |
|--------|------|------|
| 路由关键词来源 | CapabilityTriggers | 复用现有字段 |
| successIndicators | 自动推导 | 不需要配置 |
| 路由匹配机制 | RuleMatcher 动态加载 | 从 CapabilityTriggers 加载 |

## 重构方案

### 方案 1: 添加 routing_keywords 配置（推荐）

在专家配置中扩展 `self_description.capability_triggers` 字段：

```yaml
# config/parrots/schedule.yaml
self_description:
  name: schedule
  capabilities:
    - "创建日程"
    - "查询日程"
  capability_triggers:
    "查询日程": ["查看", "日程", "安排", "有什么"]
    "创建日程": ["创建", "添加", "安排"]

# config/parrots/memo.yaml
self_description:
  name: memo
  capabilities:
    - "搜索笔记"
  capability_triggers:
    "搜索笔记": ["笔记", "搜索", "查找", "memo"]
```

#### 修改点

**Phase 1: 配置层**

1. 扩展 `ParrotSelfCognition` 结构（如果需要）

**Phase 2: RuleMatcher 重构**

2. 修改 `RuleMatcher.Match()` 从 CapabilityMap 获取关键词
3. 移除硬编码 `coreKeywordsByCategory`
4. 移除硬编码 `scheduleKeywords` / `memoKeywords`

**Phase 3: 快速路径重构**

5. 修改时间模式快速路径：从 CapabilityMap 判断是否有日程相关 capability
6. 或完全移除快速路径，统一走评分机制

**Phase 4: Executor 自动推导**

7. 修改 `react_executor.go`:
   - 工具执行成功后，检查返回 JSON 中的 `success` 字段
   - 或检查特定工具名称的成功模式（如 `*_add`, `*_create` → 成功）
   - 不再硬编码 successIndicators

### 方案 2: 完全动态化（备选）

移除所有硬编码，完全依赖配置驱动：

1. 路由层：完全从 IntentRegistry + YAML 配置加载
2. Executor：从工具元数据自动推导成功模式

## 实现步骤

### Step 1: 扩展配置（1-2 文件）

```go
// ai/agents/base_parrot.go - 如果需要
type ParrotSelfCognition struct {
    // 现有字段...
    RoutingKeywords map[string][]string `json:"routing_keywords,omitempty"`
}
```

### Step 2: RuleMatcher 动态化（核心改动）

```go
// ai/routing/rule_matcher.go
type RuleMatcher struct {
    // 替换硬编码
    // 原: scheduleKeywords map[string]int
    // 新: 从 ExpertRegistry 动态获取
    capabilityMap *CapabilityMap // 注入
}
```

### Step 3: 移除快速路径硬编码

```go
// ai/routing/rule_matcher.go
func (m *RuleMatcher) Match(input string) (Intent, float32, bool) {
    // 移除硬编码的 "笔记" 检查
    // 改为: 从 capabilityMap 查询是否有任何专家包含"笔记" capability
}
```

### Step 4: Executor 自动推导

```go
// ai/agents/universal/react_executor.go
func shouldEarlyStop(toolResult string) bool {
    // 原: 硬编码 "✓ 已创建", "schedule created"
    // 新: 检查返回 JSON 中的 success 字段
    // 或检查工具名模式: *_add, *_create 成功 = true
}
```

## 影响范围

| 文件 | 改动类型 | 风险 |
|------|----------|------|
| `routing/rule_matcher.go` | 重构 | 高 |
| `universal/react_executor.go` | 重构 | 中 |
| `config/parrots/*.yaml` | 配置 | 低 |
| `agents/base_parrot.go` | 扩展 | 低 |

## 测试策略

1. **现有测试**: 确保重构后行为一致
2. **新增专家测试**: 新增一个"测试专家"验证动态加载
3. **边界测试**: 空配置、降级方案

## 风险与回退

- **风险**: RuleMatcher 性能可能下降（需要 Benchmark）
- **回退**: 保留旧的硬编码路径作为降级方案

## 待确认

- [ ] CapabilityTriggers 字段是否需要重命名为 routing_keywords？
- [ ] 是否需要保留硬编码作为降级方案？
- [ ] 自动推导 successIndicators 的具体实现方式？
