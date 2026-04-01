import { AuctionItem as EntityAuctionItem } from '@entities/auction';

export type ItemStatus = 'Pending' | 'Sold' | 'Unsold' | 'Bidding';

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

export const toAuctionItem = (entity: EntityAuctionItem): AuctionItem => ({
  id: entity.id,
  auctionId: entity.auctionId,
  fishermanId: entity.fishermanId,
  fishType: entity.fishType,
  quantity: entity.quantity,
  unit: entity.unit,
  status: entity.status as ItemStatus,
  highestBid: entity.highestBid,
  highestBidderId: entity.highestBidderId,
  highestBidderName: entity.highestBidderName,
  sortOrder: entity.sortOrder,
  createdAt: entity.createdAt,
});
