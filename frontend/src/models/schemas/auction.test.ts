import { describe, it, expect } from 'vitest';
import { getVenueSchema, getAuctionSchema, getBidSchema } from './auction';
import { tMock } from './test-utils';

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
        expect(result.error.issues[0].message).toContain('会場名を入力してください');
      }
    });
  });

  describe('auctionSchema', () => {
    it('should accept valid auction data', () => {
      const data = {
        venueId: '1',
        auctionDate: '2026-03-30',
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
        expect(result.error.issues[0].message).toContain('会場を選択してください');
      }
    });

    it('should reject empty auctionDate', () => {
      const result = auctionSchema.safeParse({
        venueId: '1',
        auctionDate: '',
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('開催日を入力してください');
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
        // 先ほどのバグ修正が効いているか確認 (t('field_name.price') -> '価格')
        expect(result.error.issues[0].message).toContain('価格を入力してください');
      }
    });

    it('should reject non-numeric price', () => {
      const result = bidSchema.safeParse({ price: 'abc' });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toBe('正の数値を入力してください');
      }
    });

    it('should reject zero or negative price', () => {
      const result = bidSchema.safeParse({ price: '0' });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toBe('正の数値を入力してください');
      }
    });
  });
});
