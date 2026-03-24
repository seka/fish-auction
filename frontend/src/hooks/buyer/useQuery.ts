import { useQuery } from '@tanstack/react-query';
import { getBuyers } from '@/src/api/admin';
import { buyerKeys } from './keys';

export const useBuyerQuery = () => {
  const { data, error, isLoading } = useQuery({
    queryKey: buyerKeys.all,
    queryFn: getBuyers,
  });

  return {
    buyers: data ?? [],
    error,
    isLoading,
  };
};
