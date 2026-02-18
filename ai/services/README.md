# AI Services (`ai/services`)

`services` 包包含了 AI 模块的具体业务逻辑实现，它们组合了多个底层的 AI 能力（如 LLM, Embedding, Store）来提供完整的业务功能。

## Subpackages

### `schedule`

负责日程管理相关的智能逻辑，特别是重复规则的解析和计算。

- **Recurrence**: 处理 RRule (Recurrence Rule) 标准，支持复杂的重复日程生成
  - 频率: 每日、每周、每月
  - 间隔: 每 N 天/周/月
  - **RecurrenceIterator**: 延迟加载迭代器，内存高效地遍历重复实例
- **Timezone Handling**: 基于 UTC 的时间计算，支持用户时区转换

### `session`

管理 AI 对话会话 (Session)。

- **SessionService**: 会话持久化接口
  - `SaveContext` / `LoadContext`: 会话上下文保存和加载
  - `ListSessions`: 用户会话列表
  - `DeleteSession`: 删除会话（用户隐私控制）
  - `CleanupExpired`: 清理过期会话
- **实现**: PostgreSQL 持久化 + Redis 缓存 (30分钟 TTL)

### `stats`

统计和成本告警服务。

- **CostAlertService**: 成本阈值监控和告警
  - 会话成本超阈值告警
  - 每日预算超支告警
  - 每日预算低于 10% 预警
- **Persister**: 代理统计数据持久化

## 目录结构

```
ai/services/
├── schedule/
│   ├── recurrence.go      # 重复规则解析和迭代器
│   └── helpers.go        # 辅助函数
├── session/
│   ├── interface.go      # 会话服务接口
│   ├── store.go          # PostgreSQL + 缓存实现
│   ├── recovery.go       # 会话恢复逻辑
│   └── cleanup.go       # 过期会话清理
└── stats/
    ├── alerting.go       # 成本告警服务
    └── persister.go      # 统计持久化
```
