import { useAuctionQuery } from '@/src/data/queries/adminAuction/useQuery';
import { useAuctionMutation } from '@/src/data/queries/adminAuction/useMutation';

/**
 * 管理画面用オークションクエリフック
 */
export const useAdminAuctions = (params: { venueId?: number } = {}) => {
  return useAuctionQuery(params);
};

/**
 * 管理画面用オークションミューテーションフック
 */
export const useAdminAuctionMutations = () => {
  return useAuctionMutation();
};
