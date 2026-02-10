# API 设计规范

> DivineSense API 设计最佳实践 — REST/gRPC/Connect RPC

---

## 核心原则

1. **Protocol First**：先修改 `.proto` 文件，再生成代码
2. **版本化**：使用 `v1`、`v2` 路径版本控制
3. **向后兼容**：新版本不破坏旧版本客户端

---

## Connect RPC 规范

### 服务定义

```go
// proto/api/v1/xxx_service.proto

service XXXService {
  rpc Get(GetRequest) returns (GetResponse) {}
  rpc List(ListRequest) returns (ListResponse) {}
  rpc Create(CreateRequest) returns (CreateResponse) {}
  rpc Update(UpdateRequest) returns (UpdateResponse) {}
  rpc Delete(DeleteRequest) returns (DeleteResponse) {}
}
```

### 请求/响应模式

```protobuf
message GetRequest {
  int64 id = 1;
}

message GetResponse {
  XXX item = 1;
  // 使用 Connect 的 metadata 传递错误
}

message ListRequest {
  int32 page = 1;
  int32 per_page = 2;
  string filter = 3;
}

message ListResponse {
  repeated XXX items = 1;
  int64 total = 2;
}
```

---

## 错误处理

### 错误码规范

```go
const (
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeAlreadyExists = "ALREADY_EXISTS"
	ErrCodeInternal      = "INTERNAL_ERROR"
	ErrCodeUnauthorized  = "UNAUTHORIZED"
)
```

### 错误响应

```go
return &connect.Response[Output]{{
	Msg: &connect.Message{
		Metadata: metadata.New(
			map[string]string{
				"error-code": ErrCodeValidation,
				"error-msg":  "Invalid input: ...",
			},
		),
	},
}}
```

---

## RESTful 规范

### URL 设计

| 操作 | HTTP 方法 | URL 模式 |
|:-----|:----------|:---------|
| 列表 | GET | /api/v1/resources |
| 详情 | GET | /api/v1/resources/{id} |
| 创建 | POST | /api/v1/resources |
| 更新 | PUT/PATCH | /api/v1/resources/{id} |
| 删除 | DELETE | /api/v1/resources/{id} |

### 命名约定

- 使用复数名词：`/api/v1/users`（非 `/user`）
- 使用 kebab-case：`/api/v1/schedule-events`（非 `scheduleEvents`）
- 版本路径：`/api/v1/` → `/api/v2/`

---

## 通用模式

### 分页

```protobuf
message PaginationRequest {
  int32 page = 1;      // 从 1 开始
  int32 per_page = 2; // 默认 20，最大 100
}

message PaginationResponse {
  int64 total = 1;
  int32 page = 2;
  int32 per_page = 3;
  bool has_more = 4;
}
```

### 排序

```protobuf
message SortOption {
  string field = 1;      // 排序字段
  string order = 2;      // "asc" 或 "desc"
}
```

### 过滤

```protobuf
message FilterOption {
  string field = 1;
  string operator = 2;  // "eq", "ne", "gt", "lt", "contains"
  string value = 3;
}
```

---

*文档路径：.claude/rules/api-design.md*
