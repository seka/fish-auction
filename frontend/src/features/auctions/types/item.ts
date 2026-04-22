import { AuctionItem as EntityAuctionItem } from '@entities/auction';
import { selectNextMinimumBid } from '../selectors/selectAuction';

export interface AuctionItem {
  id: number;
  auctionId: number;
  fishermanId: number;
  fishType: string;
  quantity: {
    value: number;
    label: string;
  };
  unit: string;
  price: {
    value: number;
    label: string;
  };
  bidding: {
    highestBid: number | null;
    highestBidderId: number | null;
    highestBidderName: string | null;
    nextMinBid: {
      value: number;
      label: string;
    };
  };
}

const formatJPY = (value: number): string => {
  return `¥${value.toLocaleString('ja-JP')}`;
};

export const toAuctionItem = (entity: EntityAuctionItem): AuctionItem => {
  const highestBid = entity.highestBid || 0;
  const nextMinBidValue = selectNextMinimumBid(highestBid);

  return {
    id: entity.id,
    auctionId: entity.auctionId,
    fishermanId: entity.fishermanId,
    fishType: entity.fishType,
    quantity: {
      value: entity.quantity,
      label: `${entity.quantity} ${entity.unit}`,
    },
    unit: entity.unit,
    price: {
      value: highestBid,
      label: formatJPY(highestBid),
    },
    bidding: {
      highestBid: entity.highestBid ?? null,
      highestBidderId: entity.highestBidderId ?? null,
      highestBidderName: entity.highestBidderName ?? null,
      nextMinBid: {
        value: nextMinBidValue,
        label: formatJPY(nextMinBidValue),
      },
    },
  };
};
