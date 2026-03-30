import { useQuery } from '@tanstack/react-query';
import { getItemsByAuction } from '@/src/data/api/admin';
import { AuctionItem } from '@entities/auction';
import { itemKeys } from './keys';

export const useItemsByAuction = (auctionId: number) => {
  const {
    data: items = [],
    isLoading,
    refetch,
  } = useQuery<AuctionItem[]>({
    queryKey: itemKeys.publicByAuction(auctionId),
    queryFn: () => getItemsByAuction(auctionId),
    refetchInterval: 5000,
  });

  return { items, isLoading, refetch };
};
