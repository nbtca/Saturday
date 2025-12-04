-- Add images field to event table (client's problem images)
ALTER TABLE event ADD COLUMN images TEXT COMMENT '事件图片（JSON数组）';

-- Add images field to event_log table (member's repair images)
ALTER TABLE event_log ADD COLUMN images TEXT COMMENT '维修记录图片（JSON数组）';
