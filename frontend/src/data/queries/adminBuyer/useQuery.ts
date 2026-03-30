import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { getBuyers } from '@/src/data/api/admin';
import { Buyer } from '@entities/admin';
import { adminBuyerKeys } from './keys';

export const useBuyerQuery = <T = Buyer[]>(
  options?: Omit<UseQueryOptions<Buyer[], Error, T>, 'queryKey' | 'queryFn'>,
) => {
  return useQuery({
    queryKey: adminBuyerKeys.all,
    queryFn: getBuyers,
    ...options,
  });
};
