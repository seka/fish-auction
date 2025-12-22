import { Auction } from '@/src/models/auction';

/**
 * オークションが現在開催中（入札可能時間内）かどうかをチェックする
 */
export const isAuctionActive = (auction: Auction): boolean => {
    if (!auction.startTime || !auction.endTime) {
        // 時間が設定されていない場合は常にアクティブとみなす
        return true;
    }

    const now = new Date();
    const auctionDate = new Date(auction.auctionDate);

    // 開始時間と終了時間をパース
    const [startHour, startMin] = auction.startTime.split(':').map(Number);
    const [endHour, endMin] = auction.endTime.split(':').map(Number);

    // 開始日時と終了日時オブジェクトを作成
    const startDateTime = new Date(auctionDate);
    startDateTime.setHours(startHour, startMin, 0, 0);

    const endDateTime = new Date(auctionDate);
    endDateTime.setHours(endHour, endMin, 0, 0);

    return now >= startDateTime && now <= endDateTime;
};

/**
 * 表示用に時間をフォーマットする (HH:MM)
 */
export const formatTime = (time?: string): string => {
    if (!time) return '';
    return time.substring(0, 5); // HH:MM:SS から HH:MM を抽出
};

export const getMinimumBidIncrement = (currentPrice: number): number => {
    if (currentPrice < 1000) return 100;
    if (currentPrice < 10000) return 500;
    if (currentPrice < 100000) return 1000;
    return 5000;
};
