import type { AuctionStatus } from '@/src/types/auction';

export type { AuctionStatus };

export interface Auction {
  id: number;
  status: AuctionStatus;
  auctionDate: string;
  startTime: string | null;
  endTime: string | null;
}
