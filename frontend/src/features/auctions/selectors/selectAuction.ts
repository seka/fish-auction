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
 * 表示用に時間をフォーマットする (HH:MM - HH:MM)
 */
export const selectTimeLabel = (startTime: string | null, endTime: string | null): string => {
  const start = selectTime(startTime);
  const end = selectTime(endTime);
  if (!start && !end) return '';
  return `${start || '--:--'} ~ ${end || '--:--'}`;
};

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
 * 次の最低入札額を計算する
 */
export const selectNextMinimumBid = (currentHighestBid: number): number => {
  return currentHighestBid + selectMinimumBidIncrement(currentHighestBid);
};

/**
 * 文字列の日付と時刻を JST として Date オブジェクトに変換する
 */
export const toJSTDate = (date: string, time: string | null): Date => {
  const t = time || '00:00:00';
  // ISO 8601 形式に +09:00 を付与して JST としてパースさせる
  return new Date(`${date.replace(/\//g, '-')}T${t}+09:00`);
};

/**
 * オークションが現在開催中（入札可能時間内）かどうかをチェックする
 */
export const selectIsAuctionActive = (
  auction:
    | Pick<EntityAuction, 'status' | 'auctionDate' | 'startTime' | 'endTime'>
    | {
        status: { value: string };
        duration: {
          startAt: Date;
          endAt: Date;
          dateLabel: string;
          startTime?: string;
          endTime?: string;
        };
      },
  now = new Date(),
): boolean => {
  // ステータスの取得
  const status = typeof auction.status === 'string' ? auction.status : auction.status.value;

  // ステータスが明確に終了または中止なら非アクティブ
  if (status === 'completed' || status === 'cancelled') {
    return false;
  }

  // ステータスが開催中ならアクティブ
  if (status === 'in_progress') {
    return true;
  }

  // ViewModel (durationプロパティを持つ) の場合
  if ('duration' in auction) {
    return now >= auction.duration.startAt && now <= auction.duration.endAt;
  }

  const date = auction.auctionDate;
  const startTime = auction.startTime;
  const endTime = auction.endTime;

  if (!startTime || !endTime) {
    return status === 'scheduled';
  }

  try {
    const startDateTime = toJSTDate(date, startTime);
    const endDateTime = toJSTDate(date, endTime);

    return now >= startDateTime && now <= endDateTime;
  } catch (e) {
    console.error('Failed to parse auction time', e);
    return false;
  }
};

/**
 * 公開一覧用の表示ポリシーを適用し、ソートした結果を返す
 * 中止されていないオークションを表示対象とし、開催中を優先して開始日時順でソートする
 */
export const selectVisiblePublicAuctions = <
  T extends {
    status: {
      isInProgress: boolean;
      isCancelled: boolean;
    };
    duration: { startAt: Date };
  },
>(
  auctions: T[],
): T[] => {
  return [...auctions]
    .filter((a) => !a.status.isCancelled)
    .sort((a, b) => {
      // 開催中を優先する
      if (a.status.isInProgress && !b.status.isInProgress) return -1;
      if (!a.status.isInProgress && b.status.isInProgress) return 1;
      // 日付順でソート
      return a.duration.startAt.getTime() - b.duration.startAt.getTime();
    });
};
