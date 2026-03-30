export type AuctionStatus = 'scheduled' | 'in_progress' | 'completed' | 'cancelled';

export interface Auction {
  id: number;
  venueId: number;
  auctionDate: string;
  startTime: string | null;
  endTime: string | null;
  status: AuctionStatus;
  isActive?: boolean;
}
