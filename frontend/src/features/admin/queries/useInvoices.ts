import { InvoiceItem as EntityInvoiceItem } from '@entities/invoice';
import { useInvoiceQuery } from '@/src/data/queries/adminInvoice/useQuery';
import { InvoiceItem } from '../types/invoice';

export const toInvoice = (entity: EntityInvoiceItem): InvoiceItem => ({
  buyerId: entity.buyerId,
  buyerName: entity.buyerName,
  totalAmount: entity.totalAmount,
});

/**
 * 管理画面用の請求書一覧クエリフック
 */
export const useInvoices = () => {
  const { data: invoices, ...rest } = useInvoiceQuery({
    select: (data) => data.map(toInvoice),
  });

  return {
    ...rest,
    invoices: invoices || [],
  };
};
