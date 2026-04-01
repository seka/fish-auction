import { Buyer as EntityBuyer } from '@entities/admin';

export interface Buyer {
  id?: number;
  name: string;
}

export const toBuyer = (entity: EntityBuyer): Buyer => ({
  id: entity.id,
  name: entity.name,
});
