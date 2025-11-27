import { useQuery } from '@tanstack/react-query';

export interface InvoiceItem {
    buyer_id: number;
    buyer_name: string;
    total_amount: number;
}

export function useInvoices() {
    const { data: invoices = [], isLoading, error, refetch } = useQuery({
        queryKey: ['invoices'],
        queryFn: async () => {
            const res = await fetch('/api/invoices');
            if (!res.ok) {
                throw new Error('Failed to fetch invoices');
            }
            return res.json() as Promise<InvoiceItem[]>;
        },
    });

    return {
        invoices,
        isLoading,
        error: error ? (error as Error).message : null,
        refetch,
    };
}
