import { useQuery } from '@tanstack/react-query';
import { getVenues } from '@/src/api/venue';
import { venueKeys } from '@/src/hooks/venue/queryKey';

export const usePublicVenues = () => {
  const { data: venues } = useQuery({
    queryKey: venueKeys.all,
    queryFn: getVenues,
  });
  return { venues };
};
