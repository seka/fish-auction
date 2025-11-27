import { useState, useEffect, useCallback } from 'react';

export interface InvoiceItem {
    buyer_id: number;
    buyer_name: string;
    total_amount: number;
}

export function useInvoices() {
    const [invoices, setInvoices] = useState<InvoiceItem[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchInvoices = useCallback(async () => {
        try {
            const res = await fetch('/api/invoices');
            if (res.ok) {
                const data = await res.json();
                setInvoices(data || []);
                setError(null);
            } else {
                setError('Failed to fetch invoices');
            }
        } catch (err) {
            setError('Error fetching invoices');
            console.error('Failed to fetch invoices', err);
        } finally {
            setIsLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchInvoices();
    }, [fetchInvoices]);

    return {
        invoices,
        isLoading,
        error,
        refetch: fetchInvoices,
    };
}
