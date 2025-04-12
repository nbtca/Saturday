-- Add a new column 'github_id' to the 'member' table
ALTER TABLE public.member
    ADD COLUMN github_id character varying(50) DEFAULT ''::character varying;

-- Optionally, add a comment for the new column
COMMENT ON COLUMN public.member.github_id IS 'GitHub User ID';

-- Update the 'member_view' to include the new 'github_id' column
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
    member.logto_id,
    member.github_id -- Include the new column
   FROM ((public.member
     LEFT JOIN public.member_role_relation ON ((member.member_id = (member_role_relation.member_id)::bpchar)))
     LEFT JOIN public.role ON ((member_role_relation.role_id = role.role_id)));