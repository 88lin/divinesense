package agent

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	// Scanner buffer sizes
	// 扫描器缓冲区大小
	scannerInitialBufSize = 256 * 1024  // 256 KB
	scannerMaxBufSize     = 1024 * 1024 // 1 MB

	// Max length for non-JSON output logging
	// 非 JSON 输出日志的最大长度
	maxNonJSONOutputLength = 100
)

// buildSystemPrompt provides minimal, high-signal context for Claude Code CLI.
func buildSystemPrompt(workDir, sessionID string, userID int32, deviceContext string) string {
	osName := runtime.GOOS
	arch := runtime.GOARCH
	if osName == "darwin" {
		osName = "macOS"
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Try to parse device context for better formatting
	// 尝试解析设备上下文以便更好地格式化
	var contextMap map[string]any
	userAgent := "Unknown"
	deviceInfo := "Unknown"
	if deviceContext != "" {
		if err := json.Unmarshal([]byte(deviceContext), &contextMap); err == nil {
			if ua, ok := contextMap["userAgent"].(string); ok {
				userAgent = ua
			}
			if mobile, ok := contextMap["isMobile"].(bool); ok {
				if mobile {
					deviceInfo = "Mobile"
				} else {
					deviceInfo = "Desktop"
				}
			}
			// Add more fields if available (screen, language, etc.)
			// 如果有更多字段则添加（屏幕、语言等）
			if w, ok := contextMap["screenWidth"].(float64); ok {
				if h, ok := contextMap["screenHeight"].(float64); ok {
					deviceInfo = fmt.Sprintf("%s (%dx%d)", deviceInfo, int(w), int(h))
				}
			}
			if lang, ok := contextMap["language"].(string); ok {
				deviceInfo = fmt.Sprintf("%s, Language: %s", deviceInfo, lang)
			}
		} else {
			// Fallback: use raw string if not JSON
			userAgent = deviceContext
		}
	}

	return fmt.Sprintf(`# Context

You are running inside DivineSense, an intelligent assistant system.

**User Interaction**: Users type questions in their web browser, which invokes you via a Go backend. Your response streams back to their browser in real-time.

- **User ID**: %d
- **Client Device**: %s
- **User Agent**: %s
- **Server OS**: %s (%s)
- **Time**: %s
- **Workspace**: %s
- **Mode**: Non-interactive headless (--print)
- **Session**: %s (persists via --session-id/--resume)

---

# File Output

When you create a file, announce the filename so the user knows it was created.
`, userID, deviceInfo, userAgent, osName, arch, timestamp, workDir, sessionID)
}

// StreamMessage represents a single event in the stream-json format.
// StreamMessage 表示 stream-json 格式中的单个事件。
type StreamMessage struct {
	Type      string            `json:"type"`
	Timestamp string            `json:"timestamp,omitempty"`
	SessionID string            `json:"session_id,omitempty"`
	Role      string            `json:"role,omitempty"`
	Content   []ContentBlock    `json:"content,omitempty"`
	Message   *AssistantMessage `json:"message,omitempty"` // Nested message for "assistant" type
	Name      string            `json:"name,omitempty"`
	Input     map[string]any    `json:"input,omitempty"`
	Output    string            `json:"output,omitempty"`
	Status    string            `json:"status,omitempty"`
	Error     string            `json:"error,omitempty"`
	Duration  int               `json:"duration_ms,omitempty"`
}

// GetContentBlocks returns the content blocks, checking both direct and nested locations.
// GetContentBlocks 返回内容块，同时检查直接和嵌套位置。
func (m *StreamMessage) GetContentBlocks() []ContentBlock {
	if m.Message != nil && len(m.Message.Content) > 0 {
		return m.Message.Content
	}
	return m.Content
}

// AssistantMessage represents the nested message structure in assistant events.
// AssistantMessage 表示 assistant 事件中的嵌套消息结构。
type AssistantMessage struct {
	ID      string         `json:"id,omitempty"`
	Type    string         `json:"type,omitempty"`
	Role    string         `json:"role,omitempty"`
	Content []ContentBlock `json:"content,omitempty"`
}

// ContentBlock represents a content block in stream-json format.
// ContentBlock 表示 stream-json 格式中的内容块。
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

// GeekParrot is the Geek Mode specialist parrot (🦜 极客).
// GeekParrot 是极客模式专用鹦鹉（🦜 极客）.
// It provides DIRECT access to Claude Code CLI without any LLM processing.
// 它提供 Claude Code CLI 的直接访问，不经过任何 LLM 处理。
type GeekParrot struct {
	cliPath string
	workDir string
	userID  int32
	timeout time.Duration
	mu      sync.Mutex

	// User context
	// 用户上下文
	deviceContext string // Detailed context (JSON)

	// Session management
	// 会话管理
	sessionID  string // 会话 ID (UUID)
	firstCall  bool   // 是否首次调用
	sessionDir string // 会话目录
}

// NewGeekParrot creates a new GeekParrot instance.
// NewGeekParrot 创建一个新的 GeekParrot 实例。
func NewGeekParrot(workDir string, userID int32, sessionID string) (*GeekParrot, error) {
	cliPath, err := exec.LookPath("claude")
	if err != nil {
		return nil, fmt.Errorf("Claude Code CLI not found: %w", err)
	}

	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create work directory: %w", err)
	}

	return &GeekParrot{
		cliPath:   cliPath,
		workDir:   workDir,
		userID:    userID,
		timeout:   10 * time.Minute, // Long timeout for CLI interactions
		sessionID: sessionID,
		firstCall: true, // Default to true, adjusted in ExecuteWithCallback
	}, nil
}

// SetDeviceContext sets the full device and browser context for the parrot.
// SetDeviceContext 为鹦鹉设置完整的设备和浏览器上下文。
func (p *GeekParrot) SetDeviceContext(contextJson string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.deviceContext = contextJson
}

// Name returns the name of the parrot.
// Name 返回鹦鹉名称。
func (p *GeekParrot) Name() string {
	return "geek"
}

// ExecuteWithCallback runs Claude Code CLI with session persistence.
// ExecuteWithCallback 运行 Claude Code CLI，支持会话持久化。
func (p *GeekParrot) ExecuteWithCallback(
	ctx context.Context,
	userInput string,
	history []string, // Ignored - Claude Code manages its own history
	callback EventCallback,
) error {
	p.mu.Lock()
	firstCall := p.firstCall
	sessionID := p.sessionID
	p.mu.Unlock()

	// Determine if this is a first call or resume based on session directory existence
	// 根据会话目录是否存在决定是首次调用还是恢复
	if sessionID == "" {
		// Fallback for empty session ID (should not happen with updated handler)
		sessionID = uuid.New().String()
		p.mu.Lock()
		p.sessionID = sessionID
		p.mu.Unlock()
	}

	p.sessionDir = filepath.Join(p.workDir, ".claude", "sessions", sessionID)

	// Check if session directory exists
	// 检查会话目录是否存在
	if _, err := os.Stat(p.sessionDir); os.IsNotExist(err) {
		firstCall = true
		if err := os.MkdirAll(p.sessionDir, 0755); err != nil {
			slog.Error("GeekParrot: failed to create session directory",
				"user_id", p.userID,
				"session_id", sessionID,
				"error", err)
		}
		slog.Info("GeekParrot: Starting NEW session",
			"user_id", p.userID,
			"session_id", sessionID)
	} else {
		firstCall = false
		slog.Info("GeekParrot: Resuming EXISTING session",
			"user_id", p.userID,
			"session_id", sessionID)
	}

	p.mu.Lock()
	p.firstCall = firstCall
	p.mu.Unlock()

	slog.Info("GeekParrot: Executing Claude Code CLI",
		"user_id", p.userID,
		"session_id", sessionID,
		"first_call", firstCall,
		"input_length", len(userInput))

	// Send thinking event (use i18n key for frontend translation)
	if callback != nil {
		if err := callback(EventTypeThinking, "ai.geek_mode.thinking"); err != nil {
			return err
		}
	}

	// Execute and stream response
	if err := p.executeWithSession(ctx, userInput, firstCall, sessionID, callback); err != nil {
		// On error, reset to allow retry
		p.mu.Lock()
		p.firstCall = true
		p.sessionID = ""
		p.mu.Unlock()

		if callback != nil {
			callback(EventTypeError, err.Error())
		}
		return NewParrotError(p.Name(), "ExecuteWithCallback", err)
	}

	return nil
}

// executeWithSession executes Claude Code CLI with appropriate flags.
// executeWithSession 使用适当的标志执行 Claude Code CLI。
func (p *GeekParrot) executeWithSession(
	ctx context.Context,
	prompt string,
	firstCall bool,
	sessionID string,
	callback EventCallback,
) error {
	// Build dynamic system prompt with current context
	// 构建包含当前上下文的动态 system prompt
	p.mu.Lock()
	deviceContext := p.deviceContext
	p.mu.Unlock()
	systemPrompt := buildSystemPrompt(p.workDir, sessionID, p.userID, deviceContext)

	// Build command arguments
	// 构建命令参数
	// Note: prompt is passed as a positional argument after all flags
	// 注意：prompt 作为位置参数放在所有标志之后
	var args []string
	if firstCall {
		// First call: use --session-id
		// 首次调用：使用 --session-id
		args = []string{
			"--print",
			"--verbose",
			"--append-system-prompt", systemPrompt,
			"--session-id", sessionID,
			"--output-format", "stream-json",
			prompt, // ← 位置参数，不是选项
		}
	} else {
		// Subsequent calls: use --resume with sessionID
		// 后续调用：使用 --resume 指定会话 ID
		args = []string{
			"--print",
			"--verbose",
			"--append-system-prompt", systemPrompt,
			"--resume", sessionID,
			"--output-format", "stream-json",
			prompt, // ← 位置参数
		}
	}

	cmd := exec.CommandContext(ctx, p.cliPath, args...)
	cmd.Dir = p.workDir

	// Set environment for programmatic usage
	// 设置程序化使用环境变量
	// Note: --print flag enables headless mode; CLAUDE_HEADLESS is not documented
	// 注意：--print 标志启用无头模式；CLAUDE_HEADLESS 未在文档中
	cmd.Env = append(os.Environ(),
		"CLAUDE_DISABLE_TELEMETRY=1",
	)

	// Get pipes
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}
	defer stderr.Close()

	// Start command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start command: %w", err)
	}

	// Stream output with timeout
	// 带超时流式输出输出
	if err := p.streamOutput(ctx, stdout, stderr, callback); err != nil {
		// Try to kill the process
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return err
	}

	// Wait for command completion
	waitErr := cmd.Wait()

	// Check for specific exit codes
	if waitErr != nil {
		exitCode := 0
		if cmd.ProcessState != nil {
			exitCode = cmd.ProcessState.ExitCode()
		}
		return fmt.Errorf("command exited with code %d: %w", exitCode, waitErr)
	}

	return nil
}

// streamOutput reads and parses stream-json output from CLI.
// streamOutput 读取并解析 CLI 的 stream-json 输出。
func (p *GeekParrot) streamOutput(
	ctx context.Context,
	stdout, stderr io.ReadCloser,
	callback EventCallback,
) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	done := make(chan struct{})

	// Stream stdout
	// 流式处理 stdout
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		buf := make([]byte, 0, scannerInitialBufSize)
		scanner.Buffer(buf, scannerMaxBufSize)

		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			var msg StreamMessage
			if err := json.Unmarshal([]byte(line), &msg); err != nil {
				// Not JSON, treat as plain text
				if len(line) > maxNonJSONOutputLength {
					line = line[:maxNonJSONOutputLength]
				}
				slog.Debug("GeekParrot: non-JSON output",
					"user_id", p.userID,
					"line", line)
				if callback != nil {
					callback(EventTypeAnswer, line)
				}
				continue
			}

			// Log each message type for debugging (debug level to reduce noise)
			slog.Debug("GeekParrot: received message",
				"user_id", p.userID,
				"type", msg.Type,
				"content_blocks", len(msg.Content),
				"has_name", msg.Name != "",
				"has_output", msg.Output != "",
				"has_error", msg.Error != "")

			// Dispatch event to callback
			if callback != nil {
				if err := p.dispatchCallback(msg, callback); err != nil {
					errCh <- err
					return
				}
			}

			// Check for completion
			if msg.Type == "result" || msg.Type == "error" {
				return
			}
		}
		errCh <- scanner.Err()
	}()

	// Stream stderr to log
	// 流式处理 stderr 到日志
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			slog.Warn("GeekParrot: stderr from Claude Code CLI",
				"user_id", p.userID,
				"line", scanner.Text())
		}
		errCh <- scanner.Err()
	}()

	// Wait for completion or timeout
	// 等待完成或超时
	go func() {
		wg.Wait()
		close(done)
	}()

	// Use timer instead of time.After to avoid goroutine leak
	timer := time.NewTimer(p.timeout)
	defer timer.Stop()

	select {
	case <-done:
		// Collect any errors
		var errors []string
		for i := 0; i < 2; i++ {
			select {
			case err := <-errCh:
				if err != nil {
					errors = append(errors, err.Error())
				}
			default:
			}
		}
		if len(errors) > 0 {
			return fmt.Errorf("stream errors: %s", errors[0])
		}
		return nil
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	case <-timer.C:
		return fmt.Errorf("execution timeout after %v", p.timeout)
	}
}

// dispatchCallback dispatches stream events to the callback.
// dispatchCallback 将流事件分发给回调。
func (p *GeekParrot) dispatchCallback(msg StreamMessage, callback EventCallback) error {
	switch msg.Type {
	case "error":
		if msg.Error != "" {
			slog.Debug("GeekParrot: dispatching error",
				"user_id", p.userID,
				"error", msg.Error)
			return callback(EventTypeError, msg.Error)
		}
		// Empty error case - ignore per stream-json spec
	case "thinking", "status":
		for _, block := range msg.GetContentBlocks() {
			if block.Type == "text" && block.Text != "" {
				slog.Debug("GeekParrot: dispatching thinking",
					"user_id", p.userID,
					"text_len", len(block.Text))
				if err := callback(EventTypeThinking, block.Text); err != nil {
					return err
				}
			}
		}
	case "tool_use":
		if msg.Name != "" {
			slog.Debug("GeekParrot: dispatching tool_use",
				"user_id", p.userID,
				"name", msg.Name)
			if err := callback(EventTypeToolUse, msg.Name); err != nil {
				return err
			}
		}
	case "tool_result":
		if msg.Output != "" {
			slog.Debug("GeekParrot: dispatching tool_result",
				"user_id", p.userID,
				"output_len", len(msg.Output))
			if err := callback(EventTypeToolResult, msg.Output); err != nil {
				return err
			}
		}
	case "message", "content", "text", "delta":
		for _, block := range msg.GetContentBlocks() {
			if block.Type == "text" && block.Text != "" {
				slog.Debug("GeekParrot: dispatching answer",
					"user_id", p.userID,
					"text_len", len(block.Text))
				if err := callback(EventTypeAnswer, block.Text); err != nil {
					return err
				}
			}
		}
	case "assistant":
		// Assistant type has nested message.content structure
		for _, block := range msg.GetContentBlocks() {
			if block.Type == "text" && block.Text != "" {
				slog.Debug("GeekParrot: dispatching assistant answer",
					"user_id", p.userID,
					"text_len", len(block.Text))
				if err := callback(EventTypeAnswer, block.Text); err != nil {
					return err
				}
			}
		}
	default:
		// Try to extract any text content from both locations
		for _, block := range msg.GetContentBlocks() {
			if block.Type == "text" && block.Text != "" {
				slog.Debug("GeekParrot: dispatching default answer",
					"user_id", p.userID,
					"text_len", len(block.Text))
				callback(EventTypeAnswer, block.Text)
			}
		}
	}
	return nil
}

// ResetSession resets the session state (e.g., on error or user request).
// ResetSession 重置会话状态（例如出错或用户请求时）。
func (p *GeekParrot) ResetSession() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.firstCall = true
	p.sessionID = ""
	p.sessionDir = ""
	slog.Info("GeekParrot: Session reset",
		"user_id", p.userID)
}

// GetSessionID returns the current session ID.
// GetSessionID 返回当前会话 ID。
func (p *GeekParrot) GetSessionID() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.sessionID
}

// SelfDescribe returns the GeekParrot's metacognitive information.
// SelfDescribe 返回极客鹦鹉的元认知信息。
func (p *GeekParrot) SelfDescribe() *ParrotSelfCognition {
	return &ParrotSelfCognition{
		Name:  "geek",
		Emoji: "🦜",
		Title: "Claude Code CLI Runner",
		Capabilities: []string{
			"调用 Claude Code CLI",
			"通过 Go runner 执行",
			"服务 Web 界面用户",
			"实时流式响应",
		},
		Limitations: []string{
			"需要安装 Claude Code CLI",
			"Headless 模式运行",
		},
		WorkingStyle: "Go backend → Claude Code CLI → Web 用户",
	}
}

// IsSessionActive returns whether a session has been started.
// IsSessionActive 返回是否已启动会话。
func (p *GeekParrot) IsSessionActive() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return !p.firstCall
}

// GetWorkDir returns the working directory for Claude Code CLI.
// GetWorkDir 返回 Claude Code CLI 的工作目录。
func (p *GeekParrot) GetWorkDir() string {
	return p.workDir
}

// GetUserID returns the user ID for this parrot.
// GetUserID 返回此鹦鹉的用户 ID。
func (p *GeekParrot) GetUserID() int32 {
	return p.userID
}

// Cancel is a no-op for --continue mode (each call is independent).
// Cancel 对 --continue 模式是空操作（每次调用独立）。
func (p *GeekParrot) Cancel() {
	// No-op in --continue mode
	// Each execution is independent
}
