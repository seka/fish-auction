import { useQuery } from '@tanstack/react-query';
import { getAuction } from '@/src/api/auction';
import { getItemsByAuction } from '@/src/api/admin';
import { AuctionItem } from '@/src/models/auction';
import { auctionKeys } from '@/src/hooks/auction/keys';

export const useAuctionData = (auctionId: number) => {
  const { data: auction, isLoading: isAuctionLoading } = useQuery({
    queryKey: auctionKeys.publicDetail(auctionId),
    queryFn: () => getAuction(auctionId),
    refetchInterval: 5000,
  });

  const {
    data: items = [],
    isLoading: isItemsLoading,
    refetch: refetchItems,
  } = useQuery<AuctionItem[]>({
    queryKey: auctionKeys.publicItems(auctionId),
    queryFn: () => getItemsByAuction(auctionId),
    refetchInterval: 5000, // Poll every 5 seconds
  });

  return {
    auction,
    items,
    isLoading: isAuctionLoading || isItemsLoading,
    refetchItems,
  };
};
