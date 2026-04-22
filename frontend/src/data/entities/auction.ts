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

export interface AuctionItem {
  id: number;
  auctionId: number;
  fishermanId: number;
  fishType: string;
  quantity: number;
  unit: string;
  highestBid?: number;
  highestBidderId?: number;
  highestBidderName?: string;
  sortOrder: number;
  createdAt: string;
  deletedAt?: string;
}

export interface Bid {
  itemId: number;
  buyerId: number;
  price: number;
}
