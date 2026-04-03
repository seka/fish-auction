import { useAuctionQuery as useDataAuctionQuery } from '@/src/data/queries/publicAuction/useQuery';
import { toAuction } from '../types/auction';
import { selectVisiblePublicAuctions } from '../selectors/selectAuction';

export const usePublicAuctions = (filters?: {
  venueId?: number;
  date?: string;
  status?: string;
}) => {
  return useDataAuctionQuery(filters, {
    select: (data) => selectVisiblePublicAuctions(data).map(toAuction),
  });
};
