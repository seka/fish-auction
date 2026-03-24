import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query';
import {
  registerFisherman,
  registerBuyer,
  registerItem,
  getFishermen,
  getBuyers,
} from '@/src/api/admin';
import { RegisterItemParams } from '@/src/models';
import { BuyerFormData } from '@/src/models/schemas/admin';
import { itemKeys } from '@/src/hooks/item/queryKey';
import { fishermanKeys } from '@/src/hooks/fisherman/queryKey';
import { buyerKeys } from '@/src/hooks/buyer/queryKey';

export const useRegisterFisherman = () => {
  const mutation = useMutation({
    mutationFn: (data: { name: string }) => registerFisherman(data.name),
  });

  return {
    registerFisherman: mutation.mutateAsync,
    isLoading: mutation.isPending,
    error: mutation.error,
  };
};

export const useRegisterBuyer = () => {
  const mutation = useMutation({
    mutationFn: (data: BuyerFormData) => registerBuyer(data),
  });

  return {
    registerBuyer: mutation.mutateAsync,
    isLoading: mutation.isPending,
    error: mutation.error,
  };
};

export const useRegisterItem = () => {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: (item: RegisterItemParams) => registerItem(item),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: itemKeys.all });
    },
  });

  return {
    registerItem: mutation.mutateAsync,
    isLoading: mutation.isPending,
    error: mutation.error,
  };
};

export const useFishermen = () => {
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

export const useBuyers = () => {
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
