import { useItemQuery } from '@/src/data/queries/adminItem/useQuery';
import { useItemMutation } from '@/src/data/queries/adminItem/useMutation';

/**
 * 管理画面用商品クエリフック
 */
export const useAdminItems = (auctionId?: number) => {
  return useItemQuery(auctionId);
};

/**
 * 管理画面用商品ミューテーションフック
 */
export const useAdminItemMutations = () => {
  return useItemMutation();
};
