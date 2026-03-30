import { describe, it, expect } from 'vitest';
import { toBuyer } from './useBuyers';
import { Buyer as EntityBuyer } from '@entities/admin';

describe('admin/queries/useBuyers mapping', () => {
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
