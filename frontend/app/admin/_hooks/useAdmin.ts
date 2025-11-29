import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query';
import { registerFisherman, registerBuyer, registerItem, getFishermen, getBuyers } from '@/src/api/admin';
import { RegisterItemParams } from '@/src/models';

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
        mutationFn: (data: { name: string }) => registerBuyer(data.name),
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
            queryClient.invalidateQueries({ queryKey: ['items'] });
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
        queryKey: ['fishermen'],
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
        queryKey: ['buyers'],
        queryFn: getBuyers,
    });

    return {
        buyers: data ?? [],
        error,
        isLoading,
    };
};
