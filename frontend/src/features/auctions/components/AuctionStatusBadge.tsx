import { useTranslations } from 'next-intl';
import { Badge } from '@atoms';
import { Auction } from '../types/auction';

interface AuctionStatusBadgeProps {
  status: Auction['status'];
}

export const AuctionStatusBadge = ({ status }: AuctionStatusBadgeProps) => {
  const t = useTranslations('AuctionStatus');

  return <Badge variant={status.variant}>{t(status.labelKey)}</Badge>;
};
