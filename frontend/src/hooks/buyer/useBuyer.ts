import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getBuyers, registerBuyer, deleteBuyer } from '@/src/api/admin';
import { buyerKeys } from './queryKey';

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
    mutationFn: (data: {
      name: string;
      email: string;
      password: string;
      organization: string;
      contactInfo: string;
    }) => registerBuyer(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: buyerKeys.all });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteBuyer,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: buyerKeys.all });
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
