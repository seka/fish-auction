-- 会場とセリにCASCADE削除制約を追加
-- このマイグレーションは外部キー制約を変更し、関連レコードを自動的に削除するようにします。
--
-- このマイグレーション後の削除の挙動:
-- 1. 会場を削除すると、CASCADE削除により以下が削除されます:
--    - その会場に関連付けられたすべてのセリ
--    - それらのセリに関連付けられたすべての出品
--
-- 2. セリを削除すると、CASCADE削除により以下が削除されます:
--    - そのセリに関連付けられたすべての出品
--
-- 注意: 入札（transactions）は出品への参照を持っているため、入札が存在する出品は削除できません。

-- auctions.venue_id の外部キー制約をCASCADEに変更
ALTER TABLE auctions
DROP CONSTRAINT IF EXISTS auctions_venue_id_fkey,
ADD CONSTRAINT auctions_venue_id_fkey 
    FOREIGN KEY (venue_id) 
    REFERENCES venues(id) 
    ON DELETE CASCADE;

-- auction_items.auction_id の外部キー制約をCASCADEに変更
ALTER TABLE auction_items
DROP CONSTRAINT IF EXISTS auction_items_auction_id_fkey,
ADD CONSTRAINT auction_items_auction_id_fkey 
    FOREIGN KEY (auction_id) 
    REFERENCES auctions(id) 
    ON DELETE CASCADE;
