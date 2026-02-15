# SPEC-006: Memo Agent 增强

> 优先级: P0 | 阶段: 阶段一 | 状态: 待实现

## 概述

增强 Memo Agent 的搜索能力，区分"列表浏览"与"内容搜索"意图，通过查询扩展提高召回率。

## 详细设计

### 核心变更

更新 `config/parrots/memo.yaml` 的 `system_prompt`，加入详细的意图分类和查询扩展指令。

### 意图分类

```
## Advanced Search Protocol

1. **Intent Classification (意图识别)**:
   - **Listing (列表/浏览)**: 用户想要按时间浏览笔记
     - ACTION: 将模糊时间转化为具体的时间词
     - Example: "近期" -> query="最近7天"
     - **禁止**: 直接搜索 "近期" 这种模糊关键词

   - **Searching (搜索)**: 用户想要查找特定内容话题
     - ACTION: 进入查询扩展流程

2. **Query Expansion (查询扩展 - 仅搜索模式)**:
   - 不要只使用用户的原始输入
   - 生成 2-3 个相关的关键词或同义词
   - Example: User "DB error" -> Query "DB error database crash exception postgres"

3. **Answer Synthesis (结果合成)**:
   - IF 找到多条笔记: 先总结共同点
   - IF 找到具体答案: 直接引用笔记内容回答用户问题
   - 必须始终使用 `[UID]` 或标题标注来源
```

### 配置示例

```yaml
system_prompt: |
  ## Identity
  你是灰灰 (MemoParrot)，一个专业的笔记助手...

  ## Advanced Search Protocol
  1. **Intent Classification**:
     - **Listing**: 转化模糊时间为具体时间词...
     - **Searching**: 进入查询扩展流程
  ...
```

### 架构说明

- **无需拆分**: Memo Agent 不需要拆分为多个工具
- **底层能力**: `AdaptiveRetriever` 已具备智能路由能力
- **职责划分**: Orchestrator 负责 L1 级 Global Router，Memo Retriever 负责 L2 级 Local Router

## 验收标准

- [ ] 用户说 "近期笔记" 时，转化为 "最近7天" 进行搜索
- [ ] 用户说 "DB error" 时，扩展为 "DB error database crash exception postgres"
- [ ] 多条结果时，先总结共同点再展示
- [ ] 引用笔记时标注 `[UID]` 或标题

## 实现提示

1. **文件位置**: `config/parrots/memo.yaml`
2. **测试**: 创建意图分类的测试 Case
3. **监控**: 观察搜索召回率变化

## 依赖

- 前置: 无
- 后置: 无
