import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { getItemsByAuction } from '@/src/data/api/admin';
import { adminItemKeys } from './keys';
import { AuctionItem } from '@entities/auction';

/**
 * 管理用アイテム一覧取得フック
 * @param auctionId オークションID
 */
export const useItemQuery = <T = AuctionItem[]>(
  auctionId?: number,
  options?: Omit<UseQueryOptions<AuctionItem[], Error, T>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: auctionId ? ['admin', 'auctions', auctionId, 'items'] : adminItemKeys.all,
    queryFn: () => (auctionId ? getItemsByAuction(auctionId) : Promise.resolve([])),
    enabled: true,
    ...options,
  });
};
