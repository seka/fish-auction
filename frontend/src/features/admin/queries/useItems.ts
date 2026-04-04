'use client';
import { useItemQuery } from '@/src/data/queries/adminItem/useQuery';
import { useItemMutation } from '@/src/data/queries/adminItem/useMutation';
import { toAuctionItem } from '../types/item';

/**
 * 管理画面用商品クエリフック
 */
export const useAdminItems = (auctionId?: number) => {
  const { data: items, ...rest } = useItemQuery(auctionId, {
    select: (data) => data.map(toAuctionItem),
  });

  return {
    ...rest,
    items: items || [],
  };
};

/**
 * 管理画面用商品ミューテーションフック
 */
export const useAdminItemMutations = () => {
  return useItemMutation();
};
