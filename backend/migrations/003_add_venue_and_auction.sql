-- Migration: Add Venue and Auction models
-- Purpose: Introduce venue (会場) and auction (セリ) concepts to manage auctions by date and venue

-- Venue (会場マスタ)
CREATE TABLE IF NOT EXISTS venues (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (TRIM(name) <> ''),
    location TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Auction (セリイベント)
CREATE TABLE IF NOT EXISTS auctions (
    id SERIAL PRIMARY KEY,
    venue_id INTEGER NOT NULL REFERENCES venues(id) ON DELETE RESTRICT,
    auction_date DATE NOT NULL,
    start_time TIME,
    end_time TIME,
    status VARCHAR(50) NOT NULL DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'in_progress', 'completed', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add auction_id to auction_items
ALTER TABLE auction_items ADD COLUMN IF NOT EXISTS auction_id INTEGER REFERENCES auctions(id) ON DELETE RESTRICT;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_auctions_venue_id ON auctions(venue_id);
CREATE INDEX IF NOT EXISTS idx_auctions_date ON auctions(auction_date);
CREATE INDEX IF NOT EXISTS idx_auctions_status ON auctions(status);
CREATE INDEX IF NOT EXISTS idx_auction_items_auction_id ON auction_items(auction_id);

-- Insert default venue for development
INSERT INTO venues (name, location, description) 
VALUES ('豊洲市場', '東京都江東区豊洲6-6-1', 'デフォルト会場')
ON CONFLICT DO NOTHING;

-- Insert default auction for existing items (development only)
-- This allows existing auction_items to be associated with an auction
DO $$
DECLARE
    default_venue_id INTEGER;
    default_auction_id INTEGER;
BEGIN
    -- Get the default venue ID
    SELECT id INTO default_venue_id FROM venues WHERE name = '豊洲市場' LIMIT 1;
    
    IF default_venue_id IS NOT NULL THEN
        -- Create a default auction for today
        INSERT INTO auctions (venue_id, auction_date, status)
        VALUES (default_venue_id, CURRENT_DATE, 'in_progress')
        ON CONFLICT DO NOTHING
        RETURNING id INTO default_auction_id;
        
        -- If auction was created or already exists, associate orphan items
        IF default_auction_id IS NULL THEN
            SELECT id INTO default_auction_id FROM auctions 
            WHERE venue_id = default_venue_id AND auction_date = CURRENT_DATE 
            LIMIT 1;
        END IF;
        
        IF default_auction_id IS NOT NULL THEN
            -- Update existing auction_items that don't have an auction_id
            UPDATE auction_items 
            SET auction_id = default_auction_id 
            WHERE auction_id IS NULL;
        END IF;
    END IF;
END $$;
