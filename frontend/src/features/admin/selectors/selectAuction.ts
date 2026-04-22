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
