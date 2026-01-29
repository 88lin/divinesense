package store

import "time"

// EpisodicMemory represents an episodic memory record for AI learning.
type EpisodicMemory struct {
	Timestamp  time.Time
	AgentType  string
	UserInput  string
	Outcome    string
	Summary    string
	ID         int64
	CreatedTs  int64
	UserID     int32
	Importance float32
}

// FindEpisodicMemory specifies the conditions for finding episodic memories.
type FindEpisodicMemory struct {
	ID        *int64
	UserID    *int32
	AgentType *string
	Query     *string // For text search in user_input and summary
	Limit     int
	Offset    int
}

// DeleteEpisodicMemory specifies the conditions for deleting episodic memories.
type DeleteEpisodicMemory struct {
	ID     *int64
	UserID *int32
}
