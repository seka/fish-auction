import { describe, it, expect } from 'vitest';
import { AuctionItem as EntityAuctionItem } from '@entities/auction';
import { toAuctionItem } from './item';

describe('auctions/types/item', () => {
  describe('toAuctionItem', () => {
    it('should map item entity correctly and convert field names', () => {
      const entity: EntityAuctionItem = {
        id: 101,
        auctionId: 1,
        fishermanId: 5,
        fishType: 'Maguro',
        quantity: 100,
        unit: 'kg',
        status: 'Bidding',
        highestBid: 5000,
        highestBidderId: 2,
        highestBidderName: 'Buyer A',
        sortOrder: 1,
        createdAt: '2024-03-30',
      };

      const result = toAuctionItem(entity);

      expect(result.fishType).toBe('Maguro');
      expect(result.currentPrice).toBe(5000); // highestBid -> currentPrice
      expect(result.status).toBe('Bidding');
    });

    it('should handle missing highestBid', () => {
      const entity: EntityAuctionItem = {
        id: 102,
        auctionId: 1,
        fishermanId: 5,
        fishType: 'Saba',
        quantity: 50,
        unit: 'kg',
        status: 'Pending',
        sortOrder: 2,
      } as EntityAuctionItem;

      const result = toAuctionItem(entity);
      expect(result.currentPrice).toBe(0);
    });
  });
});
