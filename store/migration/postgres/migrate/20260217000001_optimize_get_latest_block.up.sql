-- =============================================================================
-- Optimize GetLatestAIBlock query for routing latency (P0 fix)
-- =============================================================================

-- Add dedicated index for GetLatestAIBlock query
-- This query is called during routing to check sticky routing state
-- Before: ~284ms with composite index (conversation_id, status, round_number)
-- After: Expected <5ms with this dedicated index
CREATE INDEX IF NOT EXISTS idx_ai_block_conversation_round_desc
    ON ai_block(conversation_id, round_number DESC);
