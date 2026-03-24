import { useQuery } from '@tanstack/react-query';
import { getVenues } from '@/src/api/venue';
import { adminVenueKeys } from './keys';

export const useVenueQuery = () => {
  const {
    data: venues,
    isLoading,
    error,
  } = useQuery({
    queryKey: adminVenueKeys.all,
    queryFn: getVenues,
  });

  return { venues: venues || [], isLoading, error };
};
