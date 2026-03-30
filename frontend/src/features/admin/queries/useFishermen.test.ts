import { describe, it, expect } from 'vitest';
import { toFisherman } from './useFishermen';
import { Fisherman as EntityFisherman } from '@entities/admin';

describe('admin/queries/useFishermen mapping', () => {
  it('should map EntityFisherman to Admin Fisherman model correctly', () => {
    const entity: EntityFisherman = {
      id: 5,
      name: 'Fisherman A',
    };

    const result = toFisherman(entity);

    expect(result.id).toBe(5);
    expect(result.name).toBe('Fisherman A');
  });
});
