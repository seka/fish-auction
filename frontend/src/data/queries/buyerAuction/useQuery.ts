import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { getMyAuctions } from '@/src/data/api/buyer_mypage';
import { AuctionSummary } from '@entities';
import { buyerAuctionKeys } from './keys';

export const useParticipatingAuctions = <T = AuctionSummary[]>(
  options?: Omit<UseQueryOptions<AuctionSummary[], Error, T>, 'queryKey' | 'queryFn'>,
) => {
  return useQuery({
    queryKey: buyerAuctionKeys.meAll(),
    queryFn: getMyAuctions,
    ...options,
  });
};
