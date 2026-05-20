CREATE TABLE rate_limit_counters (
    key          TEXT        PRIMARY KEY,
    count        INTEGER     NOT NULL DEFAULT 1,
    window_start TIMESTAMPTZ NOT NULL
);
