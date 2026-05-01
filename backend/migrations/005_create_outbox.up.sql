CREATE TABLE outbox (
    id             BIGSERIAL    PRIMARY KEY,
    job_type       VARCHAR(50)  NOT NULL,
    schema_version INTEGER      NOT NULL DEFAULT 1,
    payload        JSONB        NOT NULL,
    status         VARCHAR(20)  NOT NULL DEFAULT 'pending',
    attempts       INTEGER      NOT NULL DEFAULT 0,
    max_attempts   INTEGER      NOT NULL DEFAULT 5,
    available_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    claimed_at     TIMESTAMPTZ  NULL,
    claimed_by     VARCHAR(255) NULL,
    last_error     TEXT         NULL,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    processed_at   TIMESTAMPTZ  NULL
);

CREATE INDEX idx_outbox_claimable
    ON outbox (available_at, id)
    WHERE status = 'pending';
