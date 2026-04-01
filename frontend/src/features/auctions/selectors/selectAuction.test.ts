import { describe, it, expect } from 'vitest';
import { Auction as EntityAuction } from '@entities/auction';
import { selectIsAuctionActive, selectTime, selectMinimumBidIncrement } from './selectAuction';

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

  describe('selectIsAuctionActive', () => {
    it('should return true if current time is within auction hours', () => {
      const auctionDate = '2024-03-30';
      const auction: EntityAuction = {
        auctionDate,
        startTime: '10:00:00',
        endTime: '12:00:00',
      } as EntityAuction;
      
      const [year, month, day] = auctionDate.split('-').map(Number);
      const now = new Date(year, month - 1, day, 11, 0, 0);
      
      expect(selectIsAuctionActive(auction, now)).toBe(true);
    });

    it('should return false if current time is before auction starts', () => {
      const auctionDate = '2024-03-30';
      const auction: EntityAuction = {
        auctionDate,
        startTime: '10:00:00',
        endTime: '12:00:00',
      } as EntityAuction;
      
      const [year, month, day] = auctionDate.split('-').map(Number);
      const now = new Date(year, month - 1, day, 9, 0, 0);
      
      expect(selectIsAuctionActive(auction, now)).toBe(false);
    });

    it('should return false if current time is after auction ends', () => {
      const auctionDate = '2024-03-30';
      const auction: EntityAuction = {
        auctionDate,
        startTime: '10:00:00',
        endTime: '12:00:00',
      } as EntityAuction;
      
      const [year, month, day] = auctionDate.split('-').map(Number);
      const now = new Date(year, month - 1, day, 13, 0, 0);
      
      expect(selectIsAuctionActive(auction, now)).toBe(false);
    });

    it('should return false if startTime or endTime is missing', () => {
      const auction = {} as EntityAuction;
      expect(selectIsAuctionActive(auction)).toBe(false);
    });
  });
});
