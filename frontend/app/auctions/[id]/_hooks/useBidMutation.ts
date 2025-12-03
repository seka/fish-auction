import { useMutation, useQueryClient } from '@tanstack/react-query';
import { submitBid } from '@/src/api/bid';

export const useBidMutation = () => {
    const queryClient = useQueryClient();
    const mutation = useMutation({
        mutationFn: submitBid,
        onSuccess: () => {
            // Invalidate items to update status/price if needed
            // But actually submitBid returns boolean, and we refetch items manually or via interval
            // Ideally we should invalidate query keys
        },
    });
    return { submitBid: mutation.mutateAsync, isLoading: mutation.isPending };
};
