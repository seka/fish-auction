-- 004_refactor_auction_period.up.sql
-- auctions テーブルを auction_date/start_time/end_time から start_at/end_at (TIMESTAMPTZ) へ移行する。
-- フレッシュ DB では 001_init.up.sql に start_at/end_at が含まれるが、
-- 既存 DB では auctions テーブルが既に存在すると 001 が列追加を行わないため、
-- 本マイグレーション内でも必要に応じて start_at/end_at を追加する。
-- フレッシュ DB（001 から構築済み）では旧カラムが存在しないため全 DO ブロックが no-op となる。

DO $$
BEGIN
    -- 旧カラムが存在する場合のみ、データを移行して削除する
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'auctions' AND column_name = 'auction_date'
    ) THEN
        -- 既存DBアップグレード時に start_at/end_at が未作成でも失敗しないよう補完
        ALTER TABLE auctions ADD COLUMN IF NOT EXISTS start_at TIMESTAMPTZ;
        ALTER TABLE auctions ADD COLUMN IF NOT EXISTS end_at TIMESTAMPTZ;

        -- auction_date + start_time/end_time を JST タイムスタンプとして start_at/end_at へ移行
        UPDATE auctions
        SET
            start_at = COALESCE(
                start_at,
                CASE
                    WHEN auction_date IS NULL THEN NULL
                    WHEN start_time IS NOT NULL
                        THEN (auction_date::timestamp + start_time) AT TIME ZONE 'Asia/Tokyo'
                    ELSE auction_date::timestamp AT TIME ZONE 'Asia/Tokyo'
                END
            ),
            end_at = COALESCE(
                end_at,
                CASE
                    WHEN auction_date IS NOT NULL AND end_time IS NOT NULL
                        THEN (auction_date::timestamp + end_time) AT TIME ZONE 'Asia/Tokyo'
                    ELSE NULL
                END
            );

        -- 旧ユニーク制約を削除
        ALTER TABLE auctions DROP CONSTRAINT IF EXISTS auctions_venue_id_auction_date_key;
        -- 旧インデックスを削除
        DROP INDEX IF EXISTS idx_auctions_date;
        -- 旧カラムを削除
        ALTER TABLE auctions DROP COLUMN auction_date;
        ALTER TABLE auctions DROP COLUMN start_time;
        ALTER TABLE auctions DROP COLUMN end_time;
    END IF;
END $$;

-- created_at/updated_at が TIMESTAMP (タイムゾーンなし) の場合は TIMESTAMPTZ へ変換する
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'auctions'
          AND column_name = 'created_at'
          AND data_type = 'timestamp without time zone'
    ) THEN
        ALTER TABLE auctions
            ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC',
            ALTER COLUMN updated_at TYPE TIMESTAMPTZ USING updated_at AT TIME ZONE 'UTC';
    END IF;
END $$;
