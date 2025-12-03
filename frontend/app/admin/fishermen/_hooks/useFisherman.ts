import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query';
import { registerFisherman, getFishermen } from '@/src/api/admin';

export const useFishermen = () => {
    const { data, error, isLoading } = useQuery({
        queryKey: ['fishermen'],
        queryFn: getFishermen,
    });

    return {
        fishermen: data ?? [],
        error,
        isLoading,
    };
};

export const useFishermanMutations = () => {
    const queryClient = useQueryClient();

    const createMutation = useMutation({
        mutationFn: (data: { name: string }) => registerFisherman(data.name),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['fishermen'] });
        },
    });

    return {
        createFisherman: createMutation.mutateAsync,
        isCreating: createMutation.isPending,
        error: createMutation.error,
    };
};
