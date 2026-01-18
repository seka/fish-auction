-- Add deleted_at column to master tables for soft delete

ALTER TABLE fishermen ADD COLUMN deleted_at TIMESTAMP;
ALTER TABLE buyers ADD COLUMN deleted_at TIMESTAMP;
ALTER TABLE venues ADD COLUMN deleted_at TIMESTAMP;

-- Create indexes for performance
CREATE INDEX idx_fishermen_deleted_at ON fishermen(deleted_at);
CREATE INDEX idx_buyers_deleted_at ON buyers(deleted_at);
CREATE INDEX idx_venues_deleted_at ON venues(deleted_at);
