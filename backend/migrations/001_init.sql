CREATE TABLE IF NOT EXISTS fishermen (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS buyers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS auction_items (
    id SERIAL PRIMARY KEY,
    fisherman_id INTEGER REFERENCES fishermen(id),
    fish_type VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL,
    unit VARCHAR(50) NOT NULL, -- e.g., 'kg', 'box', 'piece'
    status VARCHAR(50) DEFAULT 'Pending', -- 'Pending', 'Sold'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    item_id INTEGER REFERENCES auction_items(id),
    buyer_id INTEGER REFERENCES buyers(id),
    price INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
