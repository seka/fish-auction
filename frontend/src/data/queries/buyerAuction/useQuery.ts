import { useQuery } from '@tanstack/react-query';
import { getMyAuctions } from '@/src/data/api/buyer_mypage';
import { buyerAuctionKeys } from './keys';

export const useParticipatingAuctions = () => {
  const {
    data: auctions = [],
    isLoading,
    error,
  } = useQuery({
    queryKey: buyerAuctionKeys.meAll(),
    queryFn: getMyAuctions,
  });

  return { auctions, isLoading, error };
};
