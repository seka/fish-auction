import { useAuctionDetailQuery } from '@/src/data/queries/publicAuction/useQuery';
import { useItemsByAuction } from '@/src/data/queries/publicItem/useQuery';
import { useBidMutation } from '@/src/data/queries/buyerAuction/useMutation';
import { toAuction } from '../types/auction';
import { toAuctionItem } from '../types/item';

/**
 * オークション詳細クエリフック
 * select を用いて、Feature 用のドメインモデルに変換する
 */
export const useAuctionDetailData = (auctionId: number) => {
  const { data: auction, isLoading: isAuctionLoading } = useAuctionDetailQuery(auctionId, {
    select: toAuction,
  });
  
  const { data: items, isLoading: isItemsLoading, refetch: refetchItems } = useItemsByAuction(auctionId, {
    select: (data) => data.map(toAuctionItem),
  });
  
  return {
    auction,
    items: items || [],
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
