-- 004_refactor_auction_period.down.sql
-- start_at/end_at (TIMESTAMPTZ) から auction_date/start_time/end_time へ巻き戻す。

DO $$
BEGIN
    -- start_at カラムが存在する場合のみロールバックを実行する
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'auctions' AND column_name = 'start_at'
    ) THEN
        -- 旧カラムを追加
        ALTER TABLE auctions ADD COLUMN IF NOT EXISTS auction_date DATE;
        ALTER TABLE auctions ADD COLUMN IF NOT EXISTS start_time TIME;
        ALTER TABLE auctions ADD COLUMN IF NOT EXISTS end_time TIME;

        -- start_at/end_at を JST の日付・時刻へ分解して移行
        -- NOTE: start_at が NULL の行は auction_date に CURRENT_DATE を設定する（仮作成レコードの開発用フォールバック）
        UPDATE auctions
        SET
            auction_date = COALESCE(
                (start_at AT TIME ZONE 'Asia/Tokyo')::date,
                CURRENT_DATE
            ),
            start_time = CASE
                WHEN start_at IS NOT NULL
                    THEN (start_at AT TIME ZONE 'Asia/Tokyo')::time
                ELSE NULL
            END,
            end_time = CASE
                WHEN end_at IS NOT NULL
                    THEN (end_at AT TIME ZONE 'Asia/Tokyo')::time
                ELSE NULL
            END;

        -- auction_date に NOT NULL 制約を付与（NULL の行がある場合はエラー）
        ALTER TABLE auctions ALTER COLUMN auction_date SET NOT NULL;
        -- 旧ユニーク制約を復元
        ALTER TABLE auctions ADD CONSTRAINT auctions_venue_id_auction_date_key UNIQUE (venue_id, auction_date);
        -- 旧インデックスを復元
        CREATE INDEX IF NOT EXISTS idx_auctions_date ON auctions (auction_date);

        -- 新インデックスを削除
        DROP INDEX IF EXISTS idx_auctions_start_at_jst_date;
        DROP INDEX IF EXISTS uq_auctions_venue_date;
        DROP INDEX IF EXISTS idx_auctions_start_at;

        -- 新カラムを削除
        ALTER TABLE auctions DROP COLUMN start_at;
        ALTER TABLE auctions DROP COLUMN end_at;
    END IF;
END $$;

-- created_at/updated_at を TIMESTAMPTZ から TIMESTAMP へ戻す
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'auctions'
          AND column_name = 'created_at'
          AND data_type = 'timestamp with time zone'
    ) THEN
        ALTER TABLE auctions
            ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'UTC',
            ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'UTC';
    END IF;
END $$;
