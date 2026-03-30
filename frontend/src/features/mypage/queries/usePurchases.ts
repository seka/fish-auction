import { useMyPurchases as useDataPurchases } from '@/src/data/queries/buyerPurchase/useQuery';
import { Purchase as EntityPurchase } from '@/src/data/api/buyer_mypage';
import { Purchase } from '../types/purchase';

export const toPurchase = (entity: EntityPurchase): Purchase => ({
  id: entity.id,
  itemId: entity.itemId,
  fishType: entity.fishType,
  quantity: entity.quantity,
  unit: entity.unit,
  price: entity.price,
  auctionId: entity.auctionId,
  auctionDate: entity.auctionDate,
  createdAt: entity.createdAt,
});

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
