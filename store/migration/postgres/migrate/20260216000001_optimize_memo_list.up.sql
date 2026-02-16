-- Optimize memo list query performance
-- Add composite indexes for creator_id + visibility + pinned + created_ts
--
-- This migration adds indexes to improve memo list query performance:
-- - idx_memo_creator_visibility_pinned_ts: Composite index for common query patterns
-- - idx_memo_relation_comment: Partial index for comment filtering

-- =====================================================
-- UP: Add indexes
-- =====================================================

-- Composite index: supports creator_id + visibility filter + pinned/created_ts sort
CREATE INDEX IF NOT EXISTS idx_memo_creator_visibility_pinned_ts
ON memo (creator_id, visibility, pinned DESC, created_ts DESC);

-- Partial index: for memo_relation comment filtering
CREATE INDEX IF NOT EXISTS idx_memo_relation_comment_memo
ON memo_relation (memo_id, related_memo_id)
WHERE type = 'COMMENT';

-- =====================================================
-- DOWN: Remove indexes
-- =====================================================

DROP INDEX IF EXISTS idx_memo_creator_visibility_pinned_ts;
DROP INDEX IF EXISTS idx_memo_relation_comment_memo;
