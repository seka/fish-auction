import type { ItemStatus } from '@/src/types/item';

export type { ItemStatus };

export interface AuctionItem {
  id: number;
  auctionId: number;
  fishermanId: number;
  fishType: string;
  quantity: number;
  unit: string;
  status: ItemStatus;
  highestBid?: number;
  highestBidderId?: number;
  highestBidderName?: string;
  sortOrder: number;
  createdAt: string;
}
