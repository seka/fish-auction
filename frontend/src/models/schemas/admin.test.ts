import { describe, it, expect } from 'vitest';
import { getFishermanSchema, getBuyerSchema, getItemSchema } from './admin';
import { tMock } from './test-utils';

const t = tMock;
const fishermanSchema = getFishermanSchema(t);
const buyerSchema = getBuyerSchema(t);
const itemSchema = getItemSchema(t);

describe('Admin Schemas', () => {
  describe('fishermanSchema', () => {
    it('should accept valid fisherman data', () => {
      const result = fishermanSchema.safeParse({ name: '山田 太郎' });
      expect(result.success).toBe(true);
    });

    it('should reject empty name', () => {
      const result = fishermanSchema.safeParse({ name: '' });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('漁師名を入力してください');
      }
    });
  });

  describe('buyerSchema', () => {
    it('should accept valid buyer data', () => {
      const result = buyerSchema.safeParse({
        name: '株式会社 魚市場',
        email: 'buyer@example.com',
        password: 'Password123!',
        organization: '築地組合',
        contactInfo: '03-1234-5678',
      });
      expect(result.success).toBe(true);
    });

    it('should reject empty name', () => {
      const result = buyerSchema.safeParse({
        name: '',
        email: 'buyer@example.com',
        password: 'Password123!',
        organization: '組合',
        contactInfo: '03-0000-0000',
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        const nameIssue = result.error.issues.find((i) => i.path[0] === 'name');
        expect(nameIssue?.message).toContain('中買人名を入力してください');
      }
    });
  });

  describe('itemSchema', () => {
    it('should accept valid item data', () => {
      const data = {
        auctionId: '1',
        fishermanId: '10',
        fishType: 'マグロ',
        quantity: '5',
        unit: '匹',
      };
      const result = itemSchema.safeParse(data);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(typeof result.data.auctionId).toBe('string');
        expect(result.data.auctionId).toBe('1');
      }
    });

    it('should reject invalid auctionId', () => {
      const result = itemSchema.safeParse({
        auctionId: '',
        fishermanId: '1',
        fishType: 'マグロ',
        quantity: '5',
        unit: '匹',
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('セリを選択してください');
      }
    });

    it('should reject empty fishType', () => {
      const result = itemSchema.safeParse({
        auctionId: '1',
        fishermanId: '1',
        fishType: '',
        quantity: '5',
        unit: '匹',
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('魚種を入力してください');
      }
    });
  });
});
