import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { getFishermen } from '@/src/data/api/admin';
import { Fisherman } from '@entities/admin';
import { adminFishermanKeys } from './keys';

export const useFishermanQuery = <T = Fisherman[]>(
  options?: Omit<UseQueryOptions<Fisherman[], Error, T>, 'queryKey' | 'queryFn'>,
) => {
  return useQuery({
    queryKey: adminFishermanKeys.all,
    queryFn: getFishermen,
    ...options,
  });
};
