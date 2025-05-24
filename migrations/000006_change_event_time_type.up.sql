DROP VIEW public.event_view;

ALTER TABLE event
ALTER COLUMN gmt_create
TYPE TIMESTAMP WITH TIME ZONE
USING gmt_create AT TIME ZONE 'Asia/Shanghai';

ALTER TABLE event
ALTER COLUMN gmt_modified
TYPE TIMESTAMP WITH TIME ZONE
USING gmt_create AT TIME ZONE 'Asia/Shanghai';

CREATE VIEW public.event_view AS
 SELECT event.event_id,
    event.client_id,
    event.model,
    event.phone,
    event.qq,
    event.contact_preference,
    event.problem,
    event.member_id,
    event.closed_by,
    event.gmt_create,
    event.gmt_modified,
    event.size,
    COALESCE(event_status.status, ''::character varying) AS status,
    event.github_issue_id,
    event.github_issue_number
   FROM ((public.event
     LEFT JOIN public.event_event_status_relation ON ((event.event_id = event_event_status_relation.event_id)))
     LEFT JOIN public.event_status ON ((event_event_status_relation.event_status_id = event_status.event_status_id)));

DROP VIEW public.event_log_view;

ALTER TABLE event_log
ALTER COLUMN gmt_create
TYPE TIMESTAMP WITH TIME ZONE
USING gmt_create AT TIME ZONE 'Asia/Shanghai';

CREATE VIEW public.event_log_view AS
 SELECT event_log.event_log_id,
    event_log.event_id,
    event_log.description,
    event_log.member_id,
    event_log.gmt_create,
    event_action.action
   FROM ((public.event_log
     LEFT JOIN public.event_event_action_relation ON ((event_log.event_log_id = event_event_action_relation.event_log_id)))
     LEFT JOIN public.event_action ON ((event_event_action_relation.event_action_id = event_action.event_action_id)));

DROP VIEW public.member_view;

ALTER TABLE member
ALTER COLUMN gmt_create
TYPE TIMESTAMP WITH TIME ZONE
USING gmt_create AT TIME ZONE 'Asia/Shanghai';

ALTER TABLE member
ALTER COLUMN gmt_modified
TYPE TIMESTAMP WITH TIME ZONE
USING gmt_modified AT TIME ZONE 'Asia/Shanghai';

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