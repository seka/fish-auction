import { useQuery } from '@tanstack/react-query';
import { getAuctions } from '@/src/api/auction';
import { auctionKeys } from './keys';

export const useAuctionQuery = (filters?: { venueId?: number; date?: string; status?: string }) => {
  const {
    data: auctions,
    isLoading,
    error,
  } = useQuery({
    queryKey: auctionKeys.list(filters),
    queryFn: () => getAuctions(filters),
  });

  return { auctions: auctions || [], isLoading, error };
};
