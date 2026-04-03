import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { getMyPurchases } from '@/src/data/api/buyer_mypage';
import { Purchase } from '@entities';
import { buyerPurchaseKeys } from './keys';

export const useMyPurchases = <T = Purchase[]>(
  options?: Omit<UseQueryOptions<Purchase[], Error, T>, 'queryKey' | 'queryFn'>,
) => {
  return useQuery({
    queryKey: buyerPurchaseKeys.meAll(),
    queryFn: getMyPurchases,
    ...options,
  });
};
