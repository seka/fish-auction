import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createVenue, updateVenue, deleteVenue } from '@/src/api/venue';
import { VenueFormData } from '@/src/models/schemas/auction';
import { venueKeys } from './keys';

export const useVenueMutation = () => {
  const queryClient = useQueryClient();

  const createMutation = useMutation({
    mutationFn: (data: VenueFormData) => createVenue(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: venueKeys.all });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: VenueFormData }) => updateVenue(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: venueKeys.all });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteVenue(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: venueKeys.all });
    },
  });

  return {
    createVenue: createMutation.mutateAsync,
    updateVenue: updateMutation.mutateAsync,
    deleteVenue: deleteMutation.mutateAsync,
    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
  };
};
