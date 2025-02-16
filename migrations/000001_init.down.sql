DROP VIEW IF EXISTS public.event_log_view;
DROP VIEW IF EXISTS public.event_view;
DROP VIEW IF EXISTS public.member_view;

DROP TABLE IF EXISTS public.setting;
DROP TABLE IF EXISTS public.member_role_relation;
DROP TABLE IF EXISTS public.role;
DROP TABLE IF EXISTS public.member;
DROP TABLE IF EXISTS public.event_event_status_relation;
DROP TABLE IF EXISTS public.event_event_action_relation;
DROP TABLE IF EXISTS public.event_log;
DROP TABLE IF EXISTS public.event_action;
DROP TABLE IF EXISTS public.event_status;
DROP TABLE IF EXISTS public.event;
DROP TABLE IF EXISTS public.client;

DROP SEQUENCE IF EXISTS public.client_client_id_seq;
DROP SEQUENCE IF EXISTS public.event_event_id_seq;
DROP SEQUENCE IF EXISTS public.event_event_action_relation_event_log_id_seq;
DROP SEQUENCE IF EXISTS public.event_log_event_log_id_seq;