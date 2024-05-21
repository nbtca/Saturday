ALTER TABLE "client" ALTER COLUMN "gmt_create" SET DEFAULT 'now()';
ALTER TABLE "client" ALTER COLUMN "gmt_modified" SET DEFAULT 'now()';


ALTER TABLE "event" ALTER COLUMN "gmt_create" SET DEFAULT 'now()';
ALTER TABLE "event" ALTER COLUMN "gmt_modified" SET DEFAULT 'now()';


ALTER TABLE "event_log" ALTER COLUMN "gmt_create" SET DEFAULT 'now()';

ALTER TABLE "event_log" ALTER COLUMN "gmt_create" SET DEFAULT 'now()';

ALTER TABLE "member" ALTER COLUMN "gmt_create" SET DEFAULT 'now()';
ALTER TABLE "member" ALTER COLUMN "gmt_modified" SET DEFAULT 'now()';