import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query';
import { registerFisherman, getFishermen, deleteFisherman } from '@/src/api/admin';

export const fishermanKeys = {
    all: ['fishermen'] as const,
};

export const useFishermanQuery = () => {
    const { data, error, isLoading } = useQuery({
        queryKey: fishermanKeys.all,
        queryFn: getFishermen,
    });

    return {
        fishermen: data ?? [],
        error,
        isLoading,
    };
};

export const useFishermanMutation = () => {
    const queryClient = useQueryClient();

    const createMutation = useMutation({
        mutationFn: (data: { name: string }) => registerFisherman(data.name),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: fishermanKeys.all });
        },
    });

    const deleteMutation = useMutation({
        mutationFn: deleteFisherman,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: fishermanKeys.all });
        },
    });

    return {
        createFisherman: createMutation.mutateAsync,
        isCreating: createMutation.isPending,
        deleteFisherman: deleteMutation.mutateAsync,
        isDeleting: deleteMutation.isPending,
        error: createMutation.error || deleteMutation.error,
    };
};
