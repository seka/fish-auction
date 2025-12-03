import { useMutation, useQueryClient } from '@tanstack/react-query';
import { registerItem } from '@/src/api/admin';
import { RegisterItemParams } from '@/src/models';

export const useItemMutations = () => {
    const queryClient = useQueryClient();

    const createMutation = useMutation({
        mutationFn: (item: RegisterItemParams) => registerItem(item),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['items'] });
        },
    });

    return {
        createItem: createMutation.mutateAsync,
        isCreating: createMutation.isPending,
        error: createMutation.error,
    };
};
