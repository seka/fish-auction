import { Buyer as EntityBuyer } from '@entities/admin';
import { useBuyerQuery } from '@/src/data/queries/adminBuyer/useQuery';
import { useBuyerMutation } from '@/src/data/queries/adminBuyer/useMutation';
import { Buyer } from '../types/buyer';

export const toBuyer = (entity: EntityBuyer): Buyer => ({
  id: entity.id,
  name: entity.name,
});

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
