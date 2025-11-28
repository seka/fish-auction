import { useQuery } from '@tanstack/react-query';
import { getInvoices } from '@/src/api/invoice';

export const useInvoices = () => {
    const { data, error, isLoading, refetch } = useQuery({
        queryKey: ['invoices'],
        queryFn: () => getInvoices(),
    });

    return {
        invoices: data || [],
        isLoading,
        error,
        refetch,
    };
};
