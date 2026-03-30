import { Auction as EntityAuction } from '@entities/auction';
import { useAuctionQuery } from '@/src/data/queries/adminAuction/useQuery';
import { useAuctionMutation } from '@/src/data/queries/adminAuction/useMutation';
import { Auction } from '../types/auction';

const toAuction = (entity: EntityAuction): Auction => ({
  id: entity.id,
  venueId: entity.venueId,
  auctionDate: entity.auctionDate,
  startTime: entity.startTime ?? null,
  endTime: entity.endTime ?? null,
  status: entity.status,
  createdAt: entity.createdAt,
  updatedAt: entity.updatedAt,
});

/**
 * 管理画面用オークションクエリフック
 */
export const useAdminAuctions = (params: { venueId?: number } = {}) => {
  const { data: auctions, ...rest } = useAuctionQuery(params, {
    select: (data) => data.map(toAuction),
  });

  return {
    ...rest,
    auctions: auctions || [],
  };
};

/**
 * 管理画面用オークションミューテーションフック
 */
export const useAdminAuctionMutations = () => {
  return useAuctionMutation();
};
