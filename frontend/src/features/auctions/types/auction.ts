import { Auction as EntityAuction } from '@entities/auction';
import { selectIsAuctionActive, selectTimeLabel, selectAuctionStatus } from '../selectors/selectAuction';

export interface Auction {
  id: number;
  venueId: number;
  status: {
    value: 'scheduled' | 'in_progress' | 'completed' | 'cancelled';
    labelKey: string;
    variant: 'success' | 'warning' | 'error' | 'info' | 'neutral';
    isScheduled: boolean;
    isInProgress: boolean;
    isCompleted: boolean;
    isCancelled: boolean;
  };
  duration: {
    startAt: Date | null;
    endAt: Date | null;
    label: string;
  };
  isActive: boolean;
}

export type AuctionStatus = Auction['status'];

export const toAuction = (entity: EntityAuction): Auction => {
  const startAt = entity.startAt ? new Date(entity.startAt) : null;
  const endAt = entity.endAt ? new Date(entity.endAt) : null;

  return {
    id: entity.id,
    venueId: entity.venueId,
    status: selectAuctionStatus(entity.status),
    duration: {
      startAt,
      endAt,
      label: selectTimeLabel(startAt, endAt),
    },
    isActive: selectIsAuctionActive(entity),
  };
};
