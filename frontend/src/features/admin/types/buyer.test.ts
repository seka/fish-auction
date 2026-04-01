import { describe, it, expect } from 'vitest';
import { Buyer as EntityBuyer } from '@entities/admin';
import { toBuyer } from './buyer';

describe('admin/types/buyer', () => {
  describe('toBuyer', () => {
    it('should map EntityBuyer to Admin Buyer model correctly', () => {
      const entity: EntityBuyer = {
        id: 1,
        name: 'Buyer A',
      };

      const result = toBuyer(entity);

      expect(result.id).toBe(1);
      expect(result.name).toBe('Buyer A');
    });
  });
});
