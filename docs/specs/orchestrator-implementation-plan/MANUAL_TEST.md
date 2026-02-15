# Orchestrator 手动验收测试用例

> **关联 Issue**: #213 | **状态**: 待执行

## Agent 能力概览

| Agent | 工具 | 核心能力 |
|:------|:-----|:---------|
| **memo** | `memo_search` | 搜索笔记、按时间浏览、语义查找 |
| **schedule** | `schedule_add`, `schedule_query`, `schedule_update`, `find_free_time` | 创建日程、查询日程、更新日程、查找空闲时间 |

## 测试环境

| 配置 | 值 |
|:-----|:---|
| 数据库 | PostgreSQL (divinesense-postgres-dev) |
| 前端端口 | 25173 |
| 后端端口 | 28081 |
| Trace ID | 可从日志中获取 |

---

## 测试场景

### 场景 1: 线性任务依赖（先搜笔记，再建日程）

**目的**: 验证 DAG 按依赖顺序串行执行

**Agent 组合**: `memo` → `schedule`

**测试输入**:
```
用户输入: "找到我之前记的关于 PostgreSQL 配置的内容，然后帮我安排时间复习"
```

**预期任务分解**:
```json
{
  "tasks": [
    {
      "id": "task_search_db_config",
      "agent": "memo",
      "input": "搜索关于 PostgreSQL 配置的笔记",
      "purpose": "查找数据库配置相关笔记",
      "dependencies": []
    },
    {
      "id": "task_create_review",
      "agent": "schedule",
      "input": "查找空闲时间，用于复习笔记内容：{{task_search_db_config.result}}",
      "purpose": "安排复习时间",
      "dependencies": ["task_search_db_config"]
    }
  ]
}
```

**验证要点**:
- [ ] 日志显示 `task_search_db_config` 在 `task_create_review` 之前执行
- [ ] `task_create_review` 的输入中 `{{task_search_db_config.result}}` 被替换为实际搜索结果
- [ ] 最终响应包含笔记搜索结果 + 日程创建结果

---

### 场景 2: 钻石依赖并行（同时搜索笔记和查询日程）

**目的**: 验证 DAG 最大化并行执行

**Agent 组合**: `memo` + `schedule` → 聚合

**测试输入**:
```
用户输入: "搜索关于 AI 项目的笔记，同时查看今天下午的日程安排"
```

**预期任务分解**:
```json
{
  "tasks": [
    {
      "id": "task_search_ai",
      "agent": "memo",
      "input": "搜索 AI 项目相关笔记",
      "purpose": "查找 AI 项目笔记",
      "dependencies": []
    },
    {
      "id": "task_check_today",
      "agent": "schedule",
      "input": "查询今天的日程安排",
      "purpose": "查看今日日程",
      "dependencies": []
    },
    {
      "id": "task_aggregate",
      "agent": "aggregator",
      "purpose": "整合两个任务结果",
      "dependencies": ["task_search_ai", "task_check_today"]
    }
  ]
}
```

**验证要点**:
- [ ] 日志显示 `task_search_ai` 和 `task_check_today` 并行执行（时间戳接近）
- [ ] `task_aggregate` 等待前两个任务都完成后才开始
- [ ] 最终响应同时包含笔记搜索结果和日程查询结果

---

### 场景 3: 上下文注入（基于笔记内容创建日程）

**目的**: 验证 `{{task_id.result}}` 变量替换

**Agent 组合**: `memo` → `schedule`

**测试输入**:
```
用户输入: "我之前记了产品需求文档在哪里？帮我安排时间 review"
```

**预期任务分解**:
```json
{
  "tasks": [
    {
      "id": "task_find_prd",
      "agent": "memo",
      "input": "搜索产品需求文档",
      "purpose": "找到 PRD 笔记",
      "dependencies": []
    },
    {
      "id": "task_schedule_review",
      "agent": "schedule",
      "input": "创建日程：review 产品需求文档，内容要点：{{task_find_prd.result}}",
      "purpose": "安排 review 时间",
      "dependencies": ["task_find_prd"]
    }
  ]
}
```

**验证要点**:
- [ ] `task_find_prd` 先执行并返回笔记位置/内容
- [ ] `task_schedule_review` 的输入中 `{{task_find_prd.result}}` 被替换为实际笔记内容
- [ ] 日程创建成功，包含笔记内容摘要

---

### 场景 4: 失败重试与 Exponential Backoff

**目的**: 验证任务失败后的重试机制

**测试方法**: Mock 一个必然失败的 Agent 调用

**预期日志**:
```
INFO executor: task start trace_id=xxx task_id=task_xxx agent=memo
WARN executor: task execution failed, retrying trace_id=xxx task_id=task_xxx attempt=1/3 error="xxx"
INFO executor: sleeping 1s before retry
WARN executor: task execution failed, retrying trace_id=xxx task_id=task_xxx attempt=2/3 error="xxx"
INFO executor: sleeping 2s before retry
WARN executor: task execution failed, retrying trace_id=xxx task_id=task_xxx attempt=3/3 error="xxx"
INFO executor: sleeping 4s before retry
ERROR executor: task failed trace_id=xxx task_id=task_xxx error="xxx" retry_count=3
```

**验证要点**:
- [ ] 日志显示 `attempt=1/3`, `attempt=2/3`, `attempt=3/3`
- [ ] 重试间隔符合 Exponential Backoff：1s → 2s → 4s
- [ ] 3 次重试后任务标记为 Failed

---

### 场景 5: 级联跳过（笔记搜索失败 → 后续日程取消）

**目的**: 验证上游失败后下游任务被跳过

**Agent 组合**: `memo` → `schedule`

**测试方法**: Mock `memo` Agent 必然失败

**预期任务分解**:
```json
{
  "tasks": [
    {
      "id": "task_search",
      "agent": "memo",
      "input": "搜索不存在的关键词 xyz123456",
      "purpose": "搜索笔记",
      "dependencies": []
    },
    {
      "id": "task_create_schedule",
      "agent": "schedule",
      "input": "基于搜索结果创建日程",
      "purpose": "创建日程",
      "dependencies": ["task_search"]
    }
  ]
}
```

**验证要点**:
- [ ] `task_search` 失败（无结果或错误）
- [ ] `task_create_schedule` 状态变为 `skipped`，原因包含 "upstream failure"
- [ ] 日志显示级联跳过

---

### 场景 6: 循环依赖检测

**目的**: 验证 DAG 能检测循环依赖并报错

**测试方法**: 构造循环依赖

**预期任务分解**:
```json
{
  "tasks": [
    {
      "id": "task_a",
      "agent": "memo",
      "dependencies": ["task_c"]  // 形成 A → B → C → A 循环
    },
    {
      "id": "task_b",
      "agent": "schedule",
      "dependencies": ["task_a"]
    },
    {
      "id": "task_c",
      "agent": "memo",
      "dependencies": ["task_b"]
    }
  ]
}
```

**验证要点**:
- [ ] 返回错误：`cycle detected or deadlock: 0/X tasks completed`
- [ ] 无任务被执行

---

## 附录: 可用 Agent 配置

### memo (MemoParrot / 灰灰)
- **工具**: `memo_search`
- **能力**: 搜索笔记内容、按时间浏览笔记、查找相关内容
- **适用**: 用户想查找之前记录的信息、浏览历史笔记、搜索特定话题

### schedule (ScheduleParrot / 时巧)
- **工具**: `schedule_add`, `schedule_query`, `schedule_update`, `find_free_time`
- **能力**: 创建日程、查询日程、更新日程、查找空闲时间
- **适用**: 时间管理、会议安排、查看日程、修改已有安排、查找空闲时间

---

## 执行清单

- [ ] 场景 1: 线性任务依赖（memo → schedule）
- [ ] 场景 2: 钻石依赖并行（memo + schedule → aggregator）
- [ ] 场景 3: 上下文注入
- [ ] 场景 4: 重试与 Backoff
- [ ] 场景 5: 级联跳过
- [ ] 场景 6: 循环依赖检测
