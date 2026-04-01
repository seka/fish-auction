import { Auction as EntityAuction } from '@entities/auction';
import { selectIsAuctionActive } from '../selectors/selectAuction';

export type AuctionStatus = 'scheduled' | 'in_progress' | 'completed' | 'cancelled';

export interface Auction {
  id: number;
  venueId: number;
  auctionDate: string;
  startTime: string | null;
  endTime: string | null;
  status: AuctionStatus;
  isActive: boolean;
}

export const toAuction = (entity: EntityAuction): Auction => {
  return {
    id: entity.id,
    venueId: entity.venueId,
    auctionDate: entity.auctionDate,
    startTime: entity.startTime ?? null,
    endTime: entity.endTime ?? null,
    status: entity.status,
    isActive: selectIsAuctionActive(entity),
  };
};
