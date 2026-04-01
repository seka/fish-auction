import { Purchase as EntityPurchase } from '@/src/data/api/buyer_mypage';

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

export const toPurchase = (entity: EntityPurchase): Purchase => ({
  id: entity.id,
  itemId: entity.itemId,
  fishType: entity.fishType,
  quantity: entity.quantity,
  unit: entity.unit,
  price: entity.price,
  auctionId: entity.auctionId,
  auctionDate: entity.auctionDate,
  createdAt: entity.createdAt,
});
