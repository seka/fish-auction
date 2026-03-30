import { describe, it, expect } from 'vitest';
import { toAuctionItem } from './useItems';
import { AuctionItem as EntityAuctionItem } from '@entities/auction';

describe('admin/queries/useItems mapping', () => {
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
    expect(result.fishType).toBe('Saba');
    expect(result.createdAt).toBe('2024-03-30T10:00:00Z');
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

    expect(result.highestBid).toBeUndefined();
    expect(result.id).toBe(102);
  });
});
