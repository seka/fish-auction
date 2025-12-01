-- Fishermen
CREATE TABLE IF NOT EXISTS fishermen (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (TRIM(name) <> '')
);

-- Buyers
CREATE TABLE IF NOT EXISTS buyers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (TRIM(name) <> ''),
    organization TEXT NOT NULL DEFAULT '',
    contact_info TEXT NOT NULL DEFAULT ''
);

-- Authentications
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

-- Venues
CREATE TABLE IF NOT EXISTS venues (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (TRIM(name) <> ''),
    location TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Auctions
CREATE TABLE IF NOT EXISTS auctions (
    id SERIAL PRIMARY KEY,
    venue_id INTEGER NOT NULL REFERENCES venues(id) ON DELETE RESTRICT,
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
    auction_id INTEGER REFERENCES auctions(id) ON DELETE RESTRICT,
    fish_type VARCHAR(255) NOT NULL CHECK (TRIM(fish_type) <> ''),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit VARCHAR(50) NOT NULL CHECK (TRIM(unit) <> ''),
    status VARCHAR(50) DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_auction_items_auction_id ON auction_items(auction_id);

-- Transactions
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL REFERENCES auction_items(id),
    buyer_id INTEGER NOT NULL REFERENCES buyers(id),
    price INTEGER NOT NULL CHECK (price > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Default Data
INSERT INTO venues (name, location, description) 
VALUES ('豊洲市場', '東京都江東区豊洲6-6-1', 'デフォルト会場')
ON CONFLICT DO NOTHING;

-- Default Auction
DO $$
DECLARE
    default_venue_id INTEGER;
BEGIN
    SELECT id INTO default_venue_id FROM venues WHERE name = '豊洲市場' LIMIT 1;
    
    IF default_venue_id IS NOT NULL THEN
        INSERT INTO auctions (venue_id, auction_date, status)
        VALUES (default_venue_id, CURRENT_DATE, 'in_progress')
        ON CONFLICT DO NOTHING;
    END IF;
END $$;
