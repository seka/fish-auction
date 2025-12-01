import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getVenues, createVenue, updateVenue, deleteVenue } from '@/src/api/venue';
import { VenueFormData } from '@/src/models/schemas/auction';

export const useVenues = () => {
    const { data: venues, isLoading, error } = useQuery({
        queryKey: ['venues'],
        queryFn: getVenues,
    });

    return { venues: venues || [], isLoading, error };
};

export const useVenueMutations = () => {
    const queryClient = useQueryClient();

    const createMutation = useMutation({
        mutationFn: (data: VenueFormData) => createVenue(data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['venues'] });
        },
    });

    const updateMutation = useMutation({
        mutationFn: ({ id, data }: { id: number; data: VenueFormData }) => updateVenue(id, data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['venues'] });
        },
    });

    const deleteMutation = useMutation({
        mutationFn: (id: number) => deleteVenue(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['venues'] });
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
