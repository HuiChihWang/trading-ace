CREATE TABLE reward_records
(
    id         SERIAL PRIMARY KEY,
    user_id    VARCHAR(255)     NOT NULL,
    points     DOUBLE PRECISION NOT NULL,
    task_id    INTEGER          NOT NULL,
    created_at TIMESTAMP        NOT NULL
);