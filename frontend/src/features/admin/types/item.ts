import { AuctionItem as EntityAuctionItem } from '@entities/auction';
import { selectNextMinimumBid } from '../selectors/selectItem';

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
  sortOrder: number;
  createdAt: string;
}

const formatJPY = (value: number): string => {
  return `¥${value.toLocaleString('ja-JP')}`;
};

export const toAuctionItem = (entity: EntityAuctionItem): AuctionItem => {
  const highestBid = entity.highestBid ?? null;
  const nextMinBidValue = selectNextMinimumBid(highestBid ?? 0);

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
      value: highestBid ?? 0,
      label: highestBid !== null ? formatJPY(highestBid) : '-',
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
    sortOrder: entity.sortOrder,
    createdAt: entity.createdAt,
  };
};
