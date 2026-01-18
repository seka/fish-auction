-- 2026-01-13 のテストデータ修正 (商品ステータスを 'Pending' に修正)
-- 豊洲市場 (ID: 13)

DO $$
DECLARE
    current_auction_id INT;
BEGIN
    SELECT id INTO current_auction_id FROM auctions WHERE venue_id = 13 AND auction_date = '2026-01-13';

    IF current_auction_id IS NOT NULL THEN
        -- 商品のステータスを入札受付中 ('Pending') に更新
        UPDATE auction_items 
        SET status = 'Pending' 
        WHERE auction_id = current_auction_id;
    END IF;
END $$;
