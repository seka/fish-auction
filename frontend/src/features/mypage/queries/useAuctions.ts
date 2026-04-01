import { useParticipatingAuctions as useDataParticipatingAuctions } from '@/src/data/queries/buyerAuction/useQuery';
import { toAuction } from '../types/auction';

export const useParticipatingAuctions = () => {
  const { data: auctions, isLoading, ...rest } = useDataParticipatingAuctions({
    select: (data) => data.map(toAuction),
  });

  return {
    auctions: auctions || [],
    isLoading,
    ...rest,
  };
};
