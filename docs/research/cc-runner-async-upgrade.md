# CCRunner 异步架构与 UX 升级调研报告

## 1. 核心需求
**目标**：消除 Geek Mode (Claude Code Runner) 的“黑盒感”，实现异步非阻塞交互，提升用户对 AI 内部行为的感知度。

## 2. 深度调研：Existing Capabilities

### 2.1 当前架构瓶颈
- **同步阻塞**: `Go backend → CCRunner → exec.Command(...).Wait()`。Web 请求必须维持连接，直至 CLI 执行完毕。
- **信息黑洞**: `CCRunner` 仅捕获 `text` 和 `tool_use.name`，丢弃了 `stdout`、`stderr` 以及具体的工具参数。

### 2.2 Claude Code CLI `stream-json` 能力挖掘
通过实际测试 (`claude --print --verbose "cmd" --output-format stream-json`)，发现其输出包含了构建 Step UI 所需的所有元数据：

**Tool Call (Assistant Message)**:
```json
{
  "type": "assistant",
  "message": {
    "content": [{
      "type": "tool_use",
      "name": "Bash",
      "input": { "command": "echo 'hello'", "description": "..." }
    }]
  }
}
```
> ✅ 可提取：具体执行的 Shell 命令。

**Tool Result (User Message)**:
```json
{
  "type": "user",
  "tool_use_result": {
    "stdout": "hello\n",
    "is_error": false
  }
}
```
> ✅ 可提取：Shell 命令的真实 stdout 输出。

**File Creation Result**:
```json
{
    "tool_use_result": {
        "type": "create", 
        "filePath": "/path/to/file",
        "content": "file content..."
    }
}
```
> ✅ 可提取：创建或修改的文件内容，支持预览。

## 3. 技术方案设计

### 3.1 后端架构升级 (Go)

```
[Web Client] <--(SSE)-- [Event Broker] <---- [Task Manager]
                             ^                    |
                             | (Status/Logs)      | (Submit Job)
                             |                    v
                        [Job Store]         [CCRunner Worker]
                             ^                    |
                             | (Update)           | (Execute)
                             ---------------------/
```

1.  **Task Manager**: 接受 `/api/v1/cc/submit`，返回 `TaskID`。立即释放 HTTP 请求。
2.  **Worker Pool**: 后台 goroutine 运行 `CCRunner.Execute`。
3.  **Event Broker**: 基于 SSE 或 WebSocket，根据 `SessionID` 推送实时事件。

### 3.2 CCRunner 解析逻辑增强

需重构 `streamOutput` 函数，实现新的状态机：

| 收到消息类型           | 提取数据             | 触发事件       | 用途                                     |
| :--------------------- | :------------------- | :------------- | :--------------------------------------- |
| `assistant` (tool_use) | `name`, `input`      | `step_start`   | 前端创建步骤卡片，显示 "Run: ls -la"     |
| `user` (tool_result)   | `stdout` / `content` | `step_output`  | 前端更新步骤卡片，显示终端输出或文件预览 |
| `assistant` (text)     | `text`               | `answer_delta` | 流式显示最终回答                         |

### 3.3 前端交互设计 (React)

1.  **任务状态卡片**: 聊天流中仅显示 "Running Task #123 [Details]"。
2.  **Execution Drawer (详情抽屉)**:
    *   **Timeline**: 垂直时间轴，串联 Thinking -> Command -> Output。
    *   **Terminal View**: 黑色背景，单色字体，还原 CLI 体验。
    *   **File Preview**: 代码高亮视图，展示 `Read` 或 `Write` 的内容。

## 4. 风险与规避
1.  **SSE 连接断开**: 实现前端重连 + 后端 Event Buffer (Replay Last-N events)。
2.  **CLI 输出变动**: Claude Code 可能会更新输出格式。需添加防御性解析代码，并锁定 CLI 版本。

## 5. 结论
技术路径清晰，完全可行。该升级将极大提升 Geek Mode 的专业感和可用性。
