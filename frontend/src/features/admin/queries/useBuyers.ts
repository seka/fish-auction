import { useBuyerQuery } from '@/src/data/queries/adminBuyer/useQuery';
import { useBuyerMutation } from '@/src/data/queries/adminBuyer/useMutation';

/**
 * 管理画面用買受人クエリフック
 */
export const useAdminBuyers = () => {
  return useBuyerQuery();
};

/**
 * 管理画面用買受人ミューテーションフック
 */
export const useAdminBuyerMutations = () => {
  return useBuyerMutation();
};
