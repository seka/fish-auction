import { Auction as EntityAuction, AuctionItem as EntityAuctionItem } from '@entities/auction';
import { Auction, AuctionItem } from '../types';

export const toAuction = (entity: EntityAuction): Auction => ({
  id: entity.id,
  venueId: entity.venueId,
  auctionDate: entity.auctionDate,
  startTime: entity.startTime ?? null,
  endTime: entity.endTime ?? null,
  status: entity.status,
});

export const toAuctionItem = (entity: EntityAuctionItem): AuctionItem => ({
  id: entity.id,
  auctionId: entity.auctionId,
  fishermanId: entity.fishermanId,
  fishType: entity.fishType,
  quantity: entity.quantity,
  unit: entity.unit,
  startPrice: 0, 
  currentPrice: entity.highestBid || 0,
  status: entity.status,
  highestBid: entity.highestBid,
  highestBidderId: entity.highestBidderId,
  highestBidderName: entity.highestBidderName,
});
