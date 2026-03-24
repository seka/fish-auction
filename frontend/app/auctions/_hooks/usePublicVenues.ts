import { useQuery } from '@tanstack/react-query';
import { getVenues } from '@/src/api/venue';
import { venueKeys } from '@/src/hooks/venue/keys';

export const usePublicVenues = () => {
  const { data: venues } = useQuery({
    queryKey: venueKeys.publicAll,
    queryFn: getVenues,
  });
  return { venues };
};
