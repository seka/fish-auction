import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query';
import { registerItem, updateItem, deleteItem, updateItemSortOrder, getItemsByAuction } from '@/src/api/admin';
import { RegisterItemParams, UpdateItemParams, UpdateItemSortOrderParams } from '@/src/models';

export const itemKeys = {
    all: ['items'] as const,
    byAuction: (auctionId: number) => [...itemKeys.all, 'auction', auctionId] as const,
};

export const useItemMutation = () => {
    const queryClient = useQueryClient();

    const createMutation = useMutation({
        mutationFn: registerItem,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: itemKeys.all });
        },
    });

    const updateMutation = useMutation({
        mutationFn: updateItem,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: itemKeys.all });
        },
    });

    const deleteMutation = useMutation({
        mutationFn: deleteItem,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: itemKeys.all });
        },
    });

    const sortMutation = useMutation({
        mutationFn: updateItemSortOrder,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: itemKeys.all });
        },
    });

    return {
        createItem: createMutation.mutateAsync,
        isCreating: createMutation.isPending,
        updateItem: updateMutation.mutateAsync,
        isUpdating: updateMutation.isPending,
        deleteItem: deleteMutation.mutateAsync,
        isDeleting: deleteMutation.isPending,
        updateSortOrder: sortMutation.mutateAsync,
        isSorting: sortMutation.isPending,
    };
};

export const useItemQuery = (auctionId?: number) => {
    return useQuery({
        queryKey: auctionId ? itemKeys.byAuction(auctionId) : itemKeys.all,
        queryFn: () => auctionId ? getItemsByAuction(auctionId) : Promise.resolve([]),
        enabled: !!auctionId,
    });
};
