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
      expect(result.status).toEqual({
        value: 'in_progress',
        labelKey: 'in_progress',
        variant: 'success',
        isScheduled: false,
        isInProgress: true,
        isCompleted: false,
        isCancelled: false,
      });
      expect(result.duration.startAt).toBeInstanceOf(Date);
      expect(result.duration.endAt).toBeInstanceOf(Date);
      // JST (UTC+9) 10:00 -> UTC 01:00
      expect(result.duration.startAt.toISOString()).toBe('2024-03-30T01:00:00.000Z');
      // JST (UTC+9) 12:00 -> UTC 03:00
      expect(result.duration.endAt.toISOString()).toBe('2024-03-30T03:00:00.000Z');
      expect(result.duration.label).toBe('10:00 ~ 12:00');
    });

    it('should handle undefined/null startTime/endTime', () => {
      const entity = {
        id: 2,
        venueId: 10,
        auctionDate: '2024-03-30',
        status: 'scheduled',
        createdAt: '2024-03-01',
        updatedAt: '2024-03-01',
      } as EntityAuction;

      const result = toAuction(entity);

      expect(result.duration.label).toBe('');
      expect(result.duration.startTime).toBeNull();
    });
  });
});
