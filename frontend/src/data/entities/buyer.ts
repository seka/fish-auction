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
  auctionDate: string;
  startTime: string | null;
  endTime: string | null;
  status: AuctionStatus;
  createdAt: string;
  updatedAt: string;
}
