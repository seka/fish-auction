import { describe, it, expect } from 'vitest';
import { toAuction } from './useAuctions';
import { AuctionSummary } from '@/src/data/api/buyer_mypage';

describe('mypage/queries/useAuctions mapping', () => {
  it('should map AuctionSummary to feature Auction model correctly', () => {
    const summary: AuctionSummary = {
      id: 1,
      venueId: 10,
      auctionDate: '2024-03-30',
      startTime: '10:00:00',
      endTime: '12:00:00',
      status: 'in_progress',
      createdAt: '2024-03-01',
      updatedAt: '2024-03-01',
    };

    const result = toAuction(summary);

    expect(result.id).toBe(1);
    expect(result.status).toBe('in_progress');
    expect(result.startTime).toBe('10:00:00');
  });

  it('should handle null startTime in AuctionSummary', () => {
    const summary: AuctionSummary = {
      id: 1,
      venueId: 10,
      auctionDate: '2024-03-30',
      startTime: null,
      endTime: null,
      status: 'scheduled',
      createdAt: '2024-03-01',
      updatedAt: '2024-03-01',
    };

    const result = toAuction(summary);

    expect(result.startTime).toBeNull();
    expect(result.endTime).toBeNull();
  });
});
