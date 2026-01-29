package schedule

import (
	"context"
	"time"

	"github.com/hrygo/divinesense/store"
)

// Service defines the core business logic interface for schedule management.
// This abstraction allows the Agent tools to call business logic directly
// without internal HTTP callbacks or code duplication.
type Service interface {
	// FindSchedules returns schedules between start and end time.
	// For recurring schedules, this method expands instances within the time range.
	FindSchedules(ctx context.Context, userID int32, start, end time.Time) ([]*ScheduleInstance, error)

	// CreateSchedule creates a new schedule with validation logic.
	// This includes conflict checking, time range validation, and permission checks.
	CreateSchedule(ctx context.Context, userID int32, create *CreateScheduleRequest) (*store.Schedule, error)

	// UpdateSchedule updates an existing schedule.
	UpdateSchedule(ctx context.Context, userID int32, id int32, update *UpdateScheduleRequest) (*store.Schedule, error)

	// DeleteSchedule deletes a schedule by ID.
	DeleteSchedule(ctx context.Context, userID int32, id int32) error

	// CheckConflicts checks for schedule conflicts within a time range.
	// Returns a list of conflicting schedules.
	CheckConflicts(ctx context.Context, userID int32, startTs int64, endTs *int64, excludeIDs []int32) ([]*store.Schedule, error)
}

// ScheduleInstance represents a specific schedule instance (expanded from recurring schedules).
type ScheduleInstance struct {
	EndTs       *int64
	UID         string
	Title       string
	Description string
	Location    string
	Timezone    string
	ParentUID   string
	StartTs     int64
	ID          int32
	AllDay      bool
	IsRecurring bool
}

// CreateScheduleRequest represents the request to create a schedule.
type CreateScheduleRequest struct {
	EndTs           *int64
	RecurrenceRule  *string
	RecurrenceEndTs *int64
	Title           string
	Description     string
	Location        string
	Timezone        string
	Reminders       []*Reminder
	StartTs         int64
	AllDay          bool
}

// UpdateScheduleRequest represents the request to update a schedule.
type UpdateScheduleRequest struct {
	Title           *string
	Description     *string
	Location        *string
	StartTs         *int64
	EndTs           *int64
	AllDay          *bool
	Timezone        *string
	RecurrenceRule  *string
	RecurrenceEndTs *int64
	RowStatus       *store.RowStatus
	Reminders       []*Reminder
}

// Reminder represents a schedule reminder.
type Reminder struct {
	Type  string
	Unit  string
	Value int32
}
