-- Rename existing table to buyer_password_reset_tokens
ALTER TABLE password_reset_tokens RENAME TO buyer_password_reset_tokens;

-- Rename indices for consistency
ALTER INDEX idx_password_reset_tokens_token_hash RENAME TO idx_buyer_password_reset_tokens_token_hash;
ALTER INDEX idx_password_reset_tokens_buyer_id RENAME TO idx_buyer_password_reset_tokens_buyer_id;

-- Create admin_password_reset_tokens table
CREATE TABLE admin_password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id INTEGER NOT NULL REFERENCES admins(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_admin_password_reset_tokens_token_hash ON admin_password_reset_tokens(token_hash);
CREATE INDEX idx_admin_password_reset_tokens_admin_id ON admin_password_reset_tokens(admin_id);
