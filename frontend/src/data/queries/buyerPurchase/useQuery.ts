import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { getMyPurchases, Purchase } from '@/src/data/api/buyer_mypage';
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
