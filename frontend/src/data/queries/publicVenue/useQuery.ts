import { useQuery } from '@tanstack/react-query';
import { getVenues } from '@/src/data/api/venue';
import { venueKeys } from '@/src/data/queries/publicVenue/keys';

export const useVenueQuery = () => {
  const { data: venues } = useQuery({
    queryKey: venueKeys.publicAll,
    queryFn: getVenues,
  });
  return { venues: venues || [] };
};
