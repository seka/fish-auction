import type { AuctionStatus } from '@/src/types/auction';

export type { AuctionStatus };

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
