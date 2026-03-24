import { useQuery } from '@tanstack/react-query';
import { getItemsByAuction } from '@/src/data/api/admin';
import { adminItemKeys } from './keys';

/**
 * 管理用アイテム一覧取得フック
 * @param auctionId オークションID（指定がない場合は全件取得。ただし現在はAPI制限あり）
 */
export const useItemQuery = (auctionId?: number) => {
  return useQuery({
    queryKey: auctionId ? ['admin', 'auctions', auctionId, 'items'] : adminItemKeys.all,
    queryFn: () => (auctionId ? getItemsByAuction(auctionId) : Promise.resolve([])),
    enabled: true,
  });
};
