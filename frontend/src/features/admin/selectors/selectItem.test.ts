import { describe, it, expect } from 'vitest';
import { selectNextMinimumBid } from './selectItem';

describe('features/admin/selectors/selectItem', () => {
  describe('selectNextMinimumBid', () => {
    it('should return 100 for price < 1000', () => {
      expect(selectNextMinimumBid(500)).toBe(600);
    });

    it('should return 500 for price < 10000', () => {
      expect(selectNextMinimumBid(5000)).toBe(5500);
    });

    it('should return 1000 for price < 100000', () => {
      expect(selectNextMinimumBid(50000)).toBe(51000);
    });

    it('should return 5000 for price >= 100000', () => {
      expect(selectNextMinimumBid(150000)).toBe(155000);
    });
  });
});
