'use client';
import { useInvoiceQuery } from '@/src/data/queries/adminInvoice/useQuery';
import { toInvoice } from '../types/invoice';

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
