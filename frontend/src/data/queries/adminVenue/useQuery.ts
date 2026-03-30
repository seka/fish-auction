import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { getVenues } from '@/src/data/api/venue';
import { Venue } from '@entities/venue';
import { adminVenueKeys } from './keys';

export const useVenueQuery = <T = Venue[]>(
  options?: Omit<UseQueryOptions<Venue[], Error, T>, 'queryKey' | 'queryFn'>,
) => {
  return useQuery({
    queryKey: adminVenueKeys.all,
    queryFn: getVenues,
    ...options,
  });
};
