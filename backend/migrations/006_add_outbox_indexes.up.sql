-- RecoverStale が WHERE status='processing' AND claimed_at < ... をスキャンするため。
CREATE INDEX idx_outbox_processing
    ON outbox (claimed_at)
    WHERE status = 'processing';

-- DeleteProcessedBefore が WHERE status='processed' AND processed_at < ... をスキャンするため。
CREATE INDEX idx_outbox_processed
    ON outbox (processed_at)
    WHERE status = 'processed';
