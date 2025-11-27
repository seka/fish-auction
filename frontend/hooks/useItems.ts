import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

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
    const { data: items = [], isLoading, error, refetch } = useQuery({
        queryKey: ['items', options.status],
        queryFn: async () => {
            const url = options.status
                ? `/api/items?status=${options.status}`
                : '/api/items';
            const res = await fetch(url);
            if (!res.ok) {
                throw new Error('Failed to fetch items');
            }
            return res.json() as Promise<AuctionItem[]>;
        },
        refetchInterval: options.pollingInterval,
    });

    return {
        items,
        isLoading,
        error: error ? (error as Error).message : null,
        refetch,
    };
}

interface RegisterItemParams {
    fisherman_id: number;
    fish_type: string;
    quantity: number;
    unit: string;
}

export function useRegisterItem() {
    const queryClient = useQueryClient();

    const { mutateAsync: registerItem, isPending: isLoading, error } = useMutation({
        mutationFn: async (params: RegisterItemParams) => {
            const res = await fetch('/api/items', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(params),
            });

            if (!res.ok) {
                throw new Error('Failed to register item');
            }
            return true;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['items'] });
        },
    });

    return {
        registerItem: async (params: RegisterItemParams) => {
            try {
                return await registerItem(params);
            } catch (e) {
                return false;
            }
        },
        isLoading,
        error: error ? (error as Error).message : null,
    };
}
