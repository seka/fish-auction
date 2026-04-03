import { Auction as EntityAuction } from '@entities/auction';
import {
  selectIsAuctionActive,
  selectTimeLabel,
  selectAuctionStatus,
  toJSTDate,
} from '../selectors/selectAuction';

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
    startAt: Date;
    endAt: Date;
    dateLabel: string;
    startTime: string | null;
    endTime: string | null;
    label: string;
  };
  isActive: boolean;
}

export const toAuction = (entity: EntityAuction): Auction => {
  const auctionDate = entity.auctionDate;
  const startTime = entity.startTime || '00:00:00';
  const endTime = entity.endTime || '23:59:59';

  return {
    id: entity.id,
    venueId: entity.venueId,
    status: selectAuctionStatus(entity.status),
    duration: {
      startAt: toJSTDate(auctionDate, startTime),
      endAt: toJSTDate(auctionDate, endTime),
      dateLabel: auctionDate,
      startTime: entity.startTime ?? null,
      endTime: entity.endTime ?? null,
      label: selectTimeLabel(entity.startTime ?? null, entity.endTime ?? null),
    },
    isActive: selectIsAuctionActive(entity),
  };
};
