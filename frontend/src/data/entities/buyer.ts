import { AuctionStatus } from './auction';

export interface Purchase {
  id: number;
  itemId: number;
  fishType: string;
  quantity: number;
  unit: string;
  price: number;
  buyerId: number;
  auctionId: number;
  auctionDate: string;
  createdAt: string;
}

export interface AuctionSummary {
  id: number;
  venueId: number;
  startAt: string | null;
  endAt: string | null;
  status: AuctionStatus;
  createdAt: string;
  updatedAt: string;
}
