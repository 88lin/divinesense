# API V1 优化重构执行方案 (Optimized Plan)

> **状态**: 待执行
> **目标**: 针对 `server/router/api/v1` 的技术债务进行“先破除巨石，后精确打击”的高效重构。
> **策略**: 激进拆分，稳妥去重，重构上帝类。

---

## 阶段一：破除巨石 (高优激进拆分)
鉴于当前无并发代码冲突，优先将巨石文件 `user_service.go` (1400+ 行) 根据业务领域进行**物理隔离**。此阶段不改变底层业务逻辑，仅做结构调整，为后续 DRY 重构扫清障碍，降低代码审查的复杂度。

*   **物理结构拆分规划**：
    *   `user_service.go`: 仅保留基础的 gRPC 接口声明、服务初始化以及公共的底层辅助逻辑。
    *   `user_service_crud.go`: 承接用户的核心增删改查逻辑 (`GetUser`, `CreateUser`, `UpdateUser`, `DeleteUser`, `ListUsers` 等)。
    *   `user_service_settings.go`: 承接用户个性化设置逻辑 (`GetUserSetting`, `UpdateUserSetting`, `ListUserSettings` 等)。
    *   `user_service_auth.go`: 承接个人访问令牌 (PAT) 相关逻辑 (`ListPersonalAccessTokens`, `CreatePersonalAccessToken`, `DeletePersonalAccessToken`)。
    *   `user_service_webhook.go`: 承接用户级别的 Webhooks 和通知设定逻辑。
    *   `user_service_converter.go`: 将所有跨域的、以 `convertUser...` 开头的 proto 转换辅助函数统一收敛到一处。

---

## 阶段二：DRY 去重与标准化提取
在完成物理拆解后，各业务域代码量大幅下降。此时开始精确打击核心重复代码，植入高复用性逻辑。

*   **2.1 统一权限守卫 (`RequireUserAccess`)**：
    *   **动作**: 新建 `server/router/api/v1/permissions.go`。
    *   **逻辑**: 提取高度重复的 `fetchCurrentUser` -> `nil` 检查 -> `所有权/管理员验证` 的样板代码。
    *   **落地**: 清洗并替换刚拆分出的 `user_service_*.go` 家族、以及庞大的 `memo_service.go` 内部的校验逻辑。
*   **2.2 AI 服务可用性检查 (`RequireAI`)**：
    *   **动作**: 将 `s.AIService == nil || !s.AIService.IsEnabled()` 的判断逻辑严密收敛为辅助包装函数或中间件拦截器。
    *   **落地**: 搜查并清理 `connect_handler.go` 及相关文件中存在的数十处冗余 `if` 块。
*   **2.3 规范化资源名称解析**：
    *   **动作**: 全面排查类似于 `strings.Split(name, "/")` 等高风险的硬编码字符串操作。
    *   **落地**: 将其统一点对点替换为 `resource_name.go` 提供的标准化资源解析函数，确保类型的强校验一致性。

---

## 阶段三：结构健康与 OCP 治理
在代码显著瘦身并消除大量坏味道后，补齐最后一环，保障系统的长期可维护性和扩展性。

*   **3.1 Schedule 更新重构**：
    *   **动作**: 消除 `schedule_service.go` (如 `UpdateSchedule`) 中由于使用 `switch struct` 导致的硬编码字段映射。
    *   **落地**: 采用 **Field Mappers (字段映射器)** 模式进行重构，提高可扩展性并使其坚决符合开闭原则 (OCP)。
*   **3.2 APIV1Service 上帝类解耦**：
    *   **动作**: 继续将 `UserService`、`MemoService`、`AttachmentService` 等大模块按照已有的 `AIService` / `ScheduleService` 模式从 `APIV1Service` 结构体中剥离解耦（采用结构体组合模式）。
    *   **落地**: 逐步引入接口层依赖而非具体实现依赖，彻底切断不同服务间的不合理紧耦合。

---

## 验证计划
在每个阶段的重构行动后，必须严格执行以下质量关卡：
1.  **运行现有测试**: `go test ./server/router/api/v1/...`
2.  **Lint 检查**: 确保没有引入新的 Linter 错误 (`golangci-lint run`)。
3.  **无损验证**: 确保对外 API 的行为特征和错误码响应保持完全一致。
