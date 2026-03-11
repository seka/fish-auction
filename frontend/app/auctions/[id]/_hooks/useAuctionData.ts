import { useQuery } from '@tanstack/react-query';
import { getAuction, getAuctionItems } from '@/src/api/auction';
import { AuctionItem } from '@/src/models/auction';

const EMPTY_ITEMS: AuctionItem[] = [];

export const useAuctionData = (auctionId: number) => {
  const { data: auction, isLoading: isAuctionLoading } = useQuery({
    queryKey: ['auction', auctionId],
    queryFn: () => getAuction(auctionId),
    refetchInterval: 5000,
  });

  const {
    data: items,
    isLoading: isItemsLoading,
    refetch: refetchItems,
  } = useQuery({
    queryKey: ['auction_items', auctionId],
    queryFn: () => getAuctionItems(auctionId),
    refetchInterval: 5000, // Poll every 5 seconds
  });

  return {
    auction,
    items: items || EMPTY_ITEMS,
    isLoading: isAuctionLoading || isItemsLoading,
    refetchItems,
  };
};
