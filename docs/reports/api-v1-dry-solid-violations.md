# API V1 代码质量与重构指南

> **状态**: 动态文档
> **目标**: 指导 API V1 的增量重构，使其符合 DRY/SOLID 原则。
> **最后更新**: 2026-02-20
> **基于扫描**: 2026-02-18

---

## 一、执行摘要

本文档分析了 `server/router/api/v1` 包内的技术债务，重点关注 DRY (Don't Repeat Yourself) 和 SOLID 原则的违约情况。

**主要发现：**
*   **重复逻辑**: 关键业务逻辑（认证/授权、AI 检查、资源解析）在数十个处理程序中被复制粘贴。
*   **上帝类**: `APIV1Service` 和 `user_service.go` 已演变成难以维护的单体，违反了单一职责原则 (SRP)。
*   **紧耦合**: 服务通常依赖于具体实现而非接口，导致测试和实现替换变得困难。

**战略建议**: 采取 **“减法重构”** 策略。与其进行复杂的重写，不如系统地将重复代码提取到共享的辅助函数和中间件中，在缩小代码库规模的同时提高稳定性。

---

## 二、重构原则

1.  **安全第一**: 没有测试覆盖就没有重构。如果缺少测试，先添加特性测试 (Characterization Tests)。
2.  **增量演进**: 优先选择小的、原子的变更（例如“提取辅助函数”），而不是“大爆炸”式的重写。
3.  **标准化**: 使用成熟的模式（拦截器、装饰器、映射器）来替换临时逻辑。

---

## 三、高优先级 DRY 违约 (速赢项)

这些违约影响较大，但修复工作量相对较低。解决这些问题将立即缩小代码库并降低 Bug 风险。

### 3.1 AI 服务可用性检查

**问题**: 每个 AI 相关方法都检查 `s.AIService == nil || !s.AIService.IsEnabled()`。
**出现频率**: `connect_handler.go` (30+ 次)。

**建议方案**: `RequireAI` 辅助函数/包装器
创建一个通用的包装器或简单的辅助方法。

```go
// 方案 A: 简单辅助函数 (侵入性最小)
func (s *ConnectServiceHandler) requireAI() error {
    if s.AIService == nil || !s.AIService.IsEnabled() {
        return connect.NewError(connect.CodeUnavailable, fmt.Errorf("AI features are disabled"))
    }
    return nil
}

// 使用示例:
func (s *ConnectServiceHandler) SuggestTags(...) (...) {
    if err := s.requireAI(); err != nil { return nil, err }
    // ...
}
```

### 3.2 用户认证与权限检查

**问题**: 重复出现的 `fetchCurrentUser` -> `nil check` -> `permission check` 模式。
**出现频率**: `user_service.go`, `memo_service.go`。

**建议方案**: 权限辅助函数
封装所有权验证逻辑。

```go
// 在新文件: server/router/api/v1/permissions.go
func (s *APIV1Service) requireUserAccess(ctx context.Context, targetUserID int32) (*store.User, error) {
    currentUser, err := s.fetchCurrentUser(ctx)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
    }
    if currentUser == nil {
        return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
    }
    // 如果用户是目标用户，或是管理员/Host，则允许访问
    if currentUser.ID != targetUserID && !isSuperUser(currentUser) {
        return nil, status.Errorf(codes.PermissionDenied, "permission denied")
    }
    return currentUser, nil
}
```

### 3.3 资源名称解析

**问题**: 手动字符串操作（`strings.Split`, `TrimPrefix`）分散在各个处理程序中。
**出现频率**: 所有服务文件。

**建议方案**: 类型化资源解析器
集中解析逻辑以确保一致性并防止解析错误。

```go
// 在 resource_name.go 中
type ResourceName struct {
    Type string
    ID   string
    // ...
}

func ParseResourceName(name string) (ResourceName, error) {
    // 集中处理逻辑
}

func ExtractID(name, prefix string) (int32, error) {
     // 集中安全解析
}
```

### 3.4 Connect 错误转换

**问题**: 每个包装器中都调用 `convertGRPCError(err)`。
**出现频率**: `connect_handler.go`。

**建议方案**: 中间件/拦截器
由于 `ConnectServiceHandler` 本身就是一个包装器，我们可以将其移动到拦截器中或简化包装器。

**建议**: 暂时保持原样（低风险），优先处理 3.1 和 3.2。

---

## 四、高优先级 SOLID 违约 (结构健康度)

这些问题需要更多投入，但对于长期可维护性至关重要。

### 4.1 上帝结构体 (`APIV1Service`)

**问题**: `APIV1Service` 实现了 *所有* v1 gRPC 服务。
**违约**: 单一职责原则 (SRP)。

**建议方案**: 服务组合
将 `APIV1Service` 拆分为更小、更专注的服务，由 `ConnectServiceHandler` 进行聚合。

```go
// 当前
type APIV1Service struct { ... } // 实现所有接口

// 建议
type UserService struct { ... }
type ScheduleService struct { ... }
type MemoService struct { ... }

type APIV1Service struct {
    User     *UserService
    Schedule *ScheduleService
    Memo     *MemoService
    // ...
}
```

### 4.2 巨型文件 (`user_service.go`)

**问题**: 包含 CRUD、设置、Token、Webhooks、通知（1400+ 行）。
**违约**: SRP。

**建议方案**: 文件拆分 (无逻辑变更)
按领域拆分文件。这保留了 Git 历史（大部分）并使代码审查更容易。
*   `user_service_crud.go`: 用户创建、获取、更新、删除。
*   `user_service_settings.go`: 用户设置。
*   `user_service_auth.go`: PAT (个人访问令牌)。
*   `user_service_webhook.go`: Webhooks。

### 4.3 硬编码更新逻辑 (`UpdateSchedule`)

**问题**: 使用 switch 语句手动将字段名映射到结构体更新。
**违约**: 开闭原则 (OCP)。

**建议方案**: 字段映射器 (Field Mappers)
定义一次映射，而不是在 switch 语句中。

```go
var scheduleFieldMappers = map[string]func(*v1pb.Schedule, *store.UpdateSchedule) {
    "title": func(pb *v1pb.Schedule, u *store.UpdateSchedule) { u.Title = &pb.Title },
    // ...
}
```

---

## 五、实施路线图

### 第一阶段：标准化 (第一周)
*   [ ] **任务 1.1**: 实现 `requireAI` 辅助函数并重构 `connect_handler.go`。
*   [ ] **任务 1.2**: 实现 `requireUserAccess` 辅助函数并重构 `user_service.go`, `memo_service.go`。
*   [ ] **任务 1.3**: 在 `schedule_service.go` 中提取 `convertRemindersToJSON` 辅助函数。

### 第二阶段：拆分 (第二周)
*   [ ] **任务 2.1**: 将 `user_service.go` 拆分为 4 个组件文件。
*   [ ] **任务 2.2**: 如果需要，拆分 `memo_service.go` (视大小/复杂度而定)。

### 第三阶段：架构改进 (第三周及以后)
*   [ ] **任务 3.1**: 定义 `AIService` 接口以解耦实现。
*   [ ] **任务 3.2**: 重构 `APIV1Service` 组合关系 (破坏性变更，需仔细协调)。

---

## 六、验证计划

对于每个重构任务：
1.  **运行现有测试**: `go test ./server/router/api/v1/...`
2.  **无逻辑变更验证**: 确保行为保持完全一致。
3.  **Lint 检查**: 确保没有引入新的 Linter 错误。
