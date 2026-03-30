import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { getAuctions, getAuction } from '@/src/data/api/auction';
import { auctionKeys } from './keys';
import { Auction } from '@entities/auction';

export const useAuctionQuery = (filters?: { venueId?: number; date?: string; status?: string }) => {
  const {
    data: auctions,
    isLoading,
    error,
  } = useQuery({
    queryKey: auctionKeys.publicList(filters),
    queryFn: () => getAuctions(filters),
  });

  return { auctions: auctions || [], isLoading, error };
};

export const useAuctionDetailQuery = <T = Auction>(
  auctionId: number,
  options?: Omit<UseQueryOptions<Auction, Error, T>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: auctionKeys.publicDetail(auctionId),
    queryFn: () => getAuction(auctionId),
    refetchInterval: 5000,
    ...options,
  });
};
