package store

type AIConversation struct {
	UID       string
	Title     string
	ParrotID  string
	RowStatus RowStatus
	CreatedTs int64
	UpdatedTs int64
	ID        int32
	CreatorID int32
	Pinned    bool
}

type FindAIConversation struct {
	ID        *int32
	UID       *string
	CreatorID *int32
	Pinned    *bool
	RowStatus *RowStatus
}

type UpdateAIConversation struct {
	Title     *string
	ParrotID  *string
	Pinned    *bool
	RowStatus *RowStatus
	UpdatedTs *int64
	ID        int32
}

type DeleteAIConversation struct {
	ID int32
}

type AIMessageRole string

const (
	AIMessageRoleUser      AIMessageRole = "USER"
	AIMessageRoleAssistant AIMessageRole = "ASSISTANT"
	AIMessageRoleSystem    AIMessageRole = "SYSTEM"
)

type AIMessageType string

const (
	AIMessageTypeMessage   AIMessageType = "MESSAGE"
	AIMessageTypeSeparator AIMessageType = "SEPARATOR"
	AIMessageTypeSummary   AIMessageType = "SUMMARY" // Conversation summary (invisible to frontend)
)

type AIMessage struct {
	UID            string
	Type           AIMessageType
	Role           AIMessageRole
	Content        string
	Metadata       string
	CreatedTs      int64
	ID             int32
	ConversationID int32
}

type FindAIMessage struct {
	ID             *int32
	UID            *string
	ConversationID *int32
}

type DeleteAIMessage struct {
	ID             *int32
	ConversationID *int32
}
