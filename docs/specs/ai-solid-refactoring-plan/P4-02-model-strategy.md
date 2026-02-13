# P4-02: 模型选择策略化

> **阶段**: Phase 4 — OCP 治理  
> **原则**: OCP (开闭原则)  
> **风险**: 🟡 中  
> **预计工作量**: 1 天  
> **前置依赖**: P3-04 (RouterService 接口拆分，ModelSelector 已独立)

## 背景

`routing/service.go` 中的 `SelectModel` 方法通过 switch 硬编码模型选择策略：

```go
func (s *Service) SelectModel(ctx context.Context, task TaskType) (ModelConfig, error) {
    switch task {
    case TaskTypeSimple: return s.lightModel, nil
    case TaskTypeComplex: return s.heavyModel, nil
    // ... 更多 case
    }
}
```

新增任务类型需修改源码，违反 OCP。

## 目标

引入 `ModelStrategy` 接口，使模型选择策略可配置、可扩展。

## 涉及文件

| 操作   | 文件                                                    |
| :----- | :------------------------------------------------------ |
| NEW    | `ai/routing/model_strategy.go`                          |
| MODIFY | `ai/routing/service.go`（`SelectModel` 委托给策略实现） |

## 改造内容

### Step 1: 定义策略接口

```go
// routing/model_strategy.go [NEW]
package routing

// ModelStrategy 定义模型选择策略
type ModelStrategy interface {
    SelectModel(ctx context.Context, task TaskType) (ModelConfig, error)
}

// ConfigDrivenModelStrategy 基于配置映射的模型选择策略
type ConfigDrivenModelStrategy struct {
    configs  map[TaskType]ModelConfig
    fallback ModelConfig
}

func NewConfigDrivenModelStrategy(configs map[TaskType]ModelConfig, fallback ModelConfig) *ConfigDrivenModelStrategy {
    return &ConfigDrivenModelStrategy{configs: configs, fallback: fallback}
}

func (s *ConfigDrivenModelStrategy) SelectModel(ctx context.Context, task TaskType) (ModelConfig, error) {
    if cfg, ok := s.configs[task]; ok {
        return cfg, nil
    }
    return s.fallback, nil
}
```

### Step 2: Service 委托给策略

```diff
 type Service struct {
     ...
-    lightModel ModelConfig
-    heavyModel ModelConfig
+    modelStrategy ModelStrategy
 }

 func (s *Service) SelectModel(ctx context.Context, task TaskType) (ModelConfig, error) {
-    switch task { ... }
+    return s.modelStrategy.SelectModel(ctx, task)
 }
```

## 验收条件

- [ ] `routing/model_strategy.go` 存在且导出 `ModelStrategy` 接口
- [ ] `ConfigDrivenModelStrategy` 实现 `ModelStrategy` 接口
- [ ] `routing.Service.SelectModel` 委托给 `ModelStrategy` 实现
- [ ] `SelectModel` 方法内不再包含 switch/case 硬编码
- [ ] 可通过配置新增任务类型的模型映射，无需修改源码
- [ ] `go build ./...` 编译通过
- [ ] `go test ./ai/routing/... -count=1` 全部通过
