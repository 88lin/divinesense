# Geek Parrot (`ai/agents/geek`)

`geek` 包实现了面向**代码执行**和**自我进化**的高级 Agent —— `GeekParrot` 和 `EvolutionParrot`。它集成了 Claude Code CLI，使其具备直接操作文件系统、运行终端命令和编写代码的能力。

## 架构与能力

### 1. GeekParrot (Geek Mode)
*   **代码执行**: 封装了 Claude Code CLI 的交互。
*   **统一执行器**: 使用 `runner.CCRunner` 作为统一的执行层。
*   **会话管理**: 支持会话持久化，跨重启恢复。
*   **流式输出**: 实时流式传输代码执行的输出结果。

### 2. EvolutionParrot (Evolution Mode)
*   **仓库感知**: 能够理解当前代码库的结构和上下文。
*   **功能实现**: 接受开发任务，自动创建并修改文件。
*   **自我进化**: 具备修改自身代码的能力（管理员专用）。

### 3. GeekMode / EvolutionMode
*   **工作目录**: 每个用户有独立的沙盒目录。
*   **权限控制**: 基于用户权限的访问控制。
*   **系统提示**: 自动生成包含用户和设备上下文的系统提示。

## 安全机制

由于赋予了 Agent 极高的权限（文件读写、命令执行），必须有严格的安全限制：

*   **Danger Detection**: 拦截高危命令 (如 `rm -rf /`, `mkfs` 等)。
*   **路径检查**: 防止访问敏感目录。
*   **绕过模式**: 仅 Evolution 模式可绕过安全检查（管理员专用）。
*   **超时控制**: 强制执行超时，防止死循环或挂起。

## 业务流程

```mermaid
sequenceDiagram
    participant User
    participant GeekParrot
    participant CCRunner
    participant DangerDetector
    participant ClaudeCLI

    User->>GeekParrot: "重构 ai/utils 包"
    GeekParrot->>GeekParrot: CheckPermission()
    GeekParrot->>CCRunner: Execute(cfg, prompt)

    loop Security Check
        CCRunner->>DangerDetector: CheckInput(prompt)
        DangerDetector-->>CCRunner: Safe / Blocked
    end

    CCRunner->>ClaudeCLI: --print --session-id xxx

    loop Streaming
        ClaudeCLI-->>CCRunner: thinking/tool_use/answer
        CCRunner-->>User: 实时流式响应
    end

    CCRunner-->>GeekParrot: SessionStats
    GeekParrot-->>User: 任务完成 + 统计
```

## Parrot 接口实现

GeekParrot 实现了 `agent.ParrotAgent` 接口：

```go
type GeekParrot struct {
    runner    *agentpkg.CCRunner
    mode      *GeekMode
    sessionID string
    userID    int32
    workDir   string
}

func (p *GeekParrot) Name() string { ... }
func (p *GeekParrot) Execute(ctx, userInput, history, callback) error { ... }
func (p *GeekParrot) SelfDescribe() *agentpkg.ParrotSelfCognition { ... }
func (p *GeekParrot) GetSessionStats() *agentpkg.NormalSessionStats { ... }
```

## 配置选项

GeekParrot 通过环境变量或配置获取工作目录：

| 环境变量 | 说明 |
| :--- | :--- |
| `DIVINESENSE_CLAUDE_CODE_WORKDIR` | Claude Code 工作目录 |
| 默认值 | `~/.divinesense/claude-code` |
