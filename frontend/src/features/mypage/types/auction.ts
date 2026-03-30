export type AuctionStatus = 'scheduled' | 'in_progress' | 'completed' | 'cancelled';

export interface Auction {
  id: number;
  status: AuctionStatus;
  auctionDate: string;
  startTime?: string | null;
  endTime?: string | null;
}
