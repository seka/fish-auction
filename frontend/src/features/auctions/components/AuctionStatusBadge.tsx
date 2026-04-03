import { useTranslations } from 'next-intl';
import { Badge } from '@atoms';
import { AuctionStatus } from '../types/auction';

interface AuctionStatusBadgeProps {
  status: AuctionStatus;
}

export const AuctionStatusBadge = ({ status }: AuctionStatusBadgeProps) => {
  const t = useTranslations('AuctionStatus');

  return <Badge variant={status.variant}>{t(status.value)}</Badge>;
};
