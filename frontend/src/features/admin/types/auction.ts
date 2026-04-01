import { Auction as EntityAuction } from '@entities/auction';

export type AuctionStatus = 'scheduled' | 'in_progress' | 'completed' | 'cancelled';

export interface Auction {
  id: number;
  venueId: number;
  auctionDate: string;
  startTime: string | null;
  endTime: string | null;
  status: AuctionStatus;
  createdAt: string;
  updatedAt: string;
}

export const toAuction = (entity: EntityAuction): Auction => {
  return {
    id: entity.id,
    venueId: entity.venueId,
    auctionDate: entity.auctionDate,
    startTime: entity.startTime ?? null,
    endTime: entity.endTime ?? null,
    status: entity.status,
    createdAt: entity.createdAt,
    updatedAt: entity.updatedAt,
  };
};
