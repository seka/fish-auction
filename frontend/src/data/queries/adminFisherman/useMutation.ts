import { useMutation, useQueryClient } from '@tanstack/react-query';
import { registerFisherman, deleteFisherman } from '@/src/data/api/admin';
import { adminFishermanKeys } from './keys';

export const useFishermanMutation = () => {
  const queryClient = useQueryClient();

  const createMutation = useMutation({
    mutationFn: (data: { name: string }) => registerFisherman(data.name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminFishermanKeys.all });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteFisherman,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminFishermanKeys.all });
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
