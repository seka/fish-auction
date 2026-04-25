-- Venues
INSERT INTO venues (name, location, description)
SELECT '豊洲市場', '東京都江東区豊洲6-6-1', 'デフォルト会場'
WHERE NOT EXISTS (SELECT 1 FROM venues WHERE name = '豊洲市場');

INSERT INTO venues (name, location, description)
SELECT '函館港', '北海道函館市', 'イカが有名'
WHERE NOT EXISTS (SELECT 1 FROM venues WHERE name = '函館港');

INSERT INTO venues (name, location, description)
SELECT '銚子港', '千葉県銚子市', '水揚げ量日本一'
WHERE NOT EXISTS (SELECT 1 FROM venues WHERE name = '銚子港');

INSERT INTO venues (name, location, description)
SELECT '焼津港', '静岡県焼津市', 'マグロ・カツオ'
WHERE NOT EXISTS (SELECT 1 FROM venues WHERE name = '焼津港');

-- Fishermen
INSERT INTO fishermen (name)
SELECT '田中 一郎'
WHERE NOT EXISTS (SELECT 1 FROM fishermen WHERE name = '田中 一郎');

INSERT INTO fishermen (name)
SELECT '鈴木 次郎'
WHERE NOT EXISTS (SELECT 1 FROM fishermen WHERE name = '鈴木 次郎');

INSERT INTO fishermen (name)
SELECT '佐藤 三郎'
WHERE NOT EXISTS (SELECT 1 FROM fishermen WHERE name = '佐藤 三郎');

-- Buyers
-- Default password for all seeded buyers is 'Password123'
-- Buyer 1
DO $$
DECLARE
    lid INTEGER;
BEGIN
    IF NOT EXISTS (SELECT 1 FROM buyers WHERE name = '株式会社 魚河岸') THEN
        INSERT INTO buyers (name, organization, contact_info)
        VALUES ('株式会社 魚河岸', '卸売', '03-1234-5678')
        RETURNING id INTO lid;

        INSERT INTO authentications (buyer_id, email, password_hash)
        VALUES (lid, 'uogashi@example.com', '$2a$10$bBUXj3ggxD38hbMPPcXRYe/CiPPfiP8pBv0x593dA5YrlEEjL.AxG');
    END IF;
END $$;

-- Buyer 2
DO $$
DECLARE
    lid INTEGER;
BEGIN
    IF NOT EXISTS (SELECT 1 FROM buyers WHERE name = 'すしざんまい') THEN
        INSERT INTO buyers (name, organization, contact_info)
        VALUES ('すしざんまい', '飲食店', '03-8765-4321')
        RETURNING id INTO lid;

        INSERT INTO authentications (buyer_id, email, password_hash)
        VALUES (lid, 'sushi@example.com', '$2a$10$bBUXj3ggxD38hbMPPcXRYe/CiPPfiP8pBv0x593dA5YrlEEjL.AxG');
    END IF;
END $$;

-- Buyer 3
DO $$
DECLARE
    lid INTEGER;
BEGIN
    IF NOT EXISTS (SELECT 1 FROM buyers WHERE name = 'スーパー玉出') THEN
        INSERT INTO buyers (name, organization, contact_info)
        VALUES ('スーパー玉出', '小売', '06-1234-5678')
        RETURNING id INTO lid;

        INSERT INTO authentications (buyer_id, email, password_hash)
        VALUES (lid, 'tamade@example.com', '$2a$10$bBUXj3ggxD38hbMPPcXRYe/CiPPfiP8pBv0x593dA5YrlEEjL.AxG');
    END IF;
END $$;

-- Auction (Create one for today at '函館港' if not exists)
DO $$
DECLARE
    vid INTEGER;
    today_jst timestamptz;
BEGIN
    SELECT id INTO vid FROM venues WHERE name = '函館港';
    today_jst := date_trunc('day', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Tokyo') AT TIME ZONE 'Asia/Tokyo';
    IF vid IS NOT NULL THEN
        INSERT INTO auctions (venue_id, start_at, end_at, status)
        SELECT vid, today_jst, today_jst + interval '1 day' - interval '1 second', 'in_progress'
        WHERE NOT EXISTS (
            SELECT 1 FROM auctions 
            WHERE venue_id = vid 
              AND start_at >= today_jst 
              AND start_at < today_jst + interval '1 day'
        );
    END IF;
END $$;

-- Auction (Create one for today at '豊洲市場' if not exists)
DO $$
DECLARE
    default_venue_id INTEGER;
    today_jst timestamptz;
BEGIN
    SELECT id INTO default_venue_id FROM venues WHERE name = '豊洲市場' LIMIT 1;
    today_jst := date_trunc('day', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Tokyo') AT TIME ZONE 'Asia/Tokyo';

    IF default_venue_id IS NOT NULL THEN
        INSERT INTO auctions (venue_id, start_at, end_at, status)
        SELECT default_venue_id, today_jst, today_jst + interval '1 day' - interval '1 second', 'in_progress'
        WHERE NOT EXISTS (
            SELECT 1 FROM auctions 
            WHERE venue_id = default_venue_id 
              AND start_at >= today_jst 
              AND start_at < today_jst + interval '1 day'
        );
    END IF;
END $$;

-- Auction Items
DO $$
DECLARE
    fid INTEGER;
    aid INTEGER;
    today_jst timestamptz;
BEGIN
    SELECT id INTO fid FROM fishermen WHERE name = '田中 一郎';
    today_jst := date_trunc('day', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Tokyo') AT TIME ZONE 'Asia/Tokyo';
    -- Get today's auction at Hakodate
    SELECT a.id INTO aid
    FROM auctions a
    JOIN venues v ON a.venue_id = v.id
    WHERE v.name = '函館港' 
      AND a.start_at >= today_jst 
      AND a.start_at < today_jst + interval '1 day'
    LIMIT 1;

    IF fid IS NOT NULL AND aid IS NOT NULL THEN
        INSERT INTO auction_items (fisherman_id, auction_id, fish_type, quantity, unit, sort_order)
        SELECT fid, aid, 'スルメイカ', 100, 'kg', 1
        WHERE NOT EXISTS (
            SELECT 1 FROM auction_items
            WHERE fisherman_id = fid AND auction_id = aid AND fish_type = 'スルメイカ'
        );

        INSERT INTO auction_items (fisherman_id, auction_id, fish_type, quantity, unit, sort_order)
        SELECT fid, aid, 'ホッケ', 50, 'kg', 2
        WHERE NOT EXISTS (
            SELECT 1 FROM auction_items
            WHERE fisherman_id = fid AND auction_id = aid AND fish_type = 'ホッケ'
        );
    END IF;
END $$;

DO $$
DECLARE
    fid INTEGER;
    aid INTEGER;
    today_jst timestamptz;
BEGIN
    SELECT id INTO fid FROM fishermen WHERE name = '鈴木 次郎';
    today_jst := date_trunc('day', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Tokyo') AT TIME ZONE 'Asia/Tokyo';
    -- Get today's auction at Hakodate
    SELECT a.id INTO aid
    FROM auctions a
    JOIN venues v ON a.venue_id = v.id
    WHERE v.name = '函館港' 
      AND a.start_at >= today_jst 
      AND a.start_at < today_jst + interval '1 day'
    LIMIT 1;

    IF fid IS NOT NULL AND aid IS NOT NULL THEN
        INSERT INTO auction_items (fisherman_id, auction_id, fish_type, quantity, unit, sort_order)
        SELECT fid, aid, 'マグロ', 200, 'kg', 3
        WHERE NOT EXISTS (
            SELECT 1 FROM auction_items
            WHERE fisherman_id = fid AND auction_id = aid AND fish_type = 'マグロ'
        );
    END IF;
END $$;

-- Auction Items (Toyosu)
DO $$
DECLARE
    fid INTEGER;
    aid INTEGER;
    today_jst timestamptz;
BEGIN
    SELECT id INTO fid FROM fishermen WHERE name = '佐藤 三郎';
    today_jst := date_trunc('day', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Tokyo') AT TIME ZONE 'Asia/Tokyo';
    -- Get today's auction at Toyosu
    SELECT a.id INTO aid
    FROM auctions a
    JOIN venues v ON a.venue_id = v.id
    WHERE v.name = '豊洲市場' 
      AND a.start_at >= today_jst 
      AND a.start_at < today_jst + interval '1 day'
    LIMIT 1;

    IF fid IS NOT NULL AND aid IS NOT NULL THEN
        INSERT INTO auction_items (fisherman_id, auction_id, fish_type, quantity, unit, sort_order)
        SELECT fid, aid, '本マグロ', 1, '本', 1
        WHERE NOT EXISTS (
            SELECT 1 FROM auction_items
            WHERE fisherman_id = fid AND auction_id = aid AND fish_type = '本マグロ'
        );

        INSERT INTO auction_items (fisherman_id, auction_id, fish_type, quantity, unit, sort_order)
        SELECT fid, aid, '天然真鯛', 20, 'kg', 2
        WHERE NOT EXISTS (
            SELECT 1 FROM auction_items
            WHERE fisherman_id = fid AND auction_id = aid AND fish_type = '天然真鯛'
        );

        INSERT INTO auction_items (fisherman_id, auction_id, fish_type, quantity, unit, sort_order)
        SELECT fid, aid, '生ウニ', 5, '板', 3
        WHERE NOT EXISTS (
            SELECT 1 FROM auction_items
            WHERE fisherman_id = fid AND auction_id = aid AND fish_type = '生ウニ'
        );
    END IF;
END $$;

-- Admins
INSERT INTO admins (email, password_hash)
SELECT 'admin@example.com', '$2a$10$bBUXj3ggxD38hbMPPcXRYe/CiPPfiP8pBv0x593dA5YrlEEjL.AxG'
WHERE NOT EXISTS (SELECT 1 FROM admins WHERE email = 'admin@example.com');

