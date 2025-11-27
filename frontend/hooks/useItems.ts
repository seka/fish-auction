import { useState, useEffect, useCallback } from 'react';

export interface AuctionItem {
    id: number;
    fisherman_id: number;
    fish_type: string;
    quantity: number;
    unit: string;
    status: string;
    created_at: string;
}

interface UseItemsOptions {
    status?: string;
    pollingInterval?: number; // in milliseconds, 0 to disable
}

export function useItems(options: UseItemsOptions = {}) {
    const [items, setItems] = useState<AuctionItem[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchItems = useCallback(async () => {
        try {
            const url = options.status
                ? `/api/items?status=${options.status}`
                : '/api/items';

            const res = await fetch(url);
            if (res.ok) {
                const data = await res.json();
                setItems(data || []);
                setError(null);
            } else {
                setError('Failed to fetch items');
            }
        } catch (err) {
            setError('Error fetching items');
            console.error('Failed to fetch items', err);
        } finally {
            setIsLoading(false);
        }
    }, [options.status]);

    useEffect(() => {
        fetchItems();

        if (options.pollingInterval && options.pollingInterval > 0) {
            const interval = setInterval(fetchItems, options.pollingInterval);
            return () => clearInterval(interval);
        }
    }, [fetchItems, options.pollingInterval]);

    return {
        items,
        isLoading,
        error,
        refetch: fetchItems,
    };
}

interface RegisterItemParams {
    fisherman_id: number;
    fish_type: string;
    quantity: number;
    unit: string;
}

export function useRegisterItem() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const registerItem = async (params: RegisterItemParams): Promise<boolean> => {
        setIsLoading(true);
        setError(null);

        try {
            const res = await fetch('/api/items', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(params),
            });

            if (res.ok) {
                setIsLoading(false);
                return true;
            } else {
                setError('Failed to register item');
                setIsLoading(false);
                return false;
            }
        } catch (err) {
            setError('Error registering item');
            setIsLoading(false);
            return false;
        }
    };

    return {
        registerItem,
        isLoading,
        error,
    };
}
