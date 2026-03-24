import { useQuery } from '@tanstack/react-query';
import { getMyPurchases } from '@/src/data/api/buyer_mypage';
import { buyerPurchaseKeys } from './keys';

export const useMyPurchases = () => {
  const {
    data: purchases = [],
    isLoading,
    error,
  } = useQuery({
    queryKey: buyerPurchaseKeys.meAll(),
    queryFn: getMyPurchases,
  });

  return { purchases, isLoading, error };
};
