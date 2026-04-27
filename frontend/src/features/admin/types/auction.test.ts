import { describe, it, expect } from 'vitest';
import { Auction as EntityAuction } from '@entities/auction';
import { toAuction } from './auction';

describe('admin/types/auction toAuction', () => {
  it('should map EntityAuction to Admin Auction model correctly', () => {
    const entity: EntityAuction = {
      id: 1,
      venueId: 10,
      startAt: '2024-03-30T10:00:00+09:00',
      endAt: '2024-03-30T12:00:00+09:00',
      status: 'in_progress',
      createdAt: '2024-03-01T00:00:00Z',
      updatedAt: '2024-03-01T00:00:00Z',
    };

    const result = toAuction(entity);

    expect(result.id).toBe(1);
    expect(result.status).toEqual({
      value: 'in_progress',
      labelKey: 'in_progress',
      variant: 'success',
      isScheduled: false,
      isInProgress: true,
      isCompleted: false,
      isCancelled: false,
    });
    expect(result.duration.startAt).toBeInstanceOf(Date);
    expect(result.duration.startAt!.toISOString()).toBe('2024-03-30T01:00:00.000Z');
    expect(result.duration.label).toBe('10:00 ~ 12:00');
    expect(result.actions).toEqual({
      canStart: false,
      canFinish: true,
    });
    expect(result.createdAt).toBe('2024-03-01T00:00:00Z');
  });

  it('should handle null startAt in admin mapping', () => {
    const entity: EntityAuction = {
      id: 1,
      venueId: 10,
      status: 'scheduled',
      createdAt: '2024-03-01T00:00:00Z',
      updatedAt: '2024-03-01T00:00:00Z',
    } as EntityAuction;

    const result = toAuction(entity);

    expect(result.duration.label).toBe('');
    expect(result.actions.canStart).toBe(true);
  });
});
