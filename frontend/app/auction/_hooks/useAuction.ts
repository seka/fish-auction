import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getItems, submitBid } from '@/src/api/auction';
import { getBuyers } from '@/src/api/admin';
import { AuctionItem, Bid } from '@/src/models';

interface UseItemsOptions {
    status?: string;
    pollingInterval?: number;
}

export const useItems = ({ status, pollingInterval }: UseItemsOptions = {}) => {
    const { data, error, isLoading, refetch } = useQuery({
        queryKey: ['items', status],
        queryFn: () => getItems(status),
        refetchInterval: pollingInterval || false,
    });

    return {
        items: data || [],
        isLoading,
        error,
        refetch,
    };
};

export const useSubmitBid = () => {
    const queryClient = useQueryClient();

    const mutation = useMutation({
        mutationFn: (bid: Bid) => submitBid(bid),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['items'] });
        },
    });

    return {
        submitBid: mutation.mutateAsync,
        isLoading: mutation.isPending,
        error: mutation.error,
    };
};

export const useBuyers = () => {
    const { data, error, isLoading } = useQuery({
        queryKey: ['buyers'],
        queryFn: getBuyers,
    });

    return {
        buyers: data ?? [],
        error,
        isLoading,
    };
};
