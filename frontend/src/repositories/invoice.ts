import { useQuery } from '@tanstack/react-query';
import { getInvoices } from '@/src/api/invoice';

export const invoiceKeys = {
    all: ['invoices'] as const,
};

export const useInvoiceQuery = () => {
    const { data: invoices, error, isLoading, refetch } = useQuery({
        queryKey: invoiceKeys.all,
        queryFn: () => getInvoices(),
    });

    return {
        invoices: invoices || [],
        isLoading,
        error,
        refetch,
    };
};
