import { describe, it, expect } from 'vitest';
import { AuctionItem as EntityAuctionItem } from '@entities/auction';
import { toAuctionItem } from './item';

describe('admin/types/item', () => {
  describe('toAuctionItem', () => {
    it('should map EntityAuctionItem to Admin AuctionItem model correctly', () => {
      const entity: EntityAuctionItem = {
        id: 101,
        auctionId: 1,
        fishermanId: 5,
        fishType: 'Saba',
        quantity: 100,
        unit: 'kg',
        status: 'Bidding',
        highestBid: 2000,
        highestBidderId: 2,
        highestBidderName: 'Buyer A',
        sortOrder: 1,
        createdAt: '2024-03-30T10:00:00Z',
      };

      const result = toAuctionItem(entity);

      expect(result.id).toBe(101);
      expect(result.status).toEqual({
        value: 'Bidding',
        labelKey: 'Bidding',
        variant: 'success',
        isPending: false,
        isBidding: true,
        isSold: false,
        isUnsold: false,
      });
      expect(result.bidding).toEqual({
        highestBid: 2000,
        highestBidderId: 2,
        highestBidderName: 'Buyer A',
        nextMinBid: {
          value: 2500,
          label: '¥2,500',
        },
      });
      expect(result.quantity).toEqual({
        value: 100,
        label: '100 kg',
      });
      expect(result.price).toEqual({
        value: 2000,
        label: '¥2,000',
      });
      expect(result.fishType).toBe('Saba');
    });

    it('should handle missing highestBid in admin item mapping', () => {
      const entity = {
        id: 102,
        auctionId: 1,
        fishType: 'Maguro',
        quantity: 50,
        unit: 'kg',
        status: 'Pending',
        sortOrder: 2,
      } as EntityAuctionItem;

      const result = toAuctionItem(entity);

      expect(result.bidding.highestBid).toBeNull();
      expect(result.bidding.nextMinBid.value).toBe(100);
      expect(result.price.label).toBe('¥0');
      expect(result.quantity.label).toBe('50 kg');
    });
  });
});
