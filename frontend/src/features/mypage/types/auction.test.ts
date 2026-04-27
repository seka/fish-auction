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
        startAt: '2024-03-30T10:00:00+09:00',
        endAt: '2024-03-30T12:00:00+09:00',
        createdAt: '2024-03-01T00:00:00Z',
        updatedAt: '2024-03-01T00:00:00Z',
      };

      const result = toAuction(entity);

      expect(result.id).toBe(1);
      expect(result.status).toBe('in_progress');
      expect(result.startAt).toBe('2024-03-30T10:00:00+09:00');
    });

    it('should handle null startAt', () => {
      const entity: AuctionSummary = {
        id: 2,
        venueId: 10,
        status: 'scheduled',
        createdAt: '2024-03-01T00:00:00Z',
        updatedAt: '2024-03-01T00:00:00Z',
      } as AuctionSummary;

      const result = toAuction(entity);

      expect(result.startAt).toBeNull();
      expect(result.endAt).toBeNull();
    });
  });
});
