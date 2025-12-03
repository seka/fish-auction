import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query';
import { registerBuyer, getBuyers } from '@/src/api/admin';

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

export const useBuyerMutations = () => {
    const queryClient = useQueryClient();

    const createMutation = useMutation({
        mutationFn: (data: { name: string }) => registerBuyer(data.name),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['buyers'] });
        },
    });

    return {
        createBuyer: createMutation.mutateAsync,
        isCreating: createMutation.isPending,
        error: createMutation.error,
    };
};
