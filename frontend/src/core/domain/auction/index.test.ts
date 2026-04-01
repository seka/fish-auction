import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { isAuctionActive } from './index';

describe('isAuctionActive', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('should return false for completed or cancelled status regardless of time', () => {
    const auctionDate = '2024-03-30';
    const startTime = '10:00';
    const endTime = '12:00';
    
    // Set system time to within range
    vi.setSystemTime(new Date('2024-03-30T11:00:00'));

    expect(isAuctionActive({ status: 'completed', auctionDate, startTime, endTime })).toBe(false);
    expect(isAuctionActive({ status: 'cancelled', auctionDate, startTime, endTime })).toBe(false);
  });

  it('should return true for in_progress status regardless of time', () => {
    const auctionDate = '2024-03-30';
    const startTime = '10:00';
    const endTime = '12:00';
    
    // Set system time outside range
    vi.setSystemTime(new Date('2024-03-30T09:00:00'));

    expect(isAuctionActive({ status: 'in_progress', auctionDate, startTime, endTime })).toBe(true);
  });

  describe('scheduled status with time validation', () => {
    const auctionDate = '2024-03-30';
    const startTime = '10:00';
    const endTime = '12:00';

    it('should return true when current time is between startTime and endTime', () => {
      vi.setSystemTime(new Date('2024-03-30T11:00:00'));
      expect(isAuctionActive({ status: 'scheduled', auctionDate, startTime, endTime })).toBe(true);
    });

    it('should return false when current time is before startTime', () => {
      vi.setSystemTime(new Date('2024-03-30T09:59:59'));
      expect(isAuctionActive({ status: 'scheduled', auctionDate, startTime, endTime })).toBe(false);
    });

    it('should return false when current time is after endTime', () => {
      vi.setSystemTime(new Date('2024-03-30T12:00:01'));
      expect(isAuctionActive({ status: 'scheduled', auctionDate, startTime, endTime })).toBe(false);
    });
  });

  describe('edge cases and null safety', () => {
    const auctionDate = '2024-03-30';

    it('should return true for scheduled status when startTime is null', () => {
      // ロジック: 時刻がない場合はステータスに従う (scheduled なら true/判定可能とする)
      // 注意: 現在の実装は if (!auction.startTime || !auction.endTime) return auction.status === 'scheduled';
      expect(isAuctionActive({ status: 'scheduled', auctionDate, startTime: null, endTime: '12:00' })).toBe(true);
    });

    it('should return false for other status when startTime is null', () => {
      expect(isAuctionActive({ status: 'completed', auctionDate, startTime: null, endTime: '12:00' })).toBe(false);
    });

    it('should return false when parsing fails', () => {
      // 不正な時刻形式
      expect(isAuctionActive({ status: 'scheduled', auctionDate, startTime: 'invalid', endTime: 'invalid' })).toBe(false);
    });
  });
});
