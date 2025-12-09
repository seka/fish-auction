import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query';
import { registerBuyer, getBuyers } from '@/src/api/admin';

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
        mutationFn: (data: { name: string }) => registerBuyer(data.name),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: buyerKeys.all });
        },
    });

    return {
        createBuyer: createMutation.mutateAsync,
        isCreating: createMutation.isPending,
        error: createMutation.error,
    };
};
