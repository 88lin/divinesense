package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lithammer/shortuuid/v4"

	"github.com/hrygo/divinesense/ai/timeout"
	"github.com/hrygo/divinesense/store"
)

// MemoWriteTool creates new memos for saving ideas and inspirations.
// MemoWriteTool 创建新笔记，用于保存创意和灵感。
type MemoWriteTool struct {
	store        *store.Store
	userIDGetter func(ctx context.Context) int32
}

// NewMemoWriteTool creates a new memo write tool.
func NewMemoWriteTool(
	st *store.Store,
	userIDGetter func(ctx context.Context) int32,
) (*MemoWriteTool, error) {
	if st == nil {
		return nil, fmt.Errorf("store cannot be nil")
	}
	if userIDGetter == nil {
		return nil, fmt.Errorf("userIDGetter cannot be nil")
	}

	return &MemoWriteTool{
		store:        st,
		userIDGetter: userIDGetter,
	}, nil
}

// Name returns the name of the tool.
func (t *MemoWriteTool) Name() string {
	return "memo_write"
}

// Description returns a description of what the tool does.
func (t *MemoWriteTool) Description() string {
	return `Create a new memo to save ideas, inspirations, or notes.

Use this tool when:
- User wants to save a creative idea or inspiration
- User asks you to record something for later reference
- User says phrases like "save this idea", "write this down", "record this"

INPUT FORMAT:
{"content": "The content to save", "title": "Optional title"}

OUTPUT:
- Success: "✓ Memo created: [UID] - [Title/First line]"
- Error: "Error: [error message]"

NOTE: This tool is for Ideation Agent to save creative outputs.`
}

// MemoWriteInput represents the input for memo creation.
type MemoWriteInput struct {
	Content string `json:"content"`
	Title   string `json:"title,omitempty"`
}

// InputType returns the JSON schema for the input.
func (t *MemoWriteTool) InputType() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"content": map[string]interface{}{
				"type":        "string",
				"description": "The content to save as a memo (required, non-empty)",
			},
			"title": map[string]interface{}{
				"type":        "string",
				"description": "Optional title for the memo. If provided, it will be added as a heading.",
			},
		},
		"required": []string{"content"},
	}
}

// Run executes the memo write tool.
func (t *MemoWriteTool) Run(ctx context.Context, input string) (string, error) {
	// Add timeout protection
	ctx, cancel := context.WithTimeout(ctx, timeout.ToolExecutionTimeout)
	defer cancel()

	// Parse input
	var writeInput MemoWriteInput
	if err := json.Unmarshal([]byte(input), &writeInput); err != nil {
		return "", fmt.Errorf("invalid JSON input: %w", err)
	}

	// Validate content
	content := strings.TrimSpace(writeInput.Content)
	if content == "" {
		return "", fmt.Errorf("content is required and cannot be empty")
	}

	// Build full content with optional title
	if writeInput.Title != "" {
		title := strings.TrimSpace(writeInput.Title)
		content = fmt.Sprintf("# %s\n\n%s", title, content)
	}

	// Get user ID
	userID := t.userIDGetter(ctx)
	if userID == 0 {
		return "", fmt.Errorf("user not authenticated")
	}

	// Generate UID
	memoUID := shortuuid.New()

	// Create memo
	create := &store.Memo{
		UID:        memoUID,
		CreatorID:  userID,
		Content:    content,
		Visibility: store.Private, // Default to private for ideation
		CreatedTs:  time.Now().Unix(),
		UpdatedTs:  time.Now().Unix(),
	}

	// Save to store
	memo, err := t.store.CreateMemo(ctx, create)
	if err != nil {
		return "", fmt.Errorf("failed to create memo: %w", err)
	}

	// Extract display title (first line or title)
	displayTitle := writeInput.Title
	if displayTitle == "" {
		// Extract first line as title
		lines := strings.Split(strings.TrimSpace(writeInput.Content), "\n")
		if len(lines) > 0 {
			displayTitle = strings.TrimSpace(lines[0])
			if len(displayTitle) > 50 {
				displayTitle = displayTitle[:47] + "..."
			}
		}
	}

	return fmt.Sprintf("✓ Memo created: %s - %s", memo.UID, displayTitle), nil
}

// ensure MemoWriteTool implements required interfaces
var _ interface {
	Name() string
	Description() string
	Run(context.Context, string) (string, error)
} = (*MemoWriteTool)(nil)
