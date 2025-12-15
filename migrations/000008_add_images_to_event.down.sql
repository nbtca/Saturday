-- Recreate event_log_view without images column
DROP VIEW public.event_log_view;
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

-- Recreate event_view without images column
DROP VIEW public.event_view;
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

-- Remove images field from event_log table
ALTER TABLE event_log DROP COLUMN images;

-- Remove images field from event table
ALTER TABLE event DROP COLUMN images;
