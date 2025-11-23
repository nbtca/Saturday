-- Remove index
DROP INDEX IF EXISTS idx_member_notification_preferences;

-- Remove notification_preferences column from member table
ALTER TABLE public.member
DROP COLUMN IF EXISTS notification_preferences;
