import { useQuery } from '@tanstack/react-query';
import { getMyInvoices } from '@/src/data/api/invoice';
import { invoiceKeys } from './keys';

export const useMyInvoiceQuery = () => {
  const {
    data: invoices,
    isLoading,
    error,
  } = useQuery({
    queryKey: invoiceKeys.meAll(),
    queryFn: getMyInvoices,
  });

  return { invoices: invoices || [], isLoading, error };
};
