import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query';
import { registerBuyer, getBuyers, deleteBuyer } from '@/src/api/admin';

export const buyerKeys = {
    all: ['buyers'] as const,
};

export const useBuyerQuery = () => {
    const { data, error, isLoading } = useQuery({
        queryKey: buyerKeys.all,
        queryFn: getBuyers,
    });

    return {
        buyers: data ?? [],
        error,
        isLoading,
    };
};

export const useBuyerMutation = () => {
    const queryClient = useQueryClient();

    const createMutation = useMutation({
        mutationFn: (data: { name: string; email: string; password: string; organization: string; contactInfo: string }) => registerBuyer(data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: buyerKeys.all });
        },
    });

    const deleteMutation = useMutation({
        mutationFn: deleteBuyer,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: buyerKeys.all });
        },
    });

    return {
        createBuyer: createMutation.mutateAsync,
        isCreating: createMutation.isPending,
        deleteBuyer: deleteMutation.mutateAsync,
        isDeleting: deleteMutation.isPending,
        error: createMutation.error || deleteMutation.error,
    };
};
