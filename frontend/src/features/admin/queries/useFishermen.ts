import { useFishermanQuery } from '@/src/data/queries/adminFisherman/useQuery';
import { useFishermanMutation } from '@/src/data/queries/adminFisherman/useMutation';
import { toFisherman } from '../types/fisherman';

/**
 * 管理画面用漁師クエリフック
 */
export const useAdminFishermen = () => {
  const { data: fishermen, ...rest } = useFishermanQuery({
    select: (data) => data.map(toFisherman),
  });

  return {
    ...rest,
    fishermen: fishermen || [],
  };
};

/**
 * 管理画面用漁師ミューテーションフック
 */
export const useAdminFishermanMutations = () => {
  return useFishermanMutation();
};
