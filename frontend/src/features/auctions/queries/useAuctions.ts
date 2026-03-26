import { useAuctionDetailQuery } from '@/src/data/queries/publicAuction/useQuery';
import { useItemsByAuction } from '@/src/data/queries/publicItem/useQuery';
import { useBidMutation } from '@/src/data/queries/buyerAuction/useMutation';

import { isAuctionActive } from '@/src/utils/auction';

/**
 * オークション詳細クエリフック
 * select を用いて、ビジネスロジックを反映したモデルへ変換する例
 */
export const useAuctionDetailData = (auctionId: number) => {
  const { data: auctionInfo, isLoading: isAuctionLoading } = useAuctionDetailQuery(auctionId, {
    select: (data) => ({
      ...data,
      isActive: isAuctionActive(data),
    }),
  });
  
  const { items, isLoading: isItemsLoading, refetch: refetchItems } = useItemsByAuction(auctionId);
  
  return {
    auction: auctionInfo,
    items,
    isLoading: isAuctionLoading || isItemsLoading,
    refetchItems,
  };
};

/**
 * 入札ミューテーションフック
 */
export const useBidSubmit = () => {
  return useBidMutation();
};
