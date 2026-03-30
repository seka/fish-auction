import { useAuctionDetailQuery } from '@/src/data/queries/publicAuction/useQuery';
import { useItemsByAuction } from '@/src/data/queries/publicItem/useQuery';
import { useBidMutation } from '@/src/data/queries/buyerAuction/useMutation';

import { Auction as EntityAuction, AuctionItem as EntityAuctionItem } from '@entities/auction';
import { Auction, AuctionItem } from '../types';
import { isAuctionActive } from '@/src/utils/auction';

const toAuction = (entity: EntityAuction): Auction => ({
  id: entity.id,
  venueId: entity.venueId,
  auctionDate: entity.auctionDate,
  startTime: entity.startTime ?? null,
  endTime: entity.endTime ?? null,
  status: entity.status,
});

const toAuctionItem = (entity: EntityAuctionItem): AuctionItem => ({
  id: entity.id,
  auctionId: entity.auctionId,
  fishermanId: entity.fishermanId,
  fishType: entity.fishType,
  quantity: entity.quantity,
  unit: entity.unit,
  startPrice: 0, 
  currentPrice: entity.highestBid || 0,
  status: entity.status,
  highestBid: entity.highestBid,
  highestBidderId: entity.highestBidderId,
  highestBidderName: entity.highestBidderName,
});

/**
 * オークション詳細クエリフック
 * select を用いて、Feature 用のドメインモデルに変換する
 */
export const useAuctionDetailData = (auctionId: number) => {
  const { data: auction, isLoading: isAuctionLoading } = useAuctionDetailQuery(auctionId, {
    select: (data) => ({
      ...toAuction(data),
      isActive: isAuctionActive(data),
    }),
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
