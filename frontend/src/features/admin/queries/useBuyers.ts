import { useBuyerQuery } from '@/src/data/queries/adminBuyer/useQuery';
import { useBuyerMutation } from '@/src/data/queries/adminBuyer/useMutation';
import { toBuyer } from '../types/buyer';

/**
 * 管理画面用買受人クエリフック
 */
export const useAdminBuyers = () => {
  const { data: buyers, ...rest } = useBuyerQuery({
    select: (data) => data.map(toBuyer),
  });

  return {
    ...rest,
    buyers: buyers || [],
  };
};

/**
 * 管理画面用買受人ミューテーションフック
 */
export const useAdminBuyerMutations = () => {
  return useBuyerMutation();
};
