import { AuctionSummary } from '@entities';

export type AuctionStatus = 'scheduled' | 'in_progress' | 'completed' | 'cancelled';

export interface Auction {
  id: number;
  status: AuctionStatus;
  startAt: string | null;
  endAt: string | null;
}

export const toAuction = (entity: AuctionSummary): Auction => {
  return {
    id: entity.id,
    status: entity.status,
    startAt: entity.startAt ?? null,
    endAt: entity.endAt ?? null,
  };
};
