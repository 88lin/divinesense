# AI 重构 Specs 独立审计报告

> **审计日期**: 2026-02-13  
> **审计方法**: 以源码为 ground truth，逐条验证架构方案和 Spec 中的每项声明  
> **审计结论**: 🟡 **有条件通过** — 发现 3 类问题需修正

---

## 审计维度一：Spec 中的事实性错误

以下是 Spec 内容与真实源码不符的地方：

### 🔴 F-1: P1-06 truncate 数量低估

**Spec 声明**: truncate 函数存在于 3 处  
**实际源码**: truncate 类函数存在于 **7 处**

| 函数                                   | 位置                                   | Spec 已覆盖                 |
| :------------------------------------- | :------------------------------------- | :-------------------------- |
| `TruncateString(s, maxLen)`            | `agents/runner/types.go` L8            | ✅                           |
| `TruncateString(s, maxLen)`            | `agents/cc_runner.go` L250（别名）     | ✅                           |
| `truncateForLog(s, maxLen)`            | `agents/llm_intent_classifier.go` L231 | ✅                           |
| `truncate(s, maxLen)`                  | `routing/utils.go` L7                  | ✅                           |
| `truncateLog(s, maxLen)`               | **`tags/layer3_llm.go` L154**          | ❌ **遗漏**                  |
| `Truncate(content, maxLen)`            | **`duplicate/similarity.go` L105**     | ❌ **遗漏**                  |
| `truncateToTokens(content, maxTokens)` | `context/priority.go` L89              | ⚠️ 语义不同（按 token 截断） |

> **影响**: P1-06 完成后仍会残留 2 处重复。需扩展 Spec 扫描范围。

---

### 🟡 F-2: P3-05 SafeCallbackFunc 签名描述不准确

**架构方案声明**: `base_parrot.go` 中 `SafeCallbackFunc` 签名为 `func(string, interface{}) error`  
**实际源码** (`runner/runner.go` L68):

```go
type SafeCallbackFunc func(eventType string, eventData any)
```

关键差异：
1. 定义在 **`runner/runner.go`** 而非 `base_parrot.go`（`base_parrot.go` 只使用不定义）
2. 返回值为 **无返回值** 而非 `error`
3. `cc_runner.go` L148 中为类型别名 `type SafeCallbackFunc = runner.SafeCallbackFunc`

> **影响**: P3-05 中的适配方案需据此调整。EventCallback (有 error 返回) 和 SafeCallbackFunc (无 error 返回) 是两种不同层级的回调——SafeCallback 本身是 EventCallback 的包装器。

---

### 🟡 F-3: P3-07 ParrotAgent 接口描述不准确

**Spec 声明**: `Execute(ctx, input)` 与 `ExecuteWithCallback(ctx, input, callback)` 重叠  
**实际源码** (`base_parrot.go` L326-333):

```go
type ParrotAgent interface {
    Name() string
    Execute(ctx context.Context, userInput string, callback EventCallback) error
    ExecuteWithCallback(ctx context.Context, userInput string, history []string, callback EventCallback) error
    SelfDescribe() *ParrotSelfCognition
}
```

关键差异：
1. `Execute` **已包含** callback 参数
2. 两者真正的区别是 `history []string` 参数 — `Execute` 是 `ExecuteWithCallback` 的无历史简写
3. 所有实现（universal、geek、evolution）的 `Execute` 都直接调用 `ExecuteWithCallback(ctx, input, nil, callback)`

> **影响**: P3-07 的改造方案仍然正确（合并为单方法），但 rationale 需修正 — 不是 callback 重叠，而是 history 参数可选化。

---

## 审计维度二：架构方案遗漏的违规项

以下是源码中存在、但架构方案和 Spec 均未提及的违规：

### 🔴 M-1: `store` 包的广泛直接依赖（DIP 系统性违规）

`ai/` 子树中有 **17 个文件**直接导入 `github.com/hrygo/divinesense/store`：

| 子包                     | 文件                                        |
| :----------------------- | :------------------------------------------ |
| `agents/`                | `context.go`, `error_class.go` (已覆盖)     |
| `agents/geek/`           | `mode.go`, `evolution.go`                   |
| `services/stats/`        | `alerting.go`, `persister.go`               |
| `services/schedule/`     | `precheck_service.go`                       |
| `routing/`               | `postgres_storage.go`                       |
| `tags/`                  | `layer1_statistics.go`, `suggester_impl.go` |
| `duplicate/`             | `detector_impl.go`                          |
| `review/`                | `service.go`                                |
| `observability/metrics/` | `persister.go`, `service.go`                |
| `graph/`                 | `builder.go`                                |
| `core/embedding/`        | `embedder.go`                               |
| `core/retrieval/`        | `adaptive_retrieval.go`                     |

**当前 Spec 仅覆盖了 `context.go` 和 `error_class.go`**，其余 15 文件均未被任何 Spec 涉及。

> **建议**: 这是一个系统性 DIP 问题，scope 远超当前 Spec 范围。建议：
> 1. 当前阶段：在 P3-02 中标注此发现，但 **不扩大 scope**
> 2. 后续迭代：新建独立 Epic 处理 `store` 包依赖反转

---

### 🟡 M-2: `server` 包的额外引入点

除 `error_class.go` 外，还有 4 处直接导入 `server/` 包：

| 文件                                   | 导入                      |
| :------------------------------------- | :------------------------ |
| `tools/scheduler.go`                   | `server/service/schedule` |
| `tools/preference.go`                  | `server/service/schedule` |
| `tools/memo_search.go`                 | `server/queryengine`      |
| `core/retrieval/adaptive_retrieval.go` | `server/queryengine`      |

> **评估**: `scheduler.go` 和 `preference.go` 本身就是日程工具，**直接依赖 schedule Service 接口是合理的设计**（工具层直接使用业务服务）。但 `core/retrieval` 引用 `queryengine` 值得关注。  
> **建议**: 无需新增 Spec，但需在 P2-01（scheduler 拆分）中确认拆分后 `schedule.Service` 的注入方式。

---

### 🟡 M-3: `calculateCost` 硬编码模型定价 (OCP)

`agents/base_parrot.go` L141-160 中的 `calculateCost` 方法通过 `switch` 硬编码模型定价：

```go
switch {
case strings.Contains(modelLower, "deepseek"):
    inputPricePerMillion = 0.14
case strings.Contains(modelLower, "gpt-4"):
    inputPricePerMillion = 2.50
case strings.Contains(modelLower, "gpt-3.5"):
    inputPricePerMillion = 0.15
default:
    inputPricePerMillion = 0.14
}
```

> **评估**: 这是 OCP 的轻度违规（对扩展不开放），但当前仅影响成本估算的精确性，**不影响核心功能**。  
> **建议**: 不在本次重构中处理，记录为后续技术债。

---

## 审计维度三：Spec 可行性验证

| Spec  | 可行性       | 风险点                                                               |
| :---- | :----------- | :------------------------------------------------------------------- |
| P1-01 | ✅ 可行       | 无                                                                   |
| P1-02 | ✅ 可行       | 需确认 `llm.Chat` 的返回值是否覆盖 `CreateChatCompletion` 的全部场景 |
| P1-03 | ✅ 可行       | 需全量搜索外部引用                                                   |
| P1-04 | ✅ 可行       | 无                                                                   |
| P1-05 | ✅ 可行       | 三处 LRU 的 TTL 策略差异需仔细对齐                                   |
| P1-06 | ⚠️ **需修正** | 漏覆盖 `tags/` 和 `duplicate/` 中的变体                              |
| P2-01 | ✅ 可行       | 工作量最大，建议分多个 PR                                            |
| P2-02 | ✅ 可行       | 同包拆分，风险可控                                                   |
| P3-01 | ✅ 可行       | `Extensions map[string]any` 的类型安全需辅助函数保障                 |
| P3-02 | ✅ 可行       | 需让上层类型实现 `ConflictError` 接口                                |
| P3-03 | ✅ 可行       | 无                                                                   |
| P3-04 | ✅ 可行       | 需确认消费方使用模式                                                 |
| P3-05 | ⚠️ **需修正** | SafeCallbackFunc 签名描述有误，需重新评估统一策略                    |
| P3-06 | ✅ 可行       | 需确认非 routing 路径是否独立使用 agents/ 分类器                     |
| P3-07 | ⚠️ **需修正** | 接口签名描述与实际不符，核心差异是 `history` 参数                    |
| P4-01 | ✅ 可行       | 需谨慎设计注册时序                                                   |
| P4-02 | ✅ 可行       | 无                                                                   |

---

## 审计结论

### 必须修正（3 项）

| #    | 类型   | Spec  | 修正内容                                                                     |
| :--- | :----- | :---- | :--------------------------------------------------------------------------- |
| 1    | 遗漏   | P1-06 | 扩展 truncate 扫描范围：新增 `tags/layer3_llm.go`、`duplicate/similarity.go` |
| 2    | 不准确 | P3-05 | 修正 SafeCallbackFunc 签名描述和定义位置                                     |
| 3    | 不准确 | P3-07 | 修正 ParrotAgent 接口签名，改 rationale 为 history 参数可选化                |

### 建议记录但不纳入本期（2 项）

| #    | 类型       | 范围                     | 建议          |
| :--- | :--------- | :----------------------- | :------------ |
| 1    | 系统性 DIP | `store` 包 17 处直接引入 | 后续独立 Epic |
| 2    | OCP 轻度   | `calculateCost` 硬编码   | 技术债清单    |
