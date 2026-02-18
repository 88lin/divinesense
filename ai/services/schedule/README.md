# AI Schedule Service (`ai/services/schedule`)

`schedule` 子包处理与日程相关的复杂业务逻辑，核心是重复规则 (Recurrence Rule) 的处理。

## 核心功能

### 重复规则 (Recurrence)

实现了对重复日程的解析和生成，支持以下模式：

- **频率**: 每日 (daily)、每周 (weekly)、每月 (monthly)
- **间隔**: 每 N 天/周/月
- **结束条件**: 永不结束、直到某日期、执行 N 次后结束
- **例外**: 排除特定日期

### RecurrenceIterator (延迟加载迭代器)

提供内存高效的重复实例遍历：

- **缓存机制**: 预加载部分实例，避免一次性生成所有实例
- **懒加载**: 按需生成，支持大规模重复日程
- **安全限制**: 最大 10 年或 1000 个实例

```go
rule := &RecurrenceRule{
    Type:     RecurrenceTypeWeekly,
    Weekdays: []int{5}, // 周五
    Interval: 1,
}
iterator := rule.Iterator(startTs)

// 延迟加载
instances := iterator.GetUntil(endTs)
```

## 算法流程

```mermaid
flowchart TD
    Input[自然语言: 每周五下午3点开会] --> TimeParse[ai/aitime: 解析基准时间]
    TimeParse --> RRuleGen[生成 RRule 规则\nType=weekly;Weekdays=[5];Interval=1]

    RRuleGen --> Expand{计算实例}
    Expand --> Iterator[RecurrenceIterator\n延迟加载]

    Iterator --> Next[Next: 下一次时间]
    Iterator --> List[List: 未来 N 次时间]

    Next & List --> Store[存入数据库]
```

## 业务流程

当用户创建一个重复日程（如"每周五下午3点开周会"）时：

1. **AI 解析**: `aitime` 包解析出基础时间点
2. **规则生成**: 本包将自然语言描述转换为结构化的 `RecurrenceRule`
3. **实例展开**: 使用 `RecurrenceIterator` 按需计算发生时间点 (Occurrences)

## 时区处理

- **UTC 基准**: 所有 Unix timestamp 使用 UTC
- **调用者负责**: 调用方在传入前完成用户时区到 UTC 的转换
- **返回处理**: 调用方在返回结果时完成 UTC 到用户时区的转换
