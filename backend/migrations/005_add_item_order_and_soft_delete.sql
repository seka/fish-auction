-- Add sort_order and deleted_at to auction_items
ALTER TABLE auction_items ADD COLUMN IF NOT EXISTS sort_order INTEGER DEFAULT 0;
ALTER TABLE auction_items ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- Add index for soft delete and sorting
CREATE INDEX IF NOT EXISTS idx_auction_items_deleted_at ON auction_items(deleted_at);
CREATE INDEX IF NOT EXISTS idx_auction_items_auction_id_sort_order ON auction_items(auction_id, sort_order);
