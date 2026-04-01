import { useVenueQuery } from '@/src/data/queries/adminVenue/useQuery';
import { useVenueMutation } from '@/src/data/queries/adminVenue/useMutation';
import { toVenue } from '../types/venue';

/**
 * 管理画面用会場クエリフック
 */
export const useAdminVenues = () => {
  const { data: venues, ...rest } = useVenueQuery({
    select: (data) => data.map(toVenue),
  });

  return {
    ...rest,
    venues: venues || [],
  };
};

/**
 * 管理画面用会場ミューテーションフック
 */
export const useAdminVenueMutations = () => {
  return useVenueMutation();
};
