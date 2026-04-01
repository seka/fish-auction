import { Venue as EntityVenue } from '@entities/venue';

export interface Venue {
  id: number;
  name: string;
  location?: string;
  description?: string;
  createdAt: string;
}

export const toVenue = (entity: EntityVenue): Venue => ({
  id: entity.id,
  name: entity.name,
  location: entity.location,
  description: entity.description,
  createdAt: entity.createdAt,
});
