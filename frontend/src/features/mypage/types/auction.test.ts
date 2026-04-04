import { describe, it, expect } from 'vitest';
import { AuctionSummary } from '@entities';
import { toAuction } from './auction';

describe('mypage/types/auction', () => {
  describe('toAuction', () => {
    it('should map member entity to mypage model correctly', () => {
      const entity: AuctionSummary = {
        id: 1,
        venueId: 10,
        status: 'in_progress',
        auctionDate: '2024-03-30',
        startTime: '10:00:00',
        endTime: '12:00:00',
        createdAt: '2024-03-01T00:00:00Z',
        updatedAt: '2024-03-01T00:00:00Z',
      };

      const result = toAuction(entity);

      expect(result.id).toBe(1);
      expect(result.status).toBe('in_progress');
      expect(result.startTime).toBe('10:00:00');
    });

    it('should handle null startTime', () => {
      const entity: AuctionSummary = {
        id: 2,
        venueId: 10,
        status: 'scheduled',
        auctionDate: '2024-03-30',
        startTime: null,
        endTime: null,
        createdAt: '2024-03-01T00:00:00Z',
        updatedAt: '2024-03-01T00:00:00Z',
      };

      const result = toAuction(entity);

      expect(result.startTime).toBeNull();
      expect(result.endTime).toBeNull();
    });
  });
});
