import { AuctionSummary } from '@entities';

export type AuctionStatus = 'scheduled' | 'in_progress' | 'completed' | 'cancelled';

export interface Auction {
  id: number;
  status: AuctionStatus;
  auctionDate: string;
  startTime: string | null;
  endTime: string | null;
}

export const toAuction = (entity: AuctionSummary): Auction => {
  return {
    id: entity.id,
    status: entity.status,
    auctionDate: entity.auctionDate,
    startTime: entity.startTime ?? null,
    endTime: entity.endTime ?? null,
  };
};
