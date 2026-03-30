import { AuctionItem as EntityAuctionItem } from '@entities/auction';
import { useItemQuery } from '@/src/data/queries/adminItem/useQuery';
import { useItemMutation } from '@/src/data/queries/adminItem/useMutation';
import { AuctionItem } from '../types/item';

export const toAuctionItem = (entity: EntityAuctionItem): AuctionItem => ({
  id: entity.id,
  auctionId: entity.auctionId,
  fishermanId: entity.fishermanId,
  fishType: entity.fishType,
  quantity: entity.quantity,
  unit: entity.unit,
  status: entity.status,
  highestBid: entity.highestBid,
  highestBidderId: entity.highestBidderId,
  highestBidderName: entity.highestBidderName,
  sortOrder: entity.sortOrder,
  createdAt: entity.createdAt,
});

/**
 * 管理画面用商品クエリフック
 */
export const useAdminItems = (auctionId?: number) => {
  const { data: items, ...rest } = useItemQuery(auctionId, {
    select: (data) => data.map(toAuctionItem),
  });

  return {
    ...rest,
    items: items || [],
  };
};

/**
 * 管理画面用商品ミューテーションフック
 */
export const useAdminItemMutations = () => {
  return useItemMutation();
};
