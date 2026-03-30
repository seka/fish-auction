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
