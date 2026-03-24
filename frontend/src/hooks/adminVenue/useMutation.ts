import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createVenue, updateVenue, deleteVenue } from '@/src/api/venue';
import { VenueFormData } from '@/src/models/schemas/auction';
import { venueKeys } from '../venue/keys'; // For public invalidation
import { adminVenueKeys } from './keys';

export const useVenueMutation = () => {
  const queryClient = useQueryClient();

  const createMutation = useMutation({
    mutationFn: (data: VenueFormData) => createVenue(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminVenueKeys.all });
      queryClient.invalidateQueries({ queryKey: venueKeys.publicAll });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: VenueFormData }) => updateVenue(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminVenueKeys.all });
      queryClient.invalidateQueries({ queryKey: venueKeys.publicAll });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteVenue(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminVenueKeys.all });
      queryClient.invalidateQueries({ queryKey: venueKeys.publicAll });
    },
  });

  return {
    createVenue: createMutation.mutateAsync,
    isCreating: createMutation.isPending,
    updateVenue: updateMutation.mutateAsync,
    isUpdating: updateMutation.isPending,
    deleteVenue: deleteMutation.mutateAsync,
    isDeleting: deleteMutation.isPending,
  };
};
