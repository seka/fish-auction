import { useQuery } from '@tanstack/react-query';
import { getFishermen } from '@/src/api/admin';
import { adminFishermanKeys } from './keys';

export const useFishermanQuery = () => {
  const { data, error, isLoading } = useQuery({
    queryKey: adminFishermanKeys.all,
    queryFn: getFishermen,
  });

  return {
    fishermen: data ?? [],
    error,
    isLoading,
  };
};
