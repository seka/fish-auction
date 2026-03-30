import { useAuctionQuery as useDataAuctionQuery } from '@/src/data/queries/publicAuction/useQuery';
import { Auction as EntityAuction } from '@entities/auction';
import { Auction } from '../types';

export const toAuction = (entity: EntityAuction): Auction => ({
  id: entity.id,
  venueId: entity.venueId,
  auctionDate: entity.auctionDate,
  startTime: entity.startTime ?? null,
  endTime: entity.endTime ?? null,
  status: entity.status,
});

export const usePublicAuctions = (filters?: { venueId?: number; date?: string; status?: string }) => {
  return useDataAuctionQuery(filters, {
    select: (data) => data.map(toAuction),
  });
};
