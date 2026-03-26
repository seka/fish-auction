import { AuctionStatus } from '@/src/models/auction';

export interface Purchase {
  id: number;
  itemId: number;
  createdAt: string;
  fishType: string;
  quantity: number;
  unit: string;
  auctionId: number;
  auctionDate: string;
  price: number;
}

export interface Auction {
  id: number;
  status: AuctionStatus;
  auctionDate: string;
  startTime?: string | null;
  endTime?: string | null;
}


export interface PasswordMessage {
  text: string;
  type: 'info' | 'error' | 'success';
}

export type MyPageTab = 'purchases' | 'auctions' | 'settings';
