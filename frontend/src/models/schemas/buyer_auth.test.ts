import { describe, it, expect } from 'vitest';
import { getBuyerLoginSchema } from './buyer_auth';
import { tMock } from './test-utils';

const t = tMock;
const buyerLoginSchema = getBuyerLoginSchema(t);

describe('Buyer Auth Schemas', () => {
  describe('buyerLoginSchema', () => {
    it('should accept valid login credentials', () => {
      const result = buyerLoginSchema.safeParse({
        email: 'buyer@example.com',
        password: 'password123',
      });
      expect(result.success).toBe(true);
    });

    it('should reject invalid email format', () => {
      const result = buyerLoginSchema.safeParse({
        email: 'not-an-email',
        password: 'password123',
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('有効なメールアドレスを入力してください');
      }
    });

    it('should reject empty password', () => {
      const result = buyerLoginSchema.safeParse({
        email: 'buyer@example.com',
        password: '',
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('パスワードを入力してください');
      }
    });
  });
});
