import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getAuctions, createAuction, updateAuction, updateAuctionStatus, deleteAuction } from '@/src/api/auction';
import { AuctionFormData } from '@/src/models/schemas/auction';

export const useAuctions = (filters?: { venue_id?: number; date?: string; status?: string }) => {
    const { data: auctions, isLoading, error } = useQuery({
        queryKey: ['auctions', filters],
        queryFn: () => getAuctions(filters),
    });

    return { auctions: auctions || [], isLoading, error };
};

export const useAuctionMutations = () => {
    const queryClient = useQueryClient();

    const createMutation = useMutation({
        mutationFn: (data: AuctionFormData) => createAuction(data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['auctions'] });
        },
    });

    const updateMutation = useMutation({
        mutationFn: ({ id, data }: { id: number; data: AuctionFormData }) => updateAuction(id, data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['auctions'] });
        },
    });

    const updateStatusMutation = useMutation({
        mutationFn: ({ id, status }: { id: number; status: string }) => updateAuctionStatus(id, status),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['auctions'] });
        },
    });

    const deleteMutation = useMutation({
        mutationFn: (id: number) => deleteAuction(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['auctions'] });
        },
    });

    return {
        createAuction: createMutation.mutateAsync,
        updateAuction: updateMutation.mutateAsync,
        updateStatus: updateStatusMutation.mutateAsync,
        deleteAuction: deleteMutation.mutateAsync,
        isCreating: createMutation.isPending,
        isUpdating: updateMutation.isPending,
        isUpdatingStatus: updateStatusMutation.isPending,
        isDeleting: deleteMutation.isPending,
    };
};
