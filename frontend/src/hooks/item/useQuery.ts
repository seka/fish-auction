import { useQuery } from '@tanstack/react-query';
import { getItemsByAuction } from '@/src/api/admin';
import { itemKeys } from './keys';

export const useItemQuery = (auctionId?: number) => {
  return useQuery({
    queryKey: auctionId ? itemKeys.byAuction(auctionId) : itemKeys.all,
    queryFn: () => (auctionId ? getItemsByAuction(auctionId) : Promise.resolve([])),
    enabled: !!auctionId,
  });
};
