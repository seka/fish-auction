import { useMutation, useQueryClient } from '@tanstack/react-query';
import { registerBuyer, deleteBuyer } from '@/src/api/admin';
import { adminBuyerKeys } from './keys';

export const useBuyerMutation = () => {
  const queryClient = useQueryClient();

  const createMutation = useMutation({
    mutationFn: (data: {
      name: string;
      email: string;
      password: string;
      organization: string;
      contactInfo: string;
    }) => registerBuyer(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminBuyerKeys.all });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteBuyer,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminBuyerKeys.all });
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
