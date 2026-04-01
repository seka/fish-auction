import { describe, it, expect } from 'vitest';
import { Venue as EntityVenue } from '@entities/venue';
import { toVenue } from './venue';

describe('admin/types/venue', () => {
  describe('toVenue', () => {
    it('should map EntityVenue to Admin Venue model correctly', () => {
      const entity: EntityVenue = {
        id: 10,
        name: 'Venue A',
        location: 'Tokyo',
        description: 'A great venue',
        createdAt: '2024-03-01T00:00:00Z',
      };

      const result = toVenue(entity);

      expect(result.id).toBe(10);
      expect(result.location).toBe('Tokyo');
      expect(result.createdAt).toBe('2024-03-01T00:00:00Z');
    });
  });
});
