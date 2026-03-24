import { useQuery } from '@tanstack/react-query';
import { getVenues } from '@/src/api/venue';
import { venueKeys } from './keys';

export const useVenueQuery = () => {
  const {
    data: venues,
    isLoading,
    error,
  } = useQuery({
    queryKey: venueKeys.all,
    queryFn: getVenues,
  });

  return { venues: venues || [], isLoading, error };
};
