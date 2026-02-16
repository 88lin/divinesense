-- Rollback memo list query optimization
-- Remove indexes added for performance optimization

DROP INDEX IF EXISTS idx_memo_creator_visibility_pinned_ts;
DROP INDEX IF EXISTS idx_memo_relation_comment_memo;
