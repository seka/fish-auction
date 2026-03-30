import { describe, it, expect } from 'vitest';
import { toInvoice } from './useInvoices';
import { InvoiceItem as EntityInvoiceItem } from '@entities/invoice';

describe('admin/queries/useInvoices mapping', () => {
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
