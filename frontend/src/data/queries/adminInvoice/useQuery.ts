import { useQuery } from '@tanstack/react-query';
import { getInvoices } from '@/src/data/api/invoice';
import { adminInvoiceKeys } from './keys';

export const useInvoiceQuery = () => {
  const {
    data: invoices,
    isLoading,
    error,
  } = useQuery({
    queryKey: adminInvoiceKeys.all,
    queryFn: getInvoices,
  });

  return { invoices: invoices || [], isLoading, error };
};
