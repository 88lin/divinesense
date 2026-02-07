package store

// TitleSource indicates how the conversation title was created.
// - "default": System default (e.g., "New Chat" or truncated first message)
// - "auto": AI-generated title based on conversation content
// - "user": User-provided title (manual edit)
type TitleSource string

const (
	TitleSourceDefault TitleSource = "default"
	TitleSourceAuto    TitleSource = "auto"
	TitleSourceUser    TitleSource = "user"
)

type AIConversation struct {
	UID         string
	Title       string
	TitleSource TitleSource // Indicates how the title was created
	ParrotID    string
	RowStatus   RowStatus
	CreatedTs   int64
	UpdatedTs   int64
	ID          int32
	CreatorID   int32
	Pinned      bool
	BlockCount  int32 // Number of blocks in this conversation (populated by ListAIConversations with JOIN)
}

type FindAIConversation struct {
	ID        *int32
	UID       *string
	CreatorID *int32
	Pinned    *bool
	RowStatus *RowStatus
}

type UpdateAIConversation struct {
	Title       *string
	TitleSource *TitleSource
	ParrotID    *string
	Pinned      *bool
	RowStatus   *RowStatus
	UpdatedTs   *int64
	ID          int32
}

type DeleteAIConversation struct {
	ID int32
}

// AIMessage types removed: ALL IN Block!
// Use AIBlock from ai_block.go instead for all conversation persistence.
