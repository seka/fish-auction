import { AuctionStatus as EntityAuctionStatus } from '@entities/auction';
import { Auction } from '../types/auction';

/**
 * オークションステータス表示用の表示情報を取得する
 */
export const selectAuctionStatus = (status: EntityAuctionStatus): Auction['status'] => {
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
    'scheduled' | 'in_progress' | 'completed' | 'cancelled',
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
 * 文字列の日付と時刻を JST として Date オブジェクトに変換する
 */
export const toJSTDate = (date: string, time: string | null): Date => {
  const t = time || '00:00:00';
  // ISO 8601 形式に +09:00 を付与して JST としてパースさせる
  return new Date(`${date.replace(/\//g, '-')}T${t}+09:00`);
};
