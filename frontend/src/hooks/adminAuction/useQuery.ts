import { useQuery } from '@tanstack/react-query';
import { getAuctions, getAuction, getAuctionItems } from '@/src/api/auction';
import { adminAuctionKeys } from './keys';

export const useAuctionQuery = (filters?: { venueId?: number; date?: string; status?: string }) => {
  const {
    data: auctions,
    isLoading,
    error,
  } = useQuery({
    queryKey: adminAuctionKeys.list(filters),
    queryFn: () => getAuctions(filters),
  });

  return { auctions: auctions || [], isLoading, error };
};

export const useAuctionDetailQuery = (id: number) => {
  return useQuery({
    queryKey: adminAuctionKeys.detail(id),
    queryFn: () => getAuction(id),
    enabled: !!id,
  });
};

export const useAuctionItemsQuery = (id: number) => {
  return useQuery({
    queryKey: adminAuctionKeys.items(id),
    queryFn: () => getAuctionItems(id),
    enabled: !!id,
  });
};
