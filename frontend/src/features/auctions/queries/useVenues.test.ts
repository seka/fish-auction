import { describe, it, expect } from 'vitest';
import { toVenue } from './useVenues';
import { Venue as EntityVenue } from '@entities/venue';

describe('auctions/queries/useVenues mapping', () => {
  it('should map EntityVenue to Auctions Venue model correctly', () => {
    const entity: EntityVenue = {
      id: 10,
      name: 'Test Venue',
      location: 'Chiba',
      description: 'A test venue',
      createdAt: '2024-03-01T00:00:00Z',
    };

    const result = toVenue(entity);

    expect(result.id).toBe(10);
    expect(result.name).toBe('Test Venue');
    expect(result.location).toBe('Chiba');
  });
});
