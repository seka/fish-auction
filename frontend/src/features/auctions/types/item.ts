export type ItemStatus = 'Pending' | 'Sold' | 'Unsold' | 'Bidding';

export interface AuctionItem {
  id: number;
  auctionId: number;
  fishermanId: number;
  fishType: string;
  quantity: number;
  unit: string;
  startPrice: number;
  currentPrice: number;
  status: ItemStatus;
  highestBid?: number;
  highestBidderId?: number;
  highestBidderName?: string;
}
