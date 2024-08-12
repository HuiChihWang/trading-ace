CREATE TABLE tasks
(
    id           SERIAL PRIMARY KEY,
    status       VARCHAR(50)      NOT NULL,
    type         VARCHAR(50)      NOT NULL,
    user_id      VARCHAR(255)     NOT NULL,
    created_at   TIMESTAMP        NOT NULL,
    completed_at TIMESTAMP,
    swap_amount  DOUBLE PRECISION NOT NULL
);