-- Restore member_view to original state (without notification_preferences)
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
    member.github_id
   FROM ((public.member
     LEFT JOIN public.member_role_relation ON ((member.member_id = (member_role_relation.member_id)::bpchar)))
     LEFT JOIN public.role ON ((member_role_relation.role_id = role.role_id)));

-- Remove index
DROP INDEX IF EXISTS idx_member_notification_preferences;

-- Remove notification_preferences column from member table
ALTER TABLE public.member
DROP COLUMN IF EXISTS notification_preferences;
