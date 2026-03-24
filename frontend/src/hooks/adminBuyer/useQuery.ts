import { useQuery } from '@tanstack/react-query';
import { getBuyers } from '@/src/api/admin';
import { adminBuyerKeys } from './keys';

export const useBuyerQuery = () => {
  const { data, error, isLoading } = useQuery({
    queryKey: adminBuyerKeys.all,
    queryFn: getBuyers,
  });

  return {
    buyers: data ?? [],
    error,
    isLoading,
  };
};
