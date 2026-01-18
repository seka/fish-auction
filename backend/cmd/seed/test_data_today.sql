-- 2026-01-19 のテストデータ作成
-- 豊洲市場 (ID: 13) にセリを作成

DO $$
DECLARE
    new_auction_id INT;
BEGIN
    -- セリの作成（本日分）
    INSERT INTO auctions (venue_id, auction_date, start_time, end_time, status, created_at, updated_at)
    VALUES (13, '2026-01-19', '00:00:00', '23:59:59', 'in_progress', NOW(), NOW())
    ON CONFLICT (venue_id, auction_date) DO UPDATE SET status = 'in_progress', updated_at = NOW()
    RETURNING id INTO new_auction_id;

    -- 商品の追加
    DELETE FROM auction_items WHERE auction_id = new_auction_id;

    INSERT INTO auction_items (auction_id, fisherman_id, fish_type, quantity, unit, start_price, highest_bid, status, created_at, updated_at)
    VALUES 
    (new_auction_id, 1076, '【テスト】本マグロ', 120, 'kg', 400000, 400000, 'Pending', NOW(), NOW()),
    (new_auction_id, 1076, '【テスト】高級ウニ', 2, '箱', 20000, 20000, 'Pending', NOW(), NOW()),
    (new_auction_id, 1076, '【テスト】寒ブリ', 10, 'kg', 5000, 5000, 'Pending', NOW(), NOW());
END $$;
