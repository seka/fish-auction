import { Auction as EntityAuction } from '@entities/auction';
import { selectAuctionStatus, selectTimeLabel } from '../selectors/selectAuction';

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
  actions: {
    canStart: boolean;
    canFinish: boolean;
  };
  createdAt: string;
  updatedAt: string;
}

export const toAuction = (entity: EntityAuction): Auction => {
  const startAt = entity.startAt ? new Date(entity.startAt) : null;
  const endAt = entity.endAt ? new Date(entity.endAt) : null;

  const status = selectAuctionStatus(entity.status);

  return {
    id: entity.id,
    venueId: entity.venueId,
    status,
    duration: {
      startAt,
      endAt,
      label: selectTimeLabel(startAt, endAt),
    },
    actions: {
      canStart: status.isScheduled,
      canFinish: status.isInProgress,
    },
    createdAt: entity.createdAt,
    updatedAt: entity.updatedAt,
  };
};
