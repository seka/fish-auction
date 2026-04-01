import { describe, it, expect } from 'vitest';
import { toAuction, toAuctionItem } from './transformers';
import { Auction as EntityAuction, AuctionItem as EntityAuctionItem } from '@entities/auction';

describe('auctions/queries/useAuctions mapping', () => {
  describe('toAuction', () => {
    it('should map entity to feature model correctly', () => {
      const entity: EntityAuction = {
        id: 1,
        venueId: 10,
        auctionDate: '2024-03-30',
        startTime: '10:00:00',
        endTime: '12:00:00',
        status: 'in_progress',
        createdAt: '2024-03-01',
        updatedAt: '2024-03-01',
      };

      const result = toAuction(entity);

      expect(result.id).toBe(1);
      expect(result.startTime).toBe('10:00:00');
      expect(result.status).toBe('in_progress');
    });

    it('should convert undefined/null startTime to null', () => {
      const entity = {
        id: 2,
        venueId: 10,
        auctionDate: '2024-03-30',
        status: 'scheduled',
        createdAt: '2024-03-01',
        updatedAt: '2024-03-01',
      } as EntityAuction;

      const result = toAuction(entity);

      expect(result.startTime).toBeNull();
      expect(result.endTime).toBeNull();
    });
  });

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
