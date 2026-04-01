import { useMyPurchases as useDataPurchases } from '@/src/data/queries/buyerPurchase/useQuery';
import { toPurchase } from '../types/purchase';

/**
 * マイページ用 購入履歴クエリフック
 */
export const usePurchases = () => {
  const { data: purchases, ...rest } = useDataPurchases({
    select: (data) => data.map(toPurchase),
  });

  return {
    ...rest,
    purchases: purchases || [],
  };
};
