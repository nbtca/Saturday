-- Add notification_preferences column to member table
ALTER TABLE public.member
ADD COLUMN IF NOT EXISTS notification_preferences JSONB DEFAULT '{"new_event_created": false, "event_assigned_to_me": true}'::jsonb;

-- Add index for faster JSONB queries
CREATE INDEX IF NOT EXISTS idx_member_notification_preferences
ON public.member USING GIN (notification_preferences);

-- Comment on the column
COMMENT ON COLUMN public.member.notification_preferences IS '成员通知偏好设置 (Member notification preferences)';
