import { AuctionItem as EntityAuctionItem } from '@entities/auction';
import { selectNextMinimumBid } from '../selectors/selectAuction';
import { selectItemStatus } from '../selectors/selectItem';

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
  status: {
    value: 'Pending' | 'Bidding' | 'Sold' | 'Unsold';
    labelKey: string;
    variant: 'success' | 'warning' | 'error' | 'info' | 'neutral';
    isPending: boolean;
    isBidding: boolean;
    isSold: boolean;
    isUnsold: boolean;
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
  const itemStatus = entity.status;
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
    status: selectItemStatus(itemStatus),
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
