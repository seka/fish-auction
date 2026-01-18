-- 2026-01-13 のテストデータ作成 (完全修正版)
-- 豊洲市場 (ID: 13) にセリを作成

DO $$
DECLARE
    new_auction_id INT;
BEGIN
    INSERT INTO auctions (venue_id, auction_date, start_time, end_time, status, created_at, updated_at)
    VALUES (13, '2026-01-13', '00:00:00', '23:59:59', 'in_progress', NOW(), NOW())
    ON CONFLICT (venue_id, auction_date) DO UPDATE SET updated_at = NOW()
    RETURNING id INTO new_auction_id;

    -- 商品を追加
    INSERT INTO auction_items (auction_id, fisherman_id, fish_type, quantity, unit, start_price, highest_bid, status, created_at, updated_at)
    VALUES 
    (new_auction_id, 1076, '本マグロ（豊洲）', 150, 'kg', 500000, 500000, 'in_progress', NOW(), NOW()),
    (new_auction_id, 1076, '真鯛（豊洲）', 2, 'kg', 3000, 3000, 'in_progress', NOW(), NOW()),
    (new_auction_id, 1076, 'ウニ（豊洲）', 1, '箱', 15000, 15000, 'in_progress', NOW(), NOW());
END $$;
