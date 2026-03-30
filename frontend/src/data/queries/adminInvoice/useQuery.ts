import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { getInvoices } from '@/src/data/api/invoice';
import { InvoiceItem } from '@entities/invoice';
import { adminInvoiceKeys } from './keys';

export const useInvoiceQuery = <T = InvoiceItem[]>(
  options?: Omit<UseQueryOptions<InvoiceItem[], Error, T>, 'queryKey' | 'queryFn'>,
) => {
  return useQuery({
    queryKey: adminInvoiceKeys.all,
    queryFn: getInvoices,
    ...options,
  });
};
