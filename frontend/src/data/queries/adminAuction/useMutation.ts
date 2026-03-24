import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  createAuction,
  updateAuction,
  updateAuctionStatus,
  deleteAuction,
} from '@/src/data/api/auction';
import { AuctionFormData } from '@/src/models/schemas/auction';
import { auctionKeys } from '../publicAuction/keys'; // For public invalidation
import { adminAuctionKeys } from './keys';

export const useAuctionMutation = () => {
  const queryClient = useQueryClient();

  const createMutation = useMutation({
    mutationFn: (data: AuctionFormData) => createAuction(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminAuctionKeys.all });
      queryClient.invalidateQueries({ queryKey: auctionKeys.publicAll });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: AuctionFormData }) => updateAuction(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminAuctionKeys.all });
      queryClient.invalidateQueries({ queryKey: auctionKeys.publicAll });
    },
  });

  const updateStatusMutation = useMutation({
    mutationFn: ({ id, status }: { id: number; status: string }) => updateAuctionStatus(id, status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminAuctionKeys.all });
      queryClient.invalidateQueries({ queryKey: auctionKeys.publicAll });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteAuction(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminAuctionKeys.all });
      queryClient.invalidateQueries({ queryKey: auctionKeys.publicAll });
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
