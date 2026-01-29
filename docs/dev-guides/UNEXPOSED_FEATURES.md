# 已实现但缺少功能入口的功能

> 本文档记录 DivineSense 代码库中**后端已实现、前端 hooks/组件已存在**，但尚未在 UI 中暴露给用户的功能。

最后更新：2025-01-29

---

## 概述

DivineSense 拥有完善的 AI 功能基础设施，包括：
- ✅ 完整的后端 API 实现
- ✅ 前端 React hooks 和组件
- ✅ 路由配置和页面结构

但部分功能缺少 UI 入口，需要补充。

---

## 1. 编辑器增强功能

### 1.1 重复检测 (DetectDuplicates)

**实现位置：**
- 后端：`server/router/api/v1/ai_service_duplicate.go`
- 前端 Hook：`web/src/hooks/useAIQueries.ts` - `useDetectDuplicates()`
- Proto API：`DetectDuplicatesRequest` / `DetectDuplicatesResponse`

**功能描述：**
检测与当前 memo 重复或高度相似的内容，帮助用户避免创建重复笔记。

**返回数据结构：**
```typescript
{
  has_duplicate: boolean;      // 相似度 > 90%
  has_related: boolean;        // 相似度 70-90%
  duplicates: SimilarMemo[];   // 重复 memo 列表
  related: SimilarMemo[];      // 相关 memo 列表
  latency_ms: number;
}
```

**建议入口：**
```typescript
// 选项 1: 编辑器工具栏按钮
<EditorToolbar>
  <Button onClick={handleDetectDuplicates}>
    检测重复
  </Button>
</EditorToolbar>

// 选项 2: 自动检测（保存时）
const handleSave = async () => {
  const duplicates = await detectDuplicates(content);
  if (duplicates.has_duplicate) {
    showDuplicateDialog(duplicates.duplicates);
  }
};
```

---

### 1.2 合并 Memos (MergeMemos)

**实现位置：**
- 后端：`server/router/api/v1/ai_service_duplicate.go`
- 前端 Hook：`web/src/hooks/useAIQueries.ts` - `useMergeMemos()`
- Proto API：`MergeMemosRequest` / `MergeMemosResponse`

**功能描述：**
将源 memo 的内容合并到目标 memo，保留目标 memo 的唯一标识符。

**建议入口：**
```typescript
// Memo 操作菜单
<Menu>
  <MenuItem onClick={() => openMergeDialog(memoUid)}>
    合并到...
  </MenuItem>
</Menu>
```

---

### 1.3 关联 Memos (LinkMemos)

**实现位置：**
- 后端：`server/router/api/v1/ai_service_duplicate.go`
- 前端组件：`web/src/components/MemoEditor/components/LinkMemoDialog.tsx` ✅ 已实现
- Proto API：`LinkMemosRequest` / `LinkMemosResponse`

**功能描述：**
在两个 memo 之间创建双向关联关系，便于相互引用。

**状态：**
- 组件已实现，但未集成到编辑器或操作菜单

**建议入口：**
```typescript
// Memo 操作菜单
<Menu>
  <MenuItem onClick={() => openLinkMemoDialog(memoUid)}>
    关联 Memo...
  </MenuItem>
</Menu>
```

---

### 1.4 相关 Memos (GetRelatedMemos)

**实现位置：**
- 后端：`server/router/api/v1/ai_service_related.go`
- 前端 Hook：`web/src/hooks/useAIQueries.ts` - `useRelatedMemos()`
- 前端组件：`web/src/components/MemoRelatedList.tsx` ✅ 已实现

**功能描述：**
获取与指定 memo 相关的其他 memos，基于向量相似度。

**状态：**
- 组件已实现，但未在 memo 详情页集成

**建议入口：**
```typescript
// Memo 详情页底部
<MemoDetail memoUid={uid}>
  {/* 现有内容 */}
  <MemoContent />

  {/* 添加相关内容区域 */}
  <RelatedMemosSection memoUid={uid} />
</MemoDetail>
```

---

## 2. Proto API 已定义但前端未使用

### 2.1 SuggestTags（AI 标签建议）

| 状态 | 说明 |
|------|------|
| 后端 API | ✅ `SuggestTags` RPC 已实现 |
| 前端 Hook | ✅ `useSuggestTags` 已定义 |
| 前端组件 | ✅ `AITagSuggestPopover` 已实现 |
| UI 集成 | ✅ 已集成到编辑器工具栏 |

**结论：此功能已完整实现，无需额外工作。**

---

### 2.2 SemanticSearch（语义搜索）

| 状态 | 说明 |
|------|------|
| 后端 API | ✅ `SemanticSearch` RPC 已实现 |
| 前端 Hook | ✅ `useSemanticSearch` 已定义 |
| UI 入口 | ✅ 搜索栏已集成 |

**结论：此功能已完整实现，无需额外工作。**

---

## 3. 功能实现优先级

### P0 - 高优先级（用户最常用）

| 功能 | 实现复杂度 | 用户价值 |
|------|-----------|---------|
| 重复检测 | 低 | 🔴 高 - 避免重复内容 |
| 相关 Memos | 低 | 🔴 高 - 发现关联内容 |

### P1 - 中优先级（增强体验）

| 功能 | 实现复杂度 | 用户价值 |
|------|-----------|---------|
| 关联 Memos | 中 | 🟡 中 - 建立内容关联 |
| 合并 Memos | 中 | 🟡 中 - 整理重复内容 |

---

## 4. 实现建议

### 4.1 添加 Memo 操作菜单

创建统一的 memo 操作菜单组件：

```typescript
// web/src/components/MemoActionsMenu/index.tsx
import { Menu } from "@/components/ui/menu";
import {
  Link2,
  MoreVerticalIcon,
  Target,
  FileSearch,
  Link2,
  Merge,
} from "lucide-react";

interface MemoActionsMenuProps {
  memoUid: string;
  onDetectDuplicates?: () => void;
  onRelatedMemos?: () => void;
  onLinkMemo?: () => void;
  onMerge?: () => void;
}

export function MemoActionsMenu({ memoUid, ...handlers }: MemoActionsMenuProps) {
  return (
    <Menu>
      <MenuTrigger>
        <Button variant="ghost" size="icon">
          <MoreVerticalIcon className="w-4 h-4" />
        </Button>
      </MenuTrigger>
      <MenuContent>
        {handlers.onRelatedMemos && (
          <MenuItem onClick={handlers.onRelatedMemos}>
            <FileSearch className="w-4 h-4 mr-2" />
            查找相关
          </MenuItem>
        )}
        {handlers.onDetectDuplicates && (
          <MenuItem onClick={handlers.onDetectDuplicates}>
            <Target className="w-4 h-4 mr-2" />
            检测重复
          </MenuItem>
        )}
        {handlers.onLinkMemo && (
          <MenuItem onClick={handlers.onLinkMemo}>
            <Link2 className="w-4 h-4 mr-2" />
            关联 Memo
          </MenuItem>
        )}
        {handlers.onMerge && (
          <MenuItem onClick={handlers.onMerge}>
            <Merge className="w-4 h-4 mr-2" />
            合并到...
          </MenuItem>
        )}
      </MenuContent>
    </Menu>
  );
}
```

### 4.2 Memo 详情页添加相关内容区域

```typescript
// web/src/pages/MemoDetail.tsx 添加区域
import { MemoRelatedList } from "@/components/MemoRelatedList";

function MemoDetailPage() {
  // ... 现有代码

  return (
    <div className="memo-detail-page">
      {/* 现有内容 */}

      {/* 新增：相关 Memos 区域 */}
      <div className="related-memos-section mt-8">
        <MemoRelatedList memoUid={uid} />
      </div>
    </div>
  );
}
```

### 4.3 编辑器工具栏添加重复检测按钮

```typescript
// web/src/components/MemoEditor/components/EditorToolbar.tsx
import { useDetectDuplicates } from "@/hooks/useAIQueries";
import { FileSearch } from "lucide-react";

function EditorToolbar({ content, disabled }) {
  const detectDuplicates = useDetectDuplicates();

  const handleDetectDuplicates = async () => {
    const result = await detectDuplicates(
      { content, limit: 5 },
      {
        onSuccess: (data) => {
          if (data.has_duplicate || data.has_related) {
            showDuplicateDialog(data);
          } else {
            toast.success("未发现重复内容");
          }
        },
      },
    );
  };

  return (
    <Toolbar>
      {/* 现有按钮 */}
      <AITagSuggestPopover content={content} onInsertTags={handleInsertTags} />

      {/* 新增：重复检测按钮 */}
      <Button
        variant="ghost"
        size="icon"
        onClick={handleDetectDuplicates}
        disabled={disabled || !content}
        title="检测重复内容"
      >
        <FileSearch className="w-4 h-4" />
      </Button>
    </Toolbar>
  );
}
```

---

## 5. 后端 API 参考

### 5.1 DetectDuplicates

**RPC 调用：**
```go
service AIService {
  rpc DetectDuplicates(DetectDuplicatesRequest) returns (DetectDuplicatesResponse) {
    option (google.api.http) = {
      post: "/api/v1/ai/detect-duplicates"
      body: "*"
    };
  }
}
```

**请求参数：**
```protobuf
message DetectDuplicatesRequest {
  string title = 1;                      // optional
  string content = 2 [(google.api.field_behavior) = REQUIRED];
  repeated string tags = 3;                // optional
  int32 top_k = 4;                        // default: 5
}
```

### 5.2 MergeMemos

**RPC 调用：**
```go
rpc MergeMemos(MergeMemosRequest) returns (MergeMemosResponse) {
  option (google.api.http) = {
    post: "/api/v1/ai/merge-memos"
    body: "*"
  };
}
```

### 5.3 LinkMemos

**RPC 调用：**
```go
rpc LinkMemos(LinkMemosRequest) returns (LinkMemosResponse) {
  option (google.api.http) = {
    post: "/api/v1/ai/link-memos"
    body: "*"
  };
}
```

### 5.4 GetRelatedMemos

**RPC 调用：**
```go
rpc GetRelatedMemos(GetRelatedMemosRequest) returns (GetRelatedMemosResponse) {
  option (google.api.http) = {
    get: "/api/v1/{name=memos/*}/related"
  };
}
```

---

## 6. 前端 Hooks 参考

### 6.1 useDetectDuplicates

```typescript
// web/src/hooks/useAIQueries.ts
export function useDetectDuplicates() {
  return useMutation({
    mutationFn: async (params: {
      title?: string;
      content: string;
      tags?: string[];
      topK?: number;
    }) => {
      const request = create(DetectDuplicatesRequestSchema, {
        title: params.title ?? "",
        content: params.content,
        tags: params.tags ?? [],
        topK: params.topK ?? 5,
      });
      return await aiServiceClient.detectDuplicates(request);
    },
  });
}
```

### 6.2 useMergeMemos

```typescript
export function useMergeMemos() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (params: {
      sourceName: string;  // memos/{uid}
      targetName: string;  // memos/{uid}
    }) => {
      const request = create(LinkMemosRequestSchema, {
        sourceName: params.sourceName,
        targetName: params.targetName,
      });
      return await aiServiceClient.linkMemos(request);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["memos"] });
    },
  });
}
```

### 6.3 useRelatedMemos

```typescript
export function useRelatedMemos(
  name: string,  // memos/{uid}
  options: { enabled?: boolean; limit?: number } = {}
) {
  return useQuery({
    queryKey: aiKeys.related(name),
    queryFn: async () => {
      const request = create(GetRelatedMemosRequestSchema, {
        name,
        limit: options.limit ?? 5,
      });
      return await aiServiceClient.getRelatedMemos(request);
    },
    enabled: (options.enabled ?? true) && !!name && name.startsWith("memos/"),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
```

---

## 7. 已完整实现的功能（参考）

以下功能已实现完整入口，无需额外工作：

| 功能 | 入口 | 位置 |
|------|------|------|
| **AI Chat** | 侧边栏 Sparkles 图标 | `/chat` |
| **Schedule Agent** | 侧边栏 Calendar 图标 | `/schedule` |
| **Knowledge Graph** | 侧边栏 Share2 图标 | `/knowledge-graph` |
| **Review System** | 侧边栏 Target 图标 | `/review` |
| **Geek Mode** | Chat Header 切换按钮 | `/chat` |
| **AI 标签建议** | 编辑器工具栏 | Memo 编辑器 |
| **语义搜索** | 搜索栏 | 主页搜索 |

---

## 8. 更新日志

| 日期 | 更新内容 |
|------|----------|
| 2025-01-29 | 初始版本，记录 4 个未暴露功能 |
