-- Add title_source field to ai_conversation table
-- Tracks whether the title was user-provided, auto-generated, or default
-- Issue #88: Conversation smart rename feature

-- Add title_source column (default = 'default' for existing rows)
ALTER TABLE ai_conversation
ADD COLUMN title_source TEXT NOT NULL DEFAULT 'default';

-- Add check constraint for valid title_source values
ALTER TABLE ai_conversation
ADD CONSTRAINT chk_ai_conversation_title_source
CHECK (title_source IN ('default', 'auto', 'user'));

-- Create index on title_source for filtering conversations by source type
CREATE INDEX idx_ai_conversation_title_source ON ai_conversation(title_source);

-- Update existing conversations:
-- - Empty titles -> 'default'
-- - Non-empty titles that look like user edits -> 'user'
-- - All others remain 'default'
UPDATE ai_conversation
SET title_source = CASE
    WHEN title = '' OR title IS NULL THEN 'default'
    WHEN title LIKE 'New Chat%' THEN 'default'
    ELSE 'user'
END;
