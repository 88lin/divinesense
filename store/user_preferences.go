package store

// UserPreferences represents user preferences for AI personalization.
type UserPreferences struct {
	Preferences string
	CreatedTs   int64
	UpdatedTs   int64
	UserID      int32
}

// FindUserPreferences specifies the conditions for finding user preferences.
type FindUserPreferences struct {
	UserID *int32
}

// UpsertUserPreferences specifies the data for upserting user preferences.
type UpsertUserPreferences struct {
	Preferences string
	UserID      int32
}
