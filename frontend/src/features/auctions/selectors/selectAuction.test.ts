import { describe, it, expect } from 'vitest';
import { Auction as EntityAuction } from '@entities/auction';
import {
  selectIsAuctionActive,
  selectTimeLabel,
  selectTime,
  selectMinimumBidIncrement,
  selectNextMinimumBid,
  selectVisiblePublicAuctions,
} from './selectAuction';

describe('auctions/selectors/selectAuction', () => {
  describe('selectTime', () => {
    it('should format HH:MM:SS to HH:MM', () => {
      expect(selectTime('10:00:00')).toBe('10:00');
      expect(selectTime('23:59:59')).toBe('23:59');
    });

    it('should return empty string for null/undefined', () => {
      expect(selectTime(null)).toBe('');
      expect(selectTime(undefined)).toBe('');
    });
  });

  describe('selectTimeLabel', () => {
    it('should format time range with ~ separator', () => {
      expect(selectTimeLabel('10:00:00', '12:00:00')).toBe('10:00 ~ 12:00');
    });

    it('should handle null/missing times with placeholders', () => {
      expect(selectTimeLabel('10:00:00', null)).toBe('10:00 ~ --:--');
      expect(selectTimeLabel(null, '12:00:00')).toBe('--:-- ~ 12:00');
      expect(selectTimeLabel(null, null)).toBe('');
    });
  });

  describe('selectMinimumBidIncrement', () => {
    it('should return 100 for price < 1000', () => {
      expect(selectMinimumBidIncrement(500)).toBe(100);
    });
    it('should return 500 for price < 10000', () => {
      expect(selectMinimumBidIncrement(5000)).toBe(500);
    });
    it('should return 1000 for price < 100000', () => {
      expect(selectMinimumBidIncrement(50000)).toBe(1000);
    });
    it('should return 5000 for price >= 100000', () => {
      expect(selectMinimumBidIncrement(150000)).toBe(5000);
    });
  });

  describe('selectNextMinimumBid', () => {
    it('should calculate next minimum bid correctly', () => {
      expect(selectNextMinimumBid(500)).toBe(600); // 500 + 100
      expect(selectNextMinimumBid(5000)).toBe(5500); // 5000 + 500
      expect(selectNextMinimumBid(50000)).toBe(51000); // 50000 + 1000
    });
  });

  describe('selectIsAuctionActive', () => {
    it('should return true if status is in_progress and current time is within auction hours', () => {
      const auction: EntityAuction = {
        auctionDate: '2024-03-30',
        startTime: '10:00:00',
        endTime: '12:00:00',
        status: 'in_progress',
      } as EntityAuction;
      const now = new Date('2024-03-30T11:00:00+09:00');
      expect(selectIsAuctionActive(auction, now)).toBe(true);
    });

    it('should return false if status is scheduled regardless of time', () => {
      const auction: EntityAuction = {
        auctionDate: '2024-03-30',
        startTime: '10:00:00',
        endTime: '12:00:00',
        status: 'scheduled',
      } as EntityAuction;
      const now = new Date('2024-03-30T11:00:00+09:00');
      expect(selectIsAuctionActive(auction, now)).toBe(false);
    });

    it('should return false if status is completed or cancelled', () => {
      expect(selectIsAuctionActive({ status: 'completed' } as EntityAuction)).toBe(false);
      expect(selectIsAuctionActive({ status: 'cancelled' } as EntityAuction)).toBe(false);
    });

    it('should return false if status is in_progress but current time is before auction start', () => {
      const auction = {
        auctionDate: '2024-03-30',
        startTime: '10:00:00',
        endTime: '12:00:00',
        status: 'in_progress',
      } as EntityAuction;
      const now = new Date('2024-03-30T09:59:59+09:00');
      expect(selectIsAuctionActive(auction, now)).toBe(false);
    });

    it('should return false if status is in_progress but current time is after auction end', () => {
      const auction = {
        auctionDate: '2024-03-30',
        startTime: '10:00:00',
        endTime: '12:00:00',
        status: 'in_progress',
      } as EntityAuction;
      const now = new Date('2024-03-30T12:00:01+09:00');
      expect(selectIsAuctionActive(auction, now)).toBe(false);
    });

    it('should return false if status is in_progress but startTime or endTime is not set', () => {
      const auction = {
        auctionDate: '2024-03-30',
        startTime: null,
        endTime: null,
        status: 'in_progress',
      } as EntityAuction;
      expect(selectIsAuctionActive(auction)).toBe(false);
    });
  });

  describe('selectVisiblePublicAuctions', () => {
    it('should prioritize active scheduled auctions over inactive ones, even if later in date', () => {
      const now = new Date('2024-03-30T12:00:00+09:00');
      const auctions = [
        {
          id: 1,
          status: 'scheduled',
          auctionDate: '2024-03-31',
          startTime: '10:00:00',
          endTime: '12:00:00',
        },
        {
          id: 2,
          status: 'in_progress',
          auctionDate: '2024-03-30',
          startTime: '11:00:00', // Active (11:00-13:00) vs now (12:00)
          endTime: '13:00:00',
        },
        {
          id: 3,
          status: 'scheduled',
          auctionDate: '2024-03-29', // Past (inactive)
          startTime: '10:00:00',
          endTime: '12:00:00',
        },
      ] as EntityAuction[];

      const visible = selectVisiblePublicAuctions(auctions, now);

      expect(visible[0].id).toBe(2); // Active scheduled first
      expect(visible[1].id).toBe(3); // Inactive but earlier date
      expect(visible[2].id).toBe(1); // Inactive and later date
    });
  });
});
