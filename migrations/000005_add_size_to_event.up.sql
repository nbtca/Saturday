ALTER TABLE public.event
ADD COLUMN size character varying(50) DEFAULT ''::character varying NOT NULL;

COMMENT ON COLUMN public.event.size IS '工作量';

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