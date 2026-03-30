import { useParticipatingAuctions as useDataParticipatingAuctions } from '@/src/data/queries/buyerAuction/useQuery';
import { AuctionSummary } from '@/src/data/api/buyer_mypage';
import { Auction } from '../types/auction';

export const toAuction = (entity: AuctionSummary): Auction => ({
  id: entity.id,
  status: entity.status,
  auctionDate: entity.auctionDate,
  startTime: entity.startTime ?? null,
  endTime: entity.endTime ?? null,
});

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
