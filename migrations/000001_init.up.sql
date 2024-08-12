-- Active: 1716098742535@@127.0.0.1@5432@weekend
--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3 (Debian 16.3-1.pgdg120+1)
-- Dumped by pg_dump version 16.3 (Debian 16.3-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: client; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS public.client (
    client_id bigint NOT NULL,
    openid character(28) DEFAULT ''::bpchar,
    gmt_create timestamp without time zone NOT NULL,
    gmt_modified timestamp without time zone NOT NULL
);


ALTER TABLE public.client OWNER TO postgres;

--
-- Name: COLUMN client.openid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client.openid IS '微信';


--
-- Name: client_client_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.client ALTER COLUMN client_id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.client_client_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: event; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS public.event (
    client_id bigint NOT NULL,
    event_id bigint NOT NULL,
    model character varying(40) DEFAULT ''::character varying,
    phone character varying(11) DEFAULT ''::character varying NOT NULL,
    qq character varying(20) DEFAULT ''::character varying,
    contact_preference character varying(20) DEFAULT 'qq'::character varying NOT NULL,
    problem character varying(500) DEFAULT ''::character varying,
    member_id character(10) DEFAULT ''::bpchar,
    closed_by character(10) DEFAULT ''::bpchar,
    gmt_create timestamp without time zone NOT NULL,
    gmt_modified timestamp without time zone NOT NULL
);


ALTER TABLE public.event OWNER TO postgres;


COMMENT ON COLUMN public.event.model IS '型号';
COMMENT ON COLUMN public.event.contact_preference IS '联系偏好';
COMMENT ON COLUMN public.event.problem IS '事件（用户）描述';
COMMENT ON COLUMN public.event.member_id IS '最后由谁维修';
COMMENT ON COLUMN public.event.closed_by IS '由谁关闭';


--
-- Name: event_action; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS public.event_action (
    event_action_id smallint NOT NULL,
    action character varying(30) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.event_action OWNER TO postgres;

--
-- Name: event_event_action_relation; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS public.event_event_action_relation (
    event_log_id bigint NOT NULL,
    event_action_id smallint NOT NULL
);


ALTER TABLE public.event_event_action_relation OWNER TO postgres;

--
-- Name: event_event_action_relation_event_log_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.event_event_action_relation ALTER COLUMN event_log_id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.event_event_action_relation_event_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: event_event_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.event ALTER COLUMN event_id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.event_event_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: event_event_status_relation; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS public.event_event_status_relation (
    event_id bigint NOT NULL,
    event_status_id smallint NOT NULL
);


ALTER TABLE public.event_event_status_relation OWNER TO postgres;

--
-- Name: event_log; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS public.event_log (
    event_log_id bigint NOT NULL,
    event_id bigint NOT NULL,
    description character varying(255) DEFAULT ''::character varying,
    member_id character(10) DEFAULT ''::bpchar,
    gmt_create timestamp without time zone NOT NULL
);


ALTER TABLE public.event_log OWNER TO postgres;

--
-- Name: event_log_event_log_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.event_log ALTER COLUMN event_log_id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.event_log_event_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: event_log_view; Type: VIEW; Schema: public; Owner: postgres
--

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


ALTER VIEW public.event_log_view OWNER TO postgres;

--
-- Name: event_status; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS public.event_status (
    event_status_id smallint NOT NULL,
    status character varying(255) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.event_status OWNER TO postgres;

--
-- Name: event_view; Type: VIEW; Schema: public; Owner: postgres
--

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
    event_status.status
   FROM ((public.event
     LEFT JOIN public.event_event_status_relation ON ((event.event_id = event_event_status_relation.event_id)))
     LEFT JOIN public.event_status ON ((event_event_status_relation.event_status_id = event_status.event_status_id)));


ALTER VIEW public.event_view OWNER TO postgres;

--
-- Name: member; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS public.member (
    member_id character(10) NOT NULL,
    alias character varying(50) DEFAULT ''::character varying,
    password character varying(50) DEFAULT ''::character varying,
    name character varying(20) DEFAULT ''::character varying,
    section character varying(20) DEFAULT ''::character varying,
    profile character varying(1000) DEFAULT ''::character varying,
    phone character varying(11) DEFAULT ''::character varying,
    qq character varying(20) DEFAULT ''::character varying,
    avatar character varying(255) DEFAULT ''::character varying,
    created_by character(10) DEFAULT ''::bpchar,
    gmt_create timestamp without time zone NOT NULL,
    gmt_modified timestamp without time zone NOT NULL,
    logto_id character varying(50) DEFAULT ''::character varying
);


ALTER TABLE public.member OWNER TO postgres;

--
-- Name: COLUMN member.alias; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.member.alias IS '昵称';


--
-- Name: COLUMN member.section; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.member.section IS '班级（计算机196）';


--
-- Name: COLUMN member.profile; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.member.profile IS '个人简介';


--
-- Name: COLUMN member.avatar; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.member.avatar IS '头像地址';


--
-- Name: COLUMN member.created_by; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.member.created_by IS '由谁添加';


--
-- Name: member_role_relation; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS public.member_role_relation (
    member_id character varying(10) NOT NULL,
    role_id smallint NOT NULL
);


ALTER TABLE public.member_role_relation OWNER TO postgres;

--
-- Name: role; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS public.role (
    role_id smallint NOT NULL,
    role character varying(255) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.role OWNER TO postgres;

--
-- Name: member_view; Type: VIEW; Schema: public; Owner: postgres
--

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
    role.role,
    member.logto_id
   FROM ((public.member
     LEFT JOIN public.member_role_relation ON ((member.member_id = (member_role_relation.member_id)::bpchar)))
     LEFT JOIN public.role ON ((member_role_relation.role_id = role.role_id)));


ALTER VIEW public.member_view OWNER TO postgres;

--
-- Name: setting; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE IF NOT EXISTS public.setting (
    setting character varying(10000) DEFAULT ''::character varying
);


ALTER TABLE public.setting OWNER TO postgres;

--
-- Data for Name: client; Type: TABLE DATA; Schema: public; Owner: postgres
--


--
-- Data for Name: event_action; Type: TABLE DATA; Schema: public; Owner: postgres
--
INSERT INTO public.event_action (event_action_id, action) VALUES
(1, 'create'),
(2, 'accept'),
(3, 'cancel'),
(4, 'commit'),
(5, 'alterCommit'),
(6, 'drop'),
(7, 'close'),
(8, 'reject'),
(9, 'update');

--
-- Data for Name: event_status; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.event_status (event_status_id, status) VALUES
(1, 'open'),
(2, 'accepted'),
(3, 'cancelled'),
(4, 'committed'),
(5, 'closed');


--
-- Data for Name: role; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.role (role_id,role) VALUES
(0,'member_inactive'),
(1,'admin_inactive'),
(2,'member'),
(4,'admin');

--
-- Name: client client_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.client
    ADD CONSTRAINT client_pkey PRIMARY KEY (client_id);


--
-- Name: event_action event_action_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_action
    ADD CONSTRAINT event_action_pkey PRIMARY KEY (event_action_id);


--
-- Name: event_event_action_relation event_event_action_relation_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_event_action_relation
    ADD CONSTRAINT event_event_action_relation_pkey PRIMARY KEY (event_log_id);


--
-- Name: event_event_status_relation event_event_status_relation_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_event_status_relation
    ADD CONSTRAINT event_event_status_relation_pkey PRIMARY KEY (event_id);


--
-- Name: event_log event_log_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_log
    ADD CONSTRAINT event_log_pkey PRIMARY KEY (event_log_id);


--
-- Name: event event_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_pkey PRIMARY KEY (event_id);


--
-- Name: event_status event_status_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_status
    ADD CONSTRAINT event_status_pkey PRIMARY KEY (event_status_id);


--
-- Name: member member_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.member
    ADD CONSTRAINT member_pkey PRIMARY KEY (member_id);


--
-- Name: member_role_relation member_role_relation_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.member_role_relation
    ADD CONSTRAINT member_role_relation_pkey PRIMARY KEY (member_id);


--
-- Name: role role_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role
    ADD CONSTRAINT role_pkey PRIMARY KEY (role_id);


--
-- Name: event_client_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX event_client_id_idx ON public.event USING btree (client_id);


--
-- Name: event_closed_by_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX event_closed_by_idx ON public.event USING btree (closed_by);


--
-- Name: event_event_action_relation_event_action_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX event_event_action_relation_event_action_id_idx ON public.event_event_action_relation USING btree (event_action_id);


--
-- Name: event_event_status_relation_event_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX event_event_status_relation_event_id_idx ON public.event_event_status_relation USING btree (event_id);


--
-- Name: event_event_status_relation_event_status_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX event_event_status_relation_event_status_id_idx ON public.event_event_status_relation USING btree (event_status_id);


--
-- Name: event_log_event_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX event_log_event_id_idx ON public.event_log USING btree (event_id);


--
-- Name: event_log_member_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX event_log_member_id_idx ON public.event_log USING btree (member_id);


--
-- Name: event_member_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX event_member_id_idx ON public.event USING btree (member_id);


--
-- Name: member_role_relation_role_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX member_role_relation_role_id_idx ON public.member_role_relation USING btree (role_id);


--
-- Name: event event_client_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_client_id_fkey FOREIGN KEY (client_id) REFERENCES public.client(client_id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- Name: event_event_action_relation event_event_action_relation_event_action_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_event_action_relation
    ADD CONSTRAINT event_event_action_relation_event_action_id_fkey FOREIGN KEY (event_action_id) REFERENCES public.event_action(event_action_id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- Name: event_event_action_relation event_event_action_relation_event_log_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_event_action_relation
    ADD CONSTRAINT event_event_action_relation_event_log_id_fkey FOREIGN KEY (event_log_id) REFERENCES public.event_log(event_log_id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- Name: event_event_status_relation event_event_status_relation_event_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_event_status_relation
    ADD CONSTRAINT event_event_status_relation_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.event(event_id);


--
-- Name: event_event_status_relation event_event_status_relation_event_status_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_event_status_relation
    ADD CONSTRAINT event_event_status_relation_event_status_id_fkey FOREIGN KEY (event_status_id) REFERENCES public.event_status(event_status_id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- Name: event_log event_log_event_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_log
    ADD CONSTRAINT event_log_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.event(event_id);


--
-- Name: member_role_relation member_role_relation_member_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.member_role_relation
    ADD CONSTRAINT member_role_relation_member_id_fkey FOREIGN KEY (member_id) REFERENCES public.member(member_id);


--
-- Name: member_role_relation member_role_relation_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.member_role_relation
    ADD CONSTRAINT member_role_relation_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.role(role_id);


--
-- PostgreSQL database dump complete
--
