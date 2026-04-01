import { InvoiceItem as EntityInvoiceItem } from '@entities/invoice';

export interface InvoiceItem {
  buyerId: number;
  buyerName: string;
  totalAmount: number;
}

export const toInvoice = (entity: EntityInvoiceItem): InvoiceItem => ({
  buyerId: entity.buyerId,
  buyerName: entity.buyerName,
  totalAmount: entity.totalAmount,
});
