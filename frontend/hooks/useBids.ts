import { useMutation, useQueryClient } from '@tanstack/react-query';

interface SubmitBidParams {
    item_id: number;
    buyer_id: number;
    price: number;
}

export function useSubmitBid() {
    const queryClient = useQueryClient();

    const { mutateAsync: submitBid, isPending: isLoading, error } = useMutation({
        mutationFn: async (params: SubmitBidParams) => {
            const res = await fetch('/api/bid', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(params),
            });

            if (!res.ok) {
                throw new Error('Failed to submit bid');
            }
            return true;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['items'] });
        },
    });

    return {
        submitBid: async (params: SubmitBidParams) => {
            try {
                return await submitBid(params);
            } catch (e) {
                return false;
            }
        },
        isLoading,
        error: error ? (error as Error).message : null,
    };
}
