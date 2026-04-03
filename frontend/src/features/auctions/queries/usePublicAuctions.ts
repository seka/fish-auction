import { useAuctionQuery as useDataAuctionQuery } from '@/src/data/queries/publicAuction/useQuery';
import { toAuction } from '../types/auction';

export const usePublicAuctions = (filters?: { venueId?: number; date?: string; status?: string }) => {
  return useDataAuctionQuery(filters, {
    select: (data) =>
      data
        .map(toAuction)
        .filter((a) =>
          ['scheduled', 'in_progress', 'completed', 'cancelled'].includes(a.status.value),
        )
        .sort((a, b) => {
          if (a.status.isInProgress && !b.status.isInProgress) return -1;
          if (!a.status.isInProgress && b.status.isInProgress) return 1;
          return a.duration.startAt.getTime() - b.duration.startAt.getTime();
        }),
  });
};
