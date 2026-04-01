import { AuctionItem as EntityAuctionItem } from '@entities/auction';

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

export const toAuctionItem = (entity: EntityAuctionItem): AuctionItem => ({
  id: entity.id,
  auctionId: entity.auctionId,
  fishermanId: entity.fishermanId,
  fishType: entity.fishType,
  quantity: entity.quantity,
  unit: entity.unit,
  startPrice: 0,
  currentPrice: entity.highestBid || 0,
  status: entity.status as ItemStatus,
  highestBid: entity.highestBid,
  highestBidderId: entity.highestBidderId,
  highestBidderName: entity.highestBidderName,
});
