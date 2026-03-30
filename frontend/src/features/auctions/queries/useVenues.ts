import { useVenueQuery } from '@/src/data/queries/publicVenue/useQuery';
import { Venue as EntityVenue } from '@entities/venue';
import { Venue } from '../types';

export const toVenue = (entity: EntityVenue): Venue => ({
  id: entity.id,
  name: entity.name,
  location: entity.location,
  description: entity.description,
  createdAt: entity.createdAt,
});

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
