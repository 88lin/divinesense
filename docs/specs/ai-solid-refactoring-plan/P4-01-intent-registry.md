# P4-01: Intent 注册表机制

> **阶段**: Phase 4 — OCP 治理  
> **原则**: OCP (开闭原则)  
> **风险**: 🟡 中  
> **预计工作量**: 2 天  
> **前置依赖**: P3-04 (RouterService 接口拆分), P3-01 (ConversationContext 解耦)

## 背景

当前新增 Agent 类型需联动修改多处硬编码：

1. `routing/interface.go` — 添加 `Intent` 枚举常量
2. `routing/interface.go` — 更新 `IntentToAgentType` 映射
3. `routing/interface.go` — 更新 `AgentTypeToIntent` 映射
4. `agents/chat_router.go` — 更新 `mapIntentToRouteType` switch
5. `agents/intent_classifier.go` — 添加匹配关键词

至少 **5 处联动修改**，严重违反 OCP。

## 目标

引入 `IntentRegistry` 注册表机制，使新增 Agent 类型仅需在启动时注册配置，不修改现有源码。

## 涉及文件

| 操作   | 文件                                                                                    |
| :----- | :-------------------------------------------------------------------------------------- |
| NEW    | `ai/routing/intent_registry.go`                                                         |
| MODIFY | `ai/routing/interface.go`（`IntentToAgentType` / `AgentTypeToIntent` 改为从注册表读取） |
| MODIFY | `ai/routing/service.go`（`ClassifyIntent` 使用注册表匹配）                              |
| MODIFY | `ai/agents/chat_router.go`（`mapIntentToRouteType` 从注册表读取）                       |
| MODIFY | 启动入口（注册内置 Intent）                                                             |

## 改造内容

### Step 1: 定义注册表

```go
// routing/intent_registry.go [NEW]
type IntentRegistry struct {
    mu       sync.RWMutex
    intents  map[Intent]IntentConfig
    mappings map[Intent]AgentType
}

type IntentConfig struct {
    Intent      Intent
    AgentType   AgentType
    Keywords    []string          // 快速匹配关键词
    Patterns    []*regexp.Regexp  // 正则匹配模式
    Priority    int               // 匹配优先级
    RouteType   string            // 对应的路由类型
}

func NewIntentRegistry() *IntentRegistry { ... }
func (r *IntentRegistry) Register(cfg IntentConfig) { ... }
func (r *IntentRegistry) Match(input string) (Intent, float32, bool) { ... }
func (r *IntentRegistry) GetAgentType(intent Intent) (AgentType, bool) { ... }
func (r *IntentRegistry) GetRouteType(intent Intent) (string, bool) { ... }
```

### Step 2: 启动时注册内置 Intent

```go
registry := routing.NewIntentRegistry()
registry.Register(IntentConfig{
    Intent:    IntentScheduleCreate,
    AgentType: AgentTypeScheduler,
    Keywords:  []string{"安排", "预约", "创建日程"},
    RouteType: "schedule",
    Priority:  100,
})
// ... 注册其他内置 intent
```

### Step 3: 消费方改为查询注册表

```diff
 // routing/service.go
 func (s *Service) ClassifyIntent(ctx, input) (Intent, float32, bool, error) {
-    // 硬编码规则匹配
-    if containsKeyword(input, scheduleKeywords) { ... }
+    intent, conf, found := s.registry.Match(input)
+    if found { return intent, conf, false, nil }
     // 降级到 LLM
 }
```

## 验收条件

- [ ] `routing/intent_registry.go` 存在且导出 `IntentRegistry` 和 `IntentConfig`
- [ ] `IntentToAgentType` 和 `AgentTypeToIntent` 映射由注册表生成，不再硬编码
- [ ] `chat_router.go` 中的 `mapIntentToRouteType` 从注册表读取
- [ ] `intent_classifier.go` 中的关键词列表从注册表读取
- [ ] 新增一个测试 Intent 仅通过 `Register` 即可被匹配到（无需修改源码）
- [ ] `go build ./...` 编译通过
- [ ] `go test ./ai/... -count=1` 全部通过

## 注意事项

- 现有 `Intent` 常量枚举可保留作为内置默认值，仅将映射逻辑委托给注册表
- 注册表初始化应在 `routing.NewService` 内完成，保证向后兼容
