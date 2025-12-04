-- Remove images field from event_log table
ALTER TABLE event_log DROP COLUMN images;

-- Remove images field from event table
ALTER TABLE event DROP COLUMN images;
