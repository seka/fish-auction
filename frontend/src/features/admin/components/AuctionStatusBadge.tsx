import { useTranslations } from 'next-intl';
import { Badge } from '@atoms';

interface AuctionStatusBadgeProps {
  status: {
    labelKey: string;
    variant: 'success' | 'warning' | 'error' | 'info' | 'neutral';
  };
}

export const AuctionStatusBadge = ({ status }: AuctionStatusBadgeProps) => {
  const t = useTranslations('AuctionStatus');

  return <Badge variant={status.variant}>{t(status.labelKey)}</Badge>;
};
