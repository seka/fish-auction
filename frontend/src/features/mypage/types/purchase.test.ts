import { describe, it, expect } from 'vitest';
import { Purchase as EntityPurchase } from '@entities';
import { toPurchase } from './purchase';

describe('mypage/types/purchase', () => {
  describe('toPurchase', () => {
    it('should map EntityPurchase to Mypage Purchase model correctly', () => {
      const entity: EntityPurchase = {
        id: 1,
        itemId: 101,
        fishType: 'Tai',
        quantity: 10,
        unit: 'kg',
        price: 12000,
        buyerId: 1,
        auctionId: 5,
        auctionDate: '2024-03-30',
        createdAt: '2024-03-30T10:00:00Z',
      };

      const result = toPurchase(entity);

      expect(result.id).toBe(1);
      expect(result.fishType).toBe('Tai');
      expect(result.price).toBe(12000);
      expect(result.auctionId).toBe(5);
    });
  });
});
