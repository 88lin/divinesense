// Package context provides context building for LLM prompts.
package context

// Default token budget values.
//
// Token allocation strategy (with retrieval):
// - System Prompt: 500 tokens (fixed)
// - User Preferences: 10%
// - Short-term Memory: 40%
// - Long-term Memory: 15%
// - Retrieval Results: 45%
//
// Reference: docs/agent-engineering/AGENT_PROMPT_BEST_PRACTICES.md
const (
	DefaultMaxTokens      = 4096
	DefaultSystemPrompt   = 500
	DefaultUserPrefsRatio = 0.10
	DefaultRetrievalRatio = 0.45 // Updated per best practices
	MinSegmentTokens      = 100
)

// TokenBudget represents the token allocation plan.
type TokenBudget struct {
	Total           int
	SystemPrompt    int
	ShortTermMemory int
	LongTermMemory  int
	Retrieval       int
	UserPrefs       int
}

// BudgetAllocator allocates token budgets.
type BudgetAllocator struct {
	systemPromptTokens int
	userPrefsRatio     float64
	retrievalRatio     float64
	profileRegistry    *ProfileRegistry
	intentResolver     *IntentResolver
}

// NewBudgetAllocator creates a new budget allocator with defaults.
func NewBudgetAllocator() *BudgetAllocator {
	return &BudgetAllocator{
		systemPromptTokens: DefaultSystemPrompt,
		userPrefsRatio:     DefaultUserPrefsRatio,
		retrievalRatio:     DefaultRetrievalRatio,
		profileRegistry:    NewProfileRegistry(),
		intentResolver:     NewIntentResolver(),
	}
}

// WithProfileRegistry sets a custom profile registry.
func (a *BudgetAllocator) WithProfileRegistry(registry *ProfileRegistry) *BudgetAllocator {
	a.profileRegistry = registry
	return a
}

// WithIntentResolver sets a custom intent resolver.
func (a *BudgetAllocator) WithIntentResolver(resolver *IntentResolver) *BudgetAllocator {
	a.intentResolver = resolver
	return a
}

// LoadProfileFromEnv loads profile overrides from environment variables.
func (a *BudgetAllocator) LoadProfileFromEnv() {
	a.profileRegistry.LoadFromEnv()
}

// Allocate allocates token budget based on total and whether retrieval is needed.
func (a *BudgetAllocator) Allocate(total int, hasRetrieval bool) *TokenBudget {
	if total <= 0 {
		total = DefaultMaxTokens
	}

	budget := &TokenBudget{
		Total:        total,
		SystemPrompt: a.systemPromptTokens,
		UserPrefs:    int(float64(total) * a.userPrefsRatio),
	}

	remaining := total - budget.SystemPrompt - budget.UserPrefs

	if hasRetrieval {
		// With retrieval: prioritize retrieval context
		// Short-term: 40%, Long-term: 15%, Retrieval: 45%
		budget.ShortTermMemory = int(float64(remaining) * 0.40)
		budget.LongTermMemory = int(float64(remaining) * 0.15)
		budget.Retrieval = int(float64(remaining) * 0.45)
	} else {
		// No retrieval: more space for memory
		// Short-term: 55%, Long-term: 30%, Retrieval: 0%
		budget.ShortTermMemory = int(float64(remaining) * 0.55)
		budget.LongTermMemory = int(float64(remaining) * 0.30)
		budget.Retrieval = 0
	}

	return budget
}

// AllocateForAgent allocates token budget based on agent type (profile-based).
// Issue #93: 动态 Token 预算分配 - 按意图类型自适应调整
func (a *BudgetAllocator) AllocateForAgent(total int, hasRetrieval bool, agentType string) *TokenBudget {
	if total <= 0 {
		total = DefaultMaxTokens
	}

	// Special cases: GEEK and EVOLUTION use Claude Code CLI with their own context management
	if agentType == "GEEK" || agentType == "EVOLUTION" {
		// These modes use Claude Code CLI, no LLM budget needed
		return &TokenBudget{
			Total:           total,
			SystemPrompt:    0,
			UserPrefs:       0,
			ShortTermMemory: 0,
			LongTermMemory:  0,
			Retrieval:       0,
		}
	}

	// Resolve intent from agent type
	intent := a.intentResolver.Resolve(agentType)
	profile, found := a.profileRegistry.Get(intent)
	if !found {
		// Log warning: using fallback profile
		// profile.Get() already handles fallback, so we still have a valid profile
	}

	budget := &TokenBudget{
		Total:        total,
		SystemPrompt: a.systemPromptTokens,
		UserPrefs:    int(float64(total) * profile.UserPrefsRatio),
	}

	remaining := total - budget.SystemPrompt - budget.UserPrefs

	// Use profile ratios for remaining budget
	budget.ShortTermMemory = int(float64(remaining) * profile.ShortTermRatio)
	budget.LongTermMemory = int(float64(remaining) * profile.LongTermRatio)

	if hasRetrieval {
		budget.Retrieval = int(float64(remaining) * profile.RetrievalRatio)
	} else {
		budget.Retrieval = 0
		// Redistribute retrieval budget to short-term and long-term
		extra := int(float64(remaining) * profile.RetrievalRatio)
		totalRatio := profile.ShortTermRatio + profile.LongTermRatio
		if totalRatio > 0 {
			budget.ShortTermMemory += int(float64(extra) * (profile.ShortTermRatio / totalRatio))
			budget.LongTermMemory += int(float64(extra) * (profile.LongTermRatio / totalRatio))
		}
	}

	return budget
}

// AllocateBudget is a convenience function.
func AllocateBudget(total int, hasRetrieval bool) *TokenBudget {
	return NewBudgetAllocator().Allocate(total, hasRetrieval)
}
