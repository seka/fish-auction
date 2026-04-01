import { Auction as EntityAuction } from '@entities/auction';

/**
 * 表示用に時間をフォーマットする (HH:MM)
 */
export const selectTime = (time?: string | null): string => {
  if (!time) return '';
  return time.substring(0, 5); // HH:MM:SS から HH:MM を抽出
};

/**
 * 価格に応じた最小入札単位を取得する
 */
export const selectMinimumBidIncrement = (currentPrice: number): number => {
  if (currentPrice < 1000) return 100;
  if (currentPrice < 10000) return 500;
  if (currentPrice < 100000) return 1000;
  return 5000;
};

/**
 * オークションが現在開催中（入札可能時間内）かどうかをチェックする
 */
export const selectIsAuctionActive = (
  auction: Pick<EntityAuction, 'status' | 'auctionDate' | 'startTime' | 'endTime'>,
  now = new Date(),
): boolean => {
  // ステータスが明確に終了または中止なら非アクティブ
  if (auction.status === 'completed' || auction.status === 'cancelled') {
    return false;
  }

  // ステータスが開催中ならアクティブ
  if (auction.status === 'in_progress') {
    return true;
  }

  // それ以外（scheduled等）は時刻判定
  if (!auction.startTime || !auction.endTime) {
    return auction.status === 'scheduled';
  }

  try {
    const dateStr = auction.auctionDate.replace(/-/g, '/');
    const auctionDate = new Date(dateStr);

    const [startHour, startMin] = auction.startTime.split(':').map(Number);
    const [endHour, endMin] = auction.endTime.split(':').map(Number);

    if (isNaN(startHour) || isNaN(endHour)) return false;

    const startDateTime = new Date(auctionDate);
    startDateTime.setHours(startHour, startMin, 0, 0);

    const endDateTime = new Date(auctionDate);
    endDateTime.setHours(endHour, endMin, 0, 0);

    return now >= startDateTime && now <= endDateTime;
  } catch (e) {
    console.error('Failed to parse auction time', e);
    return false;
  }
};
