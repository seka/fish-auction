import { describe, it, expect } from 'vitest';
import { Auction as EntityAuction } from '@entities/auction';
import { toAuction } from './auction';

describe('auctions/types/auction', () => {
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
});
