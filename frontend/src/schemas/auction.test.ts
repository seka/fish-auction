import { describe, it, expect } from 'vitest';
import { getVenueSchema, getAuctionSchema, getBidSchema } from './auction';
import { tMock } from '@testings/i18n';

const t = tMock;
const venueSchema = getVenueSchema(t);
const auctionSchema = getAuctionSchema(t);
const bidSchema = getBidSchema(t);

describe('Auction Schemas', () => {
  describe('venueSchema', () => {
    it('should accept valid venue data', () => {
      const result = venueSchema.safeParse({ name: '築地市場' });
      expect(result.success).toBe(true);
    });

    it('should reject empty venue name', () => {
      const result = venueSchema.safeParse({ name: '' });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toBe('key:required(field:key:Admin.Venues.name)');
      }
    });
  });

  describe('auctionSchema', () => {
    it('should accept valid auction data', () => {
      const data = {
        venueId: '1',
        startAt: '2026-03-30',
        status: 'scheduled',
      };
      const result = auctionSchema.safeParse(data);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.venueId).toBe(1);
      }
    });

    it('should reject invalid venueId', () => {
      const result = auctionSchema.safeParse({
        venueId: '0',
        auctionDate: '2026-03-30',
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toBe(
          'key:select_required(field:key:Admin.Auctions.venue)',
        );
      }
    });

    it('should reject empty auctionDate', () => {
      const result = auctionSchema.safeParse({
        venueId: '1',
        startAt: '',
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toBe('key:required(field:key:Admin.Auctions.start_time)');
      }
    });
  });

  describe('bidSchema', () => {
    it('should accept positive price strings', () => {
      const result = bidSchema.safeParse({ price: '1500' });
      expect(result.success).toBe(true);
    });

    it('should reject empty price', () => {
      const result = bidSchema.safeParse({ price: '' });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toBe('key:required(field:key:field_name.price)');
      }
    });

    it('should reject non-numeric price', () => {
      const result = bidSchema.safeParse({ price: 'abc' });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toBe('key:positive_number');
      }
    });

    it('should reject zero or negative price', () => {
      const result = bidSchema.safeParse({ price: '0' });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toBe('key:positive_number');
      }
    });
  });
});
