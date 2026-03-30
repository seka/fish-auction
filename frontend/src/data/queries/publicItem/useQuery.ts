import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { getItemsByAuction } from '@/src/data/api/admin';
import { AuctionItem } from '@entities/auction';
import { itemKeys } from './keys';

export const useItemsByAuction = <T = AuctionItem[]>(
  auctionId: number,
  options?: Omit<UseQueryOptions<AuctionItem[], Error, T>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: itemKeys.publicByAuction(auctionId),
    queryFn: () => getItemsByAuction(auctionId),
    refetchInterval: 5000,
    ...options,
  });
};
