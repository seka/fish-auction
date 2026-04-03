import { ItemStatus as EntityItemStatus } from '@entities/auction';

/**
 * アイテムステータス表示用の表示情報を取得する
 */
export const selectItemStatus = (
  status: EntityItemStatus,
): {
  value: 'Pending' | 'Bidding' | 'Sold' | 'Unsold';
  labelKey: string;
  variant: 'success' | 'warning' | 'error' | 'info' | 'neutral';
  isPending: boolean;
  isBidding: boolean;
  isSold: boolean;
  isUnsold: boolean;
} => {
  // Entity ステータスから ViewModel 内部値へのマッピング
  const valueMap: Record<EntityItemStatus, 'Pending' | 'Bidding' | 'Sold' | 'Unsold'> = {
    Pending: 'Pending',
    Bidding: 'Bidding',
    Sold: 'Sold',
    Unsold: 'Unsold',
  };

  const value = valueMap[status] || 'Pending';

  // ViewModel 内部値から表示用設定へのマッピング
  const config: Record<
    typeof value,
    {
      labelKey: string;
      variant: 'success' | 'warning' | 'error' | 'info' | 'neutral';
    }
  > = {
    Pending: { labelKey: 'Pending', variant: 'info' },
    Bidding: { labelKey: 'Bidding', variant: 'success' },
    Sold: { labelKey: 'Sold', variant: 'neutral' },
    Unsold: { labelKey: 'Unsold', variant: 'error' },
  };

  const { labelKey, variant } = config[value];

  return {
    value,
    labelKey,
    variant,
    isPending: value === 'Pending',
    isBidding: value === 'Bidding',
    isSold: value === 'Sold',
    isUnsold: value === 'Unsold',
  };
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
