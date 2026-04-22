import { describe, it, expect } from 'vitest';
import { AuctionStatus as EntityAuctionStatus } from '@entities/auction';
import { selectAuctionStatus, selectTimeLabel } from './selectAuction';

describe('features/admin/selectors/selectAuction', () => {
  describe('selectAuctionStatus', () => {
    it('should return correct status object for each status', () => {
      expect(selectAuctionStatus('scheduled')).toEqual({
        value: 'scheduled',
        labelKey: 'scheduled',
        variant: 'info',
        isScheduled: true,
        isInProgress: false,
        isCompleted: false,
        isCancelled: false,
      });
      expect(selectAuctionStatus('in_progress')).toEqual({
        value: 'in_progress',
        labelKey: 'in_progress',
        variant: 'success',
        isScheduled: false,
        isInProgress: true,
        isCompleted: false,
        isCancelled: false,
      });
      expect(selectAuctionStatus('completed')).toEqual({
        value: 'completed',
        labelKey: 'completed',
        variant: 'neutral',
        isScheduled: false,
        isInProgress: false,
        isCompleted: true,
        isCancelled: false,
      });
      expect(selectAuctionStatus('cancelled')).toEqual({
        value: 'cancelled',
        labelKey: 'cancelled',
        variant: 'error',
        isScheduled: false,
        isInProgress: false,
        isCompleted: false,
        isCancelled: true,
      });
    });

    it('should fallback to scheduled for unknown status', () => {
      expect(selectAuctionStatus('unknown' as unknown as EntityAuctionStatus)).toEqual({
        value: 'scheduled',
        labelKey: 'scheduled',
        variant: 'info',
        isScheduled: true,
        isInProgress: false,
        isCompleted: false,
        isCancelled: false,
      });
    });
  });

  describe('selectTimeLabel', () => {
    it('should format time range with ~ separator', () => {
      const start = new Date('2024-03-30T01:00:00Z'); // JST 10:00
      const end = new Date('2024-03-30T03:00:00Z');   // JST 12:00
      expect(selectTimeLabel(start, end)).toBe('10:00 ~ 12:00');
    });

    it('should handle null/missing times with placeholders', () => {
      const start = new Date('2024-03-30T01:00:00Z'); // JST 10:00
      const end = new Date('2024-03-30T03:00:00Z');   // JST 12:00
      expect(selectTimeLabel(start, null)).toBe('10:00 ~ --:--');
      expect(selectTimeLabel(null, end)).toBe('--:-- ~ 12:00');
      expect(selectTimeLabel(null, null)).toBe('');
    });
  });
});
