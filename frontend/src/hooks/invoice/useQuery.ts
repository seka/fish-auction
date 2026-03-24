import { useQuery } from '@tanstack/react-query';
import { getInvoices } from '@/src/api/invoice';
import { invoiceKeys } from './keys';

export const useInvoiceQuery = () => {
  const {
    data: invoices,
    isLoading,
    error,
  } = useQuery({
    queryKey: invoiceKeys.all,
    queryFn: getInvoices,
  });

  return { invoices: invoices || [], isLoading, error };
};
