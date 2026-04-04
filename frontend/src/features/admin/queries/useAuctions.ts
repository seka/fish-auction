'use client';
import { useAuctionQuery } from '@/src/data/queries/adminAuction/useQuery';
import { useAuctionMutation } from '@/src/data/queries/adminAuction/useMutation';
import { toAuction } from '../types';

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
