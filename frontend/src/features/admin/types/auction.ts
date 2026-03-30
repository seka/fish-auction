export type AuctionStatus = 'scheduled' | 'in_progress' | 'completed' | 'cancelled';

export interface Auction {
  id: number;
  venueId: number;
  auctionDate: string;
  startTime?: string;
  endTime?: string;
  status: AuctionStatus;
  createdAt: string;
  updatedAt: string;
}
