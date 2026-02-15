// Package orchestrator implements the Orchestrator-Workers pattern for multi-agent coordination.
// It uses LLM to dynamically decompose tasks, dispatch to expert agents, and aggregate results.
package orchestrator

import (
	"strings"
	"sync"

	agents "github.com/hrygo/divinesense/ai/agents"
)

// Capability represents a single capability that an expert agent can provide.
// Capability represents a single capability that an expert agent can provide.
type Capability string

// ExpertInfo contains information about an expert agent.
// ExpertInfo 包含专家代理的信息。
type ExpertInfo struct {
	Name         string   `json:"name"`
	Emoji        string   `json:"emoji"`
	Title        string   `json:"title"`
	Capabilities []string `json:"capabilities"`
}

// CapabilityMap provides a thread-safe mapping from capabilities to expert agents.
// It is used at runtime to build the capability-to-expert mapping.
type CapabilityMap struct {
	mu                    sync.RWMutex
	capabilityToExperts   map[Capability][]*ExpertInfo
	keywordToCapabilities map[string][]Capability
	experts               map[string]*ExpertInfo
}

// NewCapabilityMap creates an empty CapabilityMap.
// NewCapabilityMap 创建一个空的 CapabilityMap。
func NewCapabilityMap() *CapabilityMap {
	return &CapabilityMap{
		capabilityToExperts:   make(map[Capability][]*ExpertInfo),
		keywordToCapabilities: make(map[string][]Capability),
		experts:               make(map[string]*ExpertInfo),
	}
}

// BuildFromConfigs builds the capability map from ParrotSelfCognition configurations.
// BuildFromConfigs 从 ParrotSelfCognition 配置构建能力映射。
func (cm *CapabilityMap) BuildFromConfigs(configs []*agents.ParrotSelfCognition) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Clear existing mappings
	cm.capabilityToExperts = make(map[Capability][]*ExpertInfo)
	cm.keywordToCapabilities = make(map[string][]Capability)
	cm.experts = make(map[string]*ExpertInfo)

	for _, config := range configs {
		if config == nil {
			continue
		}

		expert := &ExpertInfo{
			Name:         config.Name,
			Emoji:        config.Emoji,
			Title:        config.Title,
			Capabilities: config.Capabilities,
		}

		// Add to experts map
		cm.experts[config.Name] = expert

		// Add to capabilityToExperts map
		for _, cap := range config.Capabilities {
			normalizedCap := cm.normalizeCapability(cap)
			if normalizedCap == "" {
				continue
			}
			cm.capabilityToExperts[Capability(normalizedCap)] = append(
				cm.capabilityToExperts[Capability(normalizedCap)],
				expert,
			)
		}

		// Add to keywordToCapabilities map
		for capName, triggers := range config.CapabilityTriggers {
			normalizedCap := cm.normalizeCapability(capName)
			if normalizedCap == "" {
				continue
			}

			for _, trigger := range triggers {
				normalizedTrigger := cm.normalizeCapability(trigger) // Reuse normalization logic
				if normalizedTrigger == "" {
					continue
				}
				cm.keywordToCapabilities[normalizedTrigger] = append(
					cm.keywordToCapabilities[normalizedTrigger],
					Capability(normalizedCap),
				)
			}
		}
	}
}

// FindExpertsByCapability returns all experts that provide the given capability.
// FindExpertsByCapability 返回提供指定能力的所有专家。
func (cm *CapabilityMap) FindExpertsByCapability(capability string) []*ExpertInfo {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	normalizedCap := cm.normalizeCapability(capability)
	if normalizedCap == "" {
		return nil
	}

	experts, ok := cm.capabilityToExperts[Capability(normalizedCap)]
	if !ok {
		return nil
	}

	// Return a copy to avoid external mutation
	result := make([]*ExpertInfo, len(experts))
	copy(result, experts)
	return result
}

// FindAlternativeExperts returns all experts that provide the given capability,
// excluding the specified expert. This is useful for finding fallback experts.
// FindAlternativeExperts 返回提供指定能力的所有专家，但排除指定的专家。
// 这在寻找备用专家时很有用。
func (cm *CapabilityMap) FindAlternativeExperts(capability string, excludeExpert string) []*ExpertInfo {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	normalizedCap := cm.normalizeCapability(capability)
	if normalizedCap == "" {
		return nil
	}

	experts, ok := cm.capabilityToExperts[Capability(normalizedCap)]
	if !ok {
		return nil
	}

	// Filter out the excluded expert
	var result []*ExpertInfo
	for _, expert := range experts {
		if expert.Name != excludeExpert {
			result = append(result, expert)
		}
	}

	return result
}

// IdentifyCapabilities identifies capabilities from a text based on registered triggers.
// IdentifyCapabilities 根据注册的触发器从文本中识别能力。
func (cm *CapabilityMap) IdentifyCapabilities(text string) []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	normalizedText := strings.ToLower(text)
	var matchedCapabilities []string
	seen := make(map[string]bool)

	for trigger, caps := range cm.keywordToCapabilities {
		if cm.matchesTrigger(normalizedText, trigger) {
			for _, cap := range caps {
				capStr := string(cap)
				if !seen[capStr] {
					seen[capStr] = true
					matchedCapabilities = append(matchedCapabilities, capStr)
				}
			}
		}
	}

	return matchedCapabilities
}

// matchesTrigger checks if the trigger exists in the text.
// For ASCII triggers, it enforces word boundaries to avoid partial matches.
// For non-ASCII triggers (e.g. Chinese), it uses simple containment.
func (cm *CapabilityMap) matchesTrigger(text, trigger string) bool {
	// 1. Basic containment check
	idx := strings.Index(text, trigger)
	if idx == -1 {
		return false
	}

	// 2. If trigger contains non-ASCII characters (e.g. Chinese), containment is sufficient
	if isNonASCII(trigger) {
		return true
	}

	// 3. For ASCII triggers, verify word boundaries
	// We must check all occurrences
	for idx != -1 {
		// Check left boundary
		leftOk := (idx == 0) || !isWordChar(text[idx-1])

		// Check right boundary
		end := idx + len(trigger)
		rightOk := (end == len(text)) || !isWordChar(text[end])

		if leftOk && rightOk {
			return true
		}

		// Find next occurrence
		next := strings.Index(text[idx+1:], trigger)
		if next == -1 {
			break
		}
		idx += 1 + next
	}

	return false
}

func isNonASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 128 {
			return true
		}
	}
	return false
}

func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

// GetAllExperts returns all registered experts.
// GetAllExperts 返回所有已注册的专家。
func (cm *CapabilityMap) GetAllExperts() []*ExpertInfo {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	experts := make([]*ExpertInfo, 0, len(cm.experts))
	for _, expert := range cm.experts {
		experts = append(experts, expert)
	}
	return experts
}

// normalizeCapability normalizes a capability string for consistent lookup.
// It converts the capability to lowercase and trims whitespace.
// normalizeCapability 标准化能力字符串以进行一致的查找。
// 它将能力转换为小写并去除空白。
func (cm *CapabilityMap) normalizeCapability(cap string) string {
	return strings.ToLower(strings.TrimSpace(cap))
}
