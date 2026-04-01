import { describe, it, expect } from 'vitest';
import { Auction as EntityAuction } from '@entities/auction';
import { toAuction } from './auction';

describe('admin/types/auction toAuction', () => {
  it('should map EntityAuction to Admin Auction model correctly', () => {
    const entity: EntityAuction = {
      id: 1,
      venueId: 10,
      auctionDate: '2024-03-30',
      startTime: '10:00:00',
      endTime: '12:00:00',
      status: 'in_progress',
      createdAt: '2024-03-01T00:00:00Z',
      updatedAt: '2024-03-01T00:00:00Z',
    };

    const result = toAuction(entity);

    expect(result.id).toBe(1);
    expect(result.createdAt).toBe('2024-03-01T00:00:00Z');
    expect(result.startTime).toBe('10:00:00');
  });

  it('should handle null startTime in admin mapping', () => {
    const entity: EntityAuction = {
      id: 1,
      venueId: 10,
      auctionDate: '2024-03-30',
      status: 'scheduled',
      createdAt: '2024-03-01T00:00:00Z',
      updatedAt: '2024-03-01T00:00:00Z',
    } as EntityAuction;

    const result = toAuction(entity);

    expect(result.startTime).toBeNull();
    expect(result.endTime).toBeNull();
  });
});
