import { useFishermanQuery } from '@/src/data/queries/adminFisherman/useQuery';
import { useFishermanMutation } from '@/src/data/queries/adminFisherman/useMutation';

/**
 * 管理画面用漁師クエリフック
 */
export const useAdminFishermen = () => {
  return useFishermanQuery();
};

/**
 * 管理画面用漁師ミューテーションフック
 */
export const useAdminFishermanMutations = () => {
  return useFishermanMutation();
};
