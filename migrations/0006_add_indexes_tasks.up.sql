CREATE INDEX tasks_type_status_idx ON tasks (type, status);
CREATE INDEX tasks_user_id ON tasks (user_id);
CREATE INDEX tasks_created_at ON tasks (created_at);