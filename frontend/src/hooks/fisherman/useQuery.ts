import { useQuery } from '@tanstack/react-query';
import { getFishermen } from '@/src/api/admin';
import { fishermanKeys } from './keys';

export const useFishermanQuery = () => {
  const { data, error, isLoading } = useQuery({
    queryKey: fishermanKeys.all,
    queryFn: getFishermen,
  });

  return {
    fishermen: data ?? [],
    error,
    isLoading,
  };
};
