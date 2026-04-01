import { describe, it, expect } from 'vitest';
import { InvoiceItem as EntityInvoiceItem } from '@entities/invoice';
import { toInvoice } from './invoice';

describe('admin/types/invoice', () => {
  describe('toInvoice', () => {
    it('should map EntityInvoiceItem to Admin InvoiceItem model correctly', () => {
      const entity: EntityInvoiceItem = {
        buyerId: 101,
        buyerName: 'Buyer A',
        totalAmount: 15000,
      };

      const result = toInvoice(entity);

      expect(result.buyerId).toBe(101);
      expect(result.buyerName).toBe('Buyer A');
      expect(result.totalAmount).toBe(15000);
    });
  });
});
