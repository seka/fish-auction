import { describe, it, expect } from 'vitest';
import { ItemStatus as EntityItemStatus } from '@entities/auction';
import { selectItemStatus } from './selectItem';

describe('features/auctions/selectors/selectItem', () => {
  describe('selectItemStatus', () => {
    it('should return correct status object for each status', () => {
      expect(selectItemStatus('Pending')).toEqual({
        value: 'Pending',
        labelKey: 'Pending',
        variant: 'info',
        isPending: true,
        isBidding: false,
        isSold: false,
        isUnsold: false,
      });
      expect(selectItemStatus('Bidding')).toEqual({
        value: 'Bidding',
        labelKey: 'Bidding',
        variant: 'success',
        isPending: false,
        isBidding: true,
        isSold: false,
        isUnsold: false,
      });
      expect(selectItemStatus('Sold')).toEqual({
        value: 'Sold',
        labelKey: 'Sold',
        variant: 'neutral',
        isPending: false,
        isBidding: false,
        isSold: true,
        isUnsold: false,
      });
      expect(selectItemStatus('Unsold')).toEqual({
        value: 'Unsold',
        labelKey: 'Unsold',
        variant: 'error',
        isPending: false,
        isBidding: false,
        isSold: false,
        isUnsold: true,
      });
    });

    it('should fallback to pending for unknown status', () => {
      expect(selectItemStatus('Unknown' as unknown as EntityItemStatus)).toEqual({
        value: 'Pending',
        labelKey: 'Pending',
        variant: 'info',
        isPending: true,
        isBidding: false,
        isSold: false,
        isUnsold: false,
      });
    });
  });
});
