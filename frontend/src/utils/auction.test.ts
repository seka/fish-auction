import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { isAuctionActive, formatTime } from './auction';
import { Auction } from '@/src/models/auction';

describe('Auction Utils', () => {
    describe('isAuctionActive', () => {
        beforeEach(() => {
            // Mock system time to 2023-01-01 12:00:00
            vi.useFakeTimers();
            vi.setSystemTime(new Date('2023-01-01T12:00:00'));
        });

        afterEach(() => {
            vi.useRealTimers();
        });

        const baseAuction: Auction = {
            id: 1,
            auctionDate: '2023-01-01',
            startTime: '10:00:00',
            endTime: '15:00:00',
            venueId: 1,
            status: 'scheduled',
            createdAt: '',
            updatedAt: ''
        };

        it('returns true when current time is within auction range', () => {
            expect(isAuctionActive(baseAuction)).toBe(true);
        });

        it('returns false when current time is before auction start', () => {
            vi.setSystemTime(new Date('2023-01-01T09:00:00'));
            expect(isAuctionActive(baseAuction)).toBe(false);
        });

        it('returns false when current time is after auction end', () => {
            vi.setSystemTime(new Date('2023-01-01T16:00:00'));
            expect(isAuctionActive(baseAuction)).toBe(false);
        });

        it('returns true if time is not set', () => {
            const noTimeAuction = { ...baseAuction, startTime: '', endTime: '' };
            expect(isAuctionActive(noTimeAuction)).toBe(true);
        });
    });

    describe('formatTime', () => {
        it('formats HH:MM:SS to HH:MM', () => {
            expect(formatTime('12:34:56')).toBe('12:34');
        });

        it('returns empty string for undefined', () => {
            expect(formatTime(undefined)).toBe('');
        });
    });
});
