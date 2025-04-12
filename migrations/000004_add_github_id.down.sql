-- Remove the 'github_id' column from the 'member' table
ALTER TABLE public.member
    DROP COLUMN github_id;

-- Revert the 'member_view' to exclude the 'github_id' column
DROP VIEW public.member_view;
CREATE VIEW public.member_view AS
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
    member.logto_id
   FROM ((public.member
     LEFT JOIN public.member_role_relation ON ((member.member_id = (member_role_relation.member_id)::bpchar)))
     LEFT JOIN public.role ON ((member_role_relation.role_id = role.role_id)));