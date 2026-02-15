package orchestrator

import (
	"reflect"
	"sort"
	"testing"

	agents "github.com/hrygo/divinesense/ai/agents"
)

// ... existing tests ...

func TestCapabilityMap_IdentifyCapabilities(t *testing.T) {
	cm := NewCapabilityMap()

	// Mock configs
	configs := []*agents.ParrotSelfCognition{
		{
			Name:         "Schedule",
			Capabilities: []string{"SCHEDULE_MGMT"},
			CapabilityTriggers: map[string][]string{
				"SCHEDULE_MGMT": {"schedule", "calendar", "日程", "会议"},
			},
		},
		{
			Name:         "Memo",
			Capabilities: []string{"NOTE_SEARCH"},
			CapabilityTriggers: map[string][]string{
				"NOTE_SEARCH": {"note", "memo", "笔记"},
			},
		},
	}
	cm.BuildFromConfigs(configs)

	tests := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Exact word match (ASCII)",
			text:     "I want to check my schedule",
			expected: []string{"schedule_mgmt"},
		},
		{
			name:     "Partial match fail (ASCII)",
			text:     "This is unrelated to scheduler",
			expected: nil,
		},
		{
			name:     "Non-ASCII containment (Chinese)",
			text:     "帮我创建日程",
			expected: []string{"schedule_mgmt"},
		},
		{
			name:     "Trigger with boundaries",
			text:     "Take a note.",
			expected: []string{"note_search"},
		},
		{
			name:     "Trigger inside other word",
			text:     "keynote speaker",
			expected: nil, // 'note' inside 'keynote' shouldn't match
		},
		{
			name:     "Multiple matches",
			text:     "Check schedule and take a note",
			expected: []string{"schedule_mgmt", "note_search"},
		},
		{
			name:     "Case insensitive",
			text:     "CALENDAR access",
			expected: []string{"schedule_mgmt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cm.IdentifyCapabilities(tt.text)

			// Sort for comparison
			sort.Strings(got)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(got, tt.expected) {
				// Handle nil vs empty slice
				if len(got) == 0 && len(tt.expected) == 0 {
					return
				}
				t.Errorf("IdentifyCapabilities() = %v, want %v", got, tt.expected)
			}
		})
	}
}
