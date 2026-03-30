import { Venue as EntityVenue } from '@entities/venue';
import { useVenueQuery } from '@/src/data/queries/adminVenue/useQuery';
import { useVenueMutation } from '@/src/data/queries/adminVenue/useMutation';
import { Venue } from '../types/venue';

export const toVenue = (entity: EntityVenue): Venue => ({
  id: entity.id,
  name: entity.name,
  location: entity.location,
  description: entity.description,
  createdAt: entity.createdAt,
});

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
