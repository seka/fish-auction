-- Fishermen
CREATE TABLE IF NOT EXISTS fishermen (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (TRIM(name) <> ''),
    deleted_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_fishermen_deleted_at ON fishermen(deleted_at);

-- Buyers
CREATE TABLE IF NOT EXISTS buyers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (TRIM(name) <> ''),
    organization TEXT NOT NULL DEFAULT '',
    contact_info TEXT NOT NULL DEFAULT '',
    deleted_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_buyers_deleted_at ON buyers(deleted_at);

-- Authentications (Buyer Auth)
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

CREATE INDEX IF NOT EXISTS idx_authentications_email ON authentications(email);
CREATE INDEX IF NOT EXISTS idx_authentications_buyer_id ON authentications(buyer_id);

-- Admins
CREATE TABLE IF NOT EXISTS admins (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL CHECK (TRIM(email) <> ''),
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Venues
CREATE TABLE IF NOT EXISTS venues (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (TRIM(name) <> ''),
    location TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_venues_deleted_at ON venues(deleted_at);

-- Auctions
CREATE TABLE IF NOT EXISTS auctions (
    id SERIAL PRIMARY KEY,
    venue_id INTEGER NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    auction_date DATE NOT NULL,
    start_time TIME,
    end_time TIME,
    status VARCHAR(50) NOT NULL DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'in_progress', 'completed', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (venue_id, auction_date)
);

CREATE INDEX IF NOT EXISTS idx_auctions_venue_id ON auctions(venue_id);
CREATE INDEX IF NOT EXISTS idx_auctions_date ON auctions(auction_date);
CREATE INDEX IF NOT EXISTS idx_auctions_status ON auctions(status);

-- Auction Items
CREATE TABLE IF NOT EXISTS auction_items (
    id SERIAL PRIMARY KEY,
    fisherman_id INTEGER NOT NULL REFERENCES fishermen(id),
    auction_id INTEGER REFERENCES auctions(id) ON DELETE CASCADE,
    fish_type VARCHAR(255) NOT NULL CHECK (TRIM(fish_type) <> ''),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit VARCHAR(50) NOT NULL CHECK (TRIM(unit) <> ''),
    status VARCHAR(50) DEFAULT 'Pending',
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_auction_items_auction_id ON auction_items(auction_id);
CREATE INDEX IF NOT EXISTS idx_auction_items_deleted_at ON auction_items(deleted_at);
CREATE INDEX IF NOT EXISTS idx_auction_items_auction_id_sort_order ON auction_items(auction_id, sort_order);

-- Sort Order Trigger for Auction Items
CREATE OR REPLACE FUNCTION set_item_sort_order()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.sort_order IS NULL OR NEW.sort_order = 0 THEN
        NEW.sort_order := (
            SELECT COALESCE(MAX(sort_order), 0) + 1 
            FROM auction_items 
            WHERE auction_id = NEW.auction_id
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_item_sort_order
BEFORE INSERT ON auction_items
FOR EACH ROW
EXECUTE FUNCTION set_item_sort_order();

-- Transactions
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL REFERENCES auction_items(id),
    buyer_id INTEGER NOT NULL REFERENCES buyers(id),
    price INTEGER NOT NULL CHECK (price > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_transactions_buyer_id ON transactions(buyer_id);
CREATE INDEX IF NOT EXISTS idx_transactions_item_id ON transactions(item_id);

-- Password Reset Tokens
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL,
    user_role VARCHAR(50) NOT NULL CHECK (user_role IN ('admin', 'buyer')),
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_token_hash ON password_reset_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id_role ON password_reset_tokens(user_id, user_role);

-- Push Subscriptions
CREATE TABLE IF NOT EXISTS push_subscriptions (
    id SERIAL PRIMARY KEY,
    buyer_id INTEGER NOT NULL REFERENCES buyers(id) ON DELETE CASCADE,
    endpoint TEXT NOT NULL UNIQUE,
    p256dh TEXT NOT NULL,
    auth TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_push_subscriptions_buyer_id ON push_subscriptions(buyer_id);
