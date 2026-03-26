import { useVenueQuery } from '@/src/data/queries/adminVenue/useQuery';
import { useVenueMutation } from '@/src/data/queries/adminVenue/useMutation';

/**
 * 管理画面用会場クエリフック
 */
export const useAdminVenues = () => {
  return useVenueQuery();
};

/**
 * 管理画面用会場ミューテーションフック
 */
export const useAdminVenueMutations = () => {
  return useVenueMutation();
};
