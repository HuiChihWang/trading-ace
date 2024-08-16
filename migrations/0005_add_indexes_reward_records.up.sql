CREATE UNIQUE INDEX reward_records_task_id ON reward_records (task_id);
CREATE INDEX reward_records_user_id ON reward_records (user_id);
CREATE INDEX reward_records_created_at ON reward_records (created_at);