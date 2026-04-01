import { useAuctionQuery as useDataAuctionQuery } from '@/src/data/queries/publicAuction/useQuery';
import { toAuction } from '../types/auction';

export const usePublicAuctions = (filters?: { venueId?: number; date?: string; status?: string }) => {
  return useDataAuctionQuery(filters, {
    select: (data) => data.map(toAuction),
  });
};
