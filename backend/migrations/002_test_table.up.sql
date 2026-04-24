-- Test migration to verify dynamic loading
CREATE TABLE IF NOT EXISTS migration_test (
    id SERIAL PRIMARY KEY,
    test_val TEXT
);
