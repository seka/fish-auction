import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { getAuctions, getAuction, getAuctionItems } from '@/src/data/api/auction';
import { adminAuctionKeys } from './keys';
import { Auction, AuctionItem } from '@entities/auction';

export const useAuctionQuery = <T = Auction[]>(
  filters?: { venueId?: number; date?: string; status?: string },
  options?: Omit<UseQueryOptions<Auction[], Error, T>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: adminAuctionKeys.list(filters),
    queryFn: () => getAuctions(filters),
    ...options,
  });
};

export const useAuctionDetailQuery = <T = Auction>(
  id: number,
  options?: Omit<UseQueryOptions<Auction, Error, T>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: adminAuctionKeys.detail(id),
    queryFn: () => getAuction(id),
    enabled: !!id,
    ...options,
  });
};

export const useAuctionItemsQuery = <T = AuctionItem[]>(
  id: number,
  options?: Omit<UseQueryOptions<AuctionItem[], Error, T>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: adminAuctionKeys.items(id),
    queryFn: () => getAuctionItems(id),
    enabled: !!id,
    ...options,
  });
};
