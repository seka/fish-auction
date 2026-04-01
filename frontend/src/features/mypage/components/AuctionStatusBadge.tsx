import { useTranslations } from 'next-intl';
import { Badge } from '@atoms';
import { AuctionStatus } from '../types/auction';

interface AuctionStatusBadgeProps {
  status: AuctionStatus;
}

export const AuctionStatusBadge = ({ status }: AuctionStatusBadgeProps) => {
  const t = useTranslations('AuctionStatus');

  const config: Record<
    AuctionStatus,
    { variant: 'success' | 'warning' | 'error' | 'info' | 'neutral' }
  > = {
    scheduled: { variant: 'info' },
    in_progress: { variant: 'success' },
    completed: { variant: 'neutral' },
    cancelled: { variant: 'error' },
  };

  const { variant } = config[status] || { variant: 'neutral' };

  return <Badge variant={variant}>{t(status)}</Badge>;
};
