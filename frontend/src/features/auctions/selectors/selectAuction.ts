import { Auction as EntityAuction, AuctionStatus as EntityAuctionStatus } from '@entities/auction';

/**
 * オークションステータス表示用の表示情報を取得する
 */
export const selectAuctionStatus = (
  status: EntityAuctionStatus,
): {
  value: 'scheduled' | 'in_progress' | 'completed' | 'cancelled';
  labelKey: string;
  variant: 'success' | 'warning' | 'error' | 'info' | 'neutral';
  isScheduled: boolean;
  isInProgress: boolean;
  isCompleted: boolean;
  isCancelled: boolean;
} => {
  // Entity ステータスから ViewModel 内部値へのマッピング
  const valueMap: Record<
    EntityAuctionStatus,
    'scheduled' | 'in_progress' | 'completed' | 'cancelled'
  > = {
    scheduled: 'scheduled',
    in_progress: 'in_progress',
    completed: 'completed',
    cancelled: 'cancelled',
  };

  const value = valueMap[status] || 'scheduled';

  // ViewModel 内部値から表示用設定へのマッピング
  const config: Record<
    typeof value,
    {
      labelKey: string;
      variant: 'success' | 'warning' | 'error' | 'info' | 'neutral';
    }
  > = {
    scheduled: { labelKey: 'scheduled', variant: 'info' },
    in_progress: { labelKey: 'in_progress', variant: 'success' },
    completed: { labelKey: 'completed', variant: 'neutral' },
    cancelled: { labelKey: 'cancelled', variant: 'error' },
  };

  const { labelKey, variant } = config[value];

  return {
    value,
    labelKey,
    variant,
    isScheduled: value === 'scheduled',
    isInProgress: value === 'in_progress',
    isCompleted: value === 'completed',
    isCancelled: value === 'cancelled',
  };
};

/**
 * 表示用に時間をフォーマットする (HH:MM ~ HH:MM)
 */
export const selectTimeLabel = (startAt: Date | null, endAt: Date | null): string => {
  const format = (d: Date) =>
    d.toLocaleTimeString('ja-JP', {
      hour: '2-digit',
      minute: '2-digit',
      timeZone: 'Asia/Tokyo',
    });
  if (!startAt && !endAt) return '';
  return `${startAt ? format(startAt) : '--:--'} ~ ${endAt ? format(endAt) : '--:--'}`;
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
 * 次の最低入札額を計算する
 */
export const selectNextMinimumBid = (currentHighestBid: number): number => {
  return currentHighestBid + selectMinimumBidIncrement(currentHighestBid);
};

/**
 * オークションが現在開催中（入札可能時間内）かどうかをチェックする
 */
export const selectIsAuctionActive = (
  auction: Pick<EntityAuction, 'status' | 'startAt' | 'endAt'>,
  now = new Date(),
): boolean => {
  const { status, startAt, endAt } = auction;

  if (status !== 'in_progress') {
    return false;
  }

  if (!startAt || !endAt) {
    return false;
  }

  try {
    return now >= new Date(startAt) && now <= new Date(endAt);
  } catch (e) {
    console.error('Failed to parse auction time', e);
    return false;
  }
};

/**
 * 公開一覧用の表示ポリシーを適用し、ソートした結果を返す
 * 表示対象ステータス (scheduled / in_progress / completed) に限定し、開催中を優先して開始日時順でソートする
 */
export const selectVisiblePublicAuctions = (
  auctions: EntityAuction[],
  now = new Date(),
): EntityAuction[] => {
  return [...auctions]
    .filter(
      (a) => a.status === 'scheduled' || a.status === 'in_progress' || a.status === 'completed',
    )
    .sort((a, b) => {
      const aActive = selectIsAuctionActive(a, now);
      const bActive = selectIsAuctionActive(b, now);

      if (aActive && !bActive) return -1;
      if (!aActive && bActive) return 1;

      const aStart = a.startAt ? new Date(a.startAt).getTime() : 0;
      const bStart = b.startAt ? new Date(b.startAt).getTime() : 0;
      return aStart - bStart;
    });
};
