export type AuctionStatus = 'scheduled' | 'in_progress' | 'completed' | 'cancelled';

export interface AuctionLike {
  status: AuctionStatus;
  auctionDate: string;
  startTime?: string | null;
  endTime?: string | null;
}

/**
 * オークションが現在開催中（入札可能時間内）かどうかをチェックする
 */
export const isAuctionActive = (auction: AuctionLike): boolean => {
  // ステータスが明確に終了または中止なら非アクティブ
  if (auction.status === 'completed' || auction.status === 'cancelled') {
    return false;
  }

  // ステータスが開催中ならアクティブ
  if (auction.status === 'in_progress') {
    return true;
  }

  // それ以外（scheduled等）は時刻判定
  // 開始時間または終了時間がない場合は判定不能なため、ステータスのみに頼る
  if (!auction.startTime || !auction.endTime) {
    return auction.status === 'scheduled';
  }

  try {
    const now = new Date();
    // yyyy-mm-dd を yyyy/mm/dd に置換すると多くのブラウザでローカル時刻として解釈されやすい
    const dateStr = auction.auctionDate.replace(/-/g, '/');
    const auctionDate = new Date(dateStr);

    // 開始時間と終了時間をパース (HH:mm)
    const [startHour, startMin] = auction.startTime.split(':').map(Number);
    const [endHour, endMin] = auction.endTime.split(':').map(Number);

    if (isNaN(startHour) || isNaN(endHour)) return false;

    // 開始日時と終了日時オブジェクトを作成
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
