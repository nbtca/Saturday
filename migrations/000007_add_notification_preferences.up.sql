-- Add notification_preferences column to member table
ALTER TABLE public.member
ADD COLUMN IF NOT EXISTS notification_preferences JSONB DEFAULT '{"new_event_created": false, "event_assigned_to_me": true, "event_status_changed": true}'::jsonb;

-- Add index for faster JSONB queries
CREATE INDEX IF NOT EXISTS idx_member_notification_preferences
ON public.member USING GIN (notification_preferences);

-- Comment on the column
COMMENT ON COLUMN public.member.notification_preferences IS '成员通知偏好设置 (Member notification preferences)';

-- Update member_view to include notification_preferences
CREATE OR REPLACE VIEW public.member_view AS
 SELECT member.member_id,
    member.alias,
    member.password,
    member.name,
    member.section,
    member.profile,
    member.phone,
    member.qq,
    member.avatar,
    member.created_by,
    member.gmt_create,
    member.gmt_modified,
    COALESCE(role.role, ''::character varying) AS role,
    member.logto_id,
    member.github_id,
    member.notification_preferences
   FROM ((public.member
     LEFT JOIN public.member_role_relation ON ((member.member_id = (member_role_relation.member_id)::bpchar)))
     LEFT JOIN public.role ON ((member_role_relation.role_id = role.role_id)));

