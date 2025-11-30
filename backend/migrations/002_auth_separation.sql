-- Migration: Separate authentication information
-- Purpose: Improve security, extensibility, and support for multiple auth methods

-- Create authentications table
CREATE TABLE IF NOT EXISTS authentications (
    id SERIAL PRIMARY KEY,
    buyer_id INT NOT NULL REFERENCES buyers(id) ON DELETE CASCADE,
    email VARCHAR(255) UNIQUE NOT NULL CHECK (TRIM(email) <> ''),
    password_hash TEXT NOT NULL,
    auth_type VARCHAR(50) NOT NULL DEFAULT 'password',
    failed_attempts INT NOT NULL DEFAULT 0,
    locked_until TIMESTAMP,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_authentications_email ON authentications(email);
CREATE INDEX IF NOT EXISTS idx_authentications_buyer_id ON authentications(buyer_id);

-- Remove authentication-related columns from buyers table
ALTER TABLE buyers DROP COLUMN IF EXISTS password_hash;
ALTER TABLE buyers DROP CONSTRAINT IF EXISTS buyers_name_key;

-- Note: Existing data migration
-- If there are existing buyers with password_hash, you'll need to:
-- 1. Create corresponding authentication records with a temporary email
-- 2. Or, clear the buyers table and start fresh
-- For development, we recommend starting fresh.
