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
      expect(result.status.value).toBe('Bidding');
      expect(result.status.variant).toBe('success');
      expect(result.bidding.highestBid).toBe(5000);
      expect(result.bidding.highestBidderId).toBe(2);
      expect(result.bidding.highestBidderName).toBe('Buyer A');
      expect(result.bidding.nextMinBid.value).toBe(5500);
      expect(result.bidding.nextMinBid.label).toBe('¥5,500');
      expect(result.quantity.value).toBe(100);
      expect(result.quantity.label).toBe('100 kg');
      expect(result.price.label).toBe('¥5,000');
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
      } as EntityAuctionItem;

      const result = toAuctionItem(entity);
      expect(result.bidding.highestBid).toBeNull();
      expect(result.bidding.nextMinBid.value).toBe(100);
      expect(result.price.label).toBe('¥0');
      expect(result.quantity.label).toBe('50 kg');
    });
  });
});
