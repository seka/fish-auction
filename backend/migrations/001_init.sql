CREATE TABLE IF NOT EXISTS fishermen (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (TRIM(name) <> '')
);

CREATE TABLE IF NOT EXISTS buyers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (TRIM(name) <> '')
);

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
