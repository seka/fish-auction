import { useVenueQuery } from '@/src/data/queries/publicVenue/useQuery';
import { toVenue } from '../types/venue';

/**
 * オークション一覧用 会場クエリフック
 */
export const useVenues = () => {
  const { data: venues, ...rest } = useVenueQuery({
    select: (data) => data.map(toVenue),
  });

  return {
    ...rest,
    venues: venues || [],
  };
};
