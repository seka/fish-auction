import { AuctionStatus } from '@/src/types/auction';

export interface Auction {
  id: number;
  venueId: number;
  auctionDate: string;
  startTime: string | null;
  endTime: string | null;
  status: AuctionStatus;
  isActive?: boolean;
}
