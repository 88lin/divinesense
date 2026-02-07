# AI 聊天消息排队优化 - 调研报告

> **调研日期**: 2026-02-07  
> **Issue**: [#121](https://github.com/hrygo/divinesense/issues/121)  
> **状态**: ✅ 已完成

---

## 需求概述

优化 AI 聊天界面的消息输入交互，建立智能的 Pending Queue 机制：

1. **当前会话没有 Block** → 创建新 Block
2. **最新 Block 状态是 COMPLETED/ERROR** → 创建新 Block  
3. **最新 Block 状态是 PENDING/STREAMING**：
   - **Normal 模式** → 消息进入 pending 队列
   - **Geek/Evolution 模式** → 消息立即 appendUserInput

---

## 规格确认

### Pending Queue 规格

| 维度 | 规则 |
|:-----|:-----|
| **合并策略** | 换行符 `\n` 连接 |
| **可见性** | 在用户输入区域上方显示 |
| **队列上限** | 10 条 |
| **超限处理** | 覆盖最老的消息（FIFO） |
| **出队时机** | 当前最新 Block 状态变为 done |
| **新 Block Mode** | 使用**当前 AI 模式** |
| **切换模式** | 不影响队列 |

### 行为示例

```
场景：用户快速连发三条消息

T1: 用户发送 "帮我分析" → Block A (NORMAL) 开始 streaming
T2: 用户发送 "Redis 的" → 进入队列
T3: 用户发送 "缓存策略" → 进入队列
T4: Block A 完成
T5: 出队 → 创建 Block B (当前模式)，内容 = "Redis 的\n缓存策略"
```

---

## 技术方案

### 前端数据模型

```typescript
interface PendingMessage {
  content: string;
  timestamp: number;
}

interface PendingQueueState {
  messages: PendingMessage[];
  targetBlockId: number | null;
  maxLength: number;  // 默认 10
}
```

### 核心逻辑

```typescript
// handleSend 流程
if (activeBlock) {
  const blockMode = getBlockModeName(activeBlock.mode);
  
  if (blockMode === "normal") {
    addPendingMessage(userMessage);  // 进入队列
  } else {
    await appendUserInput(...);      // Geek/Evolution: 立即追加
  }
}

// Block done 时的出队
useEffect(() => {
  if (wasActive && isNowTerminal && pendingQueue.messages.length > 0) {
    flushPendingQueue();
  }
}, [blocks]);
```

### 文件变更清单

| 文件 | 变更类型 |
|:-----|:---------|
| `web/src/types/aichat.ts` | 新增类型定义 |
| `web/src/contexts/AIChatContext.tsx` | 新增队列状态管理 |
| `web/src/pages/AIChat.tsx` | 修改 handleSend 逻辑 |
| `web/src/components/AIChat/PendingQueueBar.tsx` | 新增 UI 组件 |
| `web/src/locales/zh-Hans.json` | 国际化 |
| `web/src/locales/en.json` | 国际化 |

---

## 元认知评估

| 维度 | 评分 | 说明 |
|:-----|:-----|:-----|
| **需求清晰度** | 5/5 | 所有边界条件已确认 |
| **技术可行性** | 5/5 | 现有 API 支持完整 |
| **用户价值** | 4/5 | 改善交互体验 |
| **实现复杂度** | 3/5 | 中等复杂度 |
| **风险控制** | 4/5 | 风险可控 |

---

## 验收标准

- [ ] `make check-all` 通过
- [ ] Normal 模式 streaming 时消息进入队列
- [ ] Geek/Evolution 模式 streaming 时消息立即追加
- [ ] Block done 后队列自动合并发送
- [ ] 队列满时自动覆盖最老消息
- [ ] PendingQueueBar 正确显示状态
- [ ] 国际化文本正确
- [ ] 切换模式不影响队列

---

*调研完成于 2026-02-07*
