import { useQuery } from '@tanstack/react-query';
import { getAuctions, getAuction } from '@/src/data/api/auction';
import { auctionKeys } from './keys';

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

export const useAuctionDetailQuery = (auctionId: number) => {
  const {
    data: auction,
    isLoading,
    error,
  } = useQuery({
    queryKey: auctionKeys.publicDetail(auctionId),
    queryFn: () => getAuction(auctionId),
    refetchInterval: 5000,
  });

  return { auction, isLoading, error };
};
