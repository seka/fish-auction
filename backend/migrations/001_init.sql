CREATE TABLE IF NOT EXISTS fishermen (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (TRIM(name) <> '')
);

CREATE TABLE IF NOT EXISTS buyers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (TRIM(name) <> ''),
    password_hash TEXT NOT NULL DEFAULT '',
    organization TEXT NOT NULL DEFAULT '',
    contact_info TEXT NOT NULL DEFAULT '',
    UNIQUE(name)
);

-- For existing tables where columns might be missing (idempotent updates)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'buyers' AND column_name = 'password_hash') THEN
        ALTER TABLE buyers ADD COLUMN password_hash TEXT NOT NULL DEFAULT '';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'buyers' AND column_name = 'organization') THEN
        ALTER TABLE buyers ADD COLUMN organization TEXT NOT NULL DEFAULT '';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'buyers' AND column_name = 'contact_info') THEN
        ALTER TABLE buyers ADD COLUMN contact_info TEXT NOT NULL DEFAULT '';
    END IF;
    -- Add unique constraint if not exists (handling this in SQL block is tricky, assuming CREATE TABLE handles it for new, 
    -- for existing we might need to check constraint. For now, relying on CREATE TABLE for new DBs. 
    -- Existing DBs might fail on duplicate names if we enforce unique. 
    -- Let's try to add unique index concurrently or just leave it for CREATE TABLE for now to avoid complex migration logic in init file.)
END $$;

CREATE TABLE IF NOT EXISTS auction_items (
    id SERIAL PRIMARY KEY,
    fisherman_id INTEGER NOT NULL REFERENCES fishermen(id),
    fish_type VARCHAR(255) NOT NULL CHECK (TRIM(fish_type) <> ''),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit VARCHAR(50) NOT NULL CHECK (TRIM(unit) <> ''),
    status VARCHAR(50) DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL REFERENCES auction_items(id),
    buyer_id INTEGER NOT NULL REFERENCES buyers(id),
    price INTEGER NOT NULL CHECK (price > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
