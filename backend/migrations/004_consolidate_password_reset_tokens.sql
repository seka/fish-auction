-- Drop old tables
DROP TABLE IF EXISTS admin_password_reset_tokens;
DROP TABLE IF EXISTS buyer_password_reset_tokens;

-- Create unified table
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL,
    user_role VARCHAR(50) NOT NULL CHECK (user_role IN ('admin', 'buyer')),
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_password_reset_tokens_token_hash ON password_reset_tokens(token_hash);
CREATE INDEX idx_password_reset_tokens_user_id_role ON password_reset_tokens(user_id, user_role);
