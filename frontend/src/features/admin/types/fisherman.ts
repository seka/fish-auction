import { Fisherman as EntityFisherman } from '@entities/admin';

export interface Fisherman {
  id?: number;
  name: string;
}

export const toFisherman = (entity: EntityFisherman): Fisherman => ({
  id: entity.id,
  name: entity.name,
});
