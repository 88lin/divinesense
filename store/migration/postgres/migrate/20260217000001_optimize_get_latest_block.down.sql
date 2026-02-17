-- =============================================================================
-- Rollback: Optimize GetLatestAIBlock query for routing latency
-- =============================================================================

DROP INDEX IF EXISTS idx_ai_block_conversation_round_desc;
