import { AuctionStatus } from '@entities/auction';

export interface Auction {
  id: number;
  venueId: number;
  auctionDate: string;
  startTime: string | null;
  endTime: string | null;
  status: AuctionStatus;
}

export interface Venue {
  id: number;
  name: string;
}

export interface AuctionItem {
  id: number;
  auctionId: number;
  fishType: string;
  quantity: number;
  unit: string;
  startPrice: number;
  currentPrice: number;
  status: string;
}
