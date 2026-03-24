import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  registerItem,
  updateItem,
  deleteItem,
  updateItemSortOrder,
  reorderItems,
} from '@/src/data/api/admin';
import { ReorderItemsParams, AuctionItem } from '@/src/models';
import { itemKeys } from '../publicItem/keys'; // For public invalidation
import { adminItemKeys } from './keys';

export const useItemMutation = () => {
  const queryClient = useQueryClient();

  const createMutation = useMutation({
    mutationFn: registerItem,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminItemKeys.all });
      queryClient.invalidateQueries({ queryKey: itemKeys.publicAll });
    },
  });

  const updateMutation = useMutation({
    mutationFn: updateItem,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminItemKeys.all });
      queryClient.invalidateQueries({ queryKey: itemKeys.publicAll });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteItem,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminItemKeys.all });
      queryClient.invalidateQueries({ queryKey: itemKeys.publicAll });
    },
  });

  const reorderMutation = useMutation({
    mutationFn: (params: ReorderItemsParams) => reorderItems(params.auctionId, params.ids),
    onMutate: async (params) => {
      const queryKey = itemKeys.publicByAuction(params.auctionId);
      await queryClient.cancelQueries({ queryKey });

      const previousItems = queryClient.getQueryData<AuctionItem[]>(queryKey);

      if (previousItems) {
        // Rearrange items based on the new ids order
        const idToIndex = new Map(params.ids.map((id, index) => [id, index]));
        const newItems = [...previousItems].sort((a, b) => {
          const indexA = idToIndex.get(a.id) ?? 999;
          const indexB = idToIndex.get(b.id) ?? 999;
          return indexA - indexB;
        });
        queryClient.setQueryData(queryKey, newItems);
      }

      return { previousItems, queryKey };
    },
    onError: (err, variables, context) => {
      if (context?.previousItems) {
        queryClient.setQueryData(context.queryKey, context.previousItems);
      }
    },
    onSettled: (data, error, variables, context) => {
      queryClient.invalidateQueries({ queryKey: context?.queryKey || adminItemKeys.all });
      queryClient.invalidateQueries({ queryKey: itemKeys.publicAll });
    },
  });

  const sortMutation = useMutation({
    mutationFn: updateItemSortOrder,
    onMutate: async (newOrder) => {
      const queryKey = itemKeys.publicByAuction(newOrder.auctionId);
      await queryClient.cancelQueries({ queryKey });

      const previousItems = queryClient.getQueryData<AuctionItem[]>(queryKey);

      if (previousItems) {
        const oldIndex = previousItems.findIndex((i) => i.id === newOrder.id);
        if (oldIndex !== -1) {
          const newItems = [...previousItems];
          const [movedItem] = newItems.splice(oldIndex, 1);
          newItems.splice(newOrder.newIndex, 0, { ...movedItem, sortOrder: newOrder.sortOrder });
          queryClient.setQueryData(queryKey, newItems);
        }
      }

      return { previousItems, queryKey };
    },
    onError: (err, newOrder, context) => {
      if (context?.previousItems) {
        queryClient.setQueryData(context.queryKey, context.previousItems);
      }
    },
    onSettled: (data, error, variables, context) => {
      queryClient.invalidateQueries({ queryKey: context?.queryKey || adminItemKeys.all });
      queryClient.invalidateQueries({ queryKey: itemKeys.publicAll });
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
    reorderItems: reorderMutation.mutateAsync,
    isReordering: reorderMutation.isPending,
  };
};
