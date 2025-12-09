import { useMutation, useQueryClient } from '@tanstack/react-query';
import { registerItem } from '@/src/api/admin';
import { RegisterItemParams } from '@/src/models';

export const itemKeys = {
    all: ['items'] as const,
};

export const useItemMutation = () => {
    const queryClient = useQueryClient();

    const createMutation = useMutation({
        mutationFn: (item: RegisterItemParams) => registerItem(item),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: itemKeys.all });
        },
    });

    return {
        createItem: createMutation.mutateAsync,
        isCreating: createMutation.isPending,
        error: createMutation.error,
    };
};
