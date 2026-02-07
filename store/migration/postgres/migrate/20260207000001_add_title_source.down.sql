-- Rollback: Remove title_source field from ai_conversation table
-- Issue #88: Conversation smart rename feature

DROP INDEX IF EXISTS idx_ai_conversation_title_source;

ALTER TABLE ai_conversation
DROP CONSTRAINT IF EXISTS chk_ai_conversation_title_source;

ALTER TABLE ai_conversation
DROP COLUMN IF EXISTS title_source;
