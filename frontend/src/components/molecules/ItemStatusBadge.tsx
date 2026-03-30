import { useTranslations } from 'next-intl';
import { Badge } from '@atoms';

type ItemStatus = 'Pending' | 'Sold' | 'Unsold' | 'Bidding';

interface ItemStatusBadgeProps {
  status: ItemStatus;
}

export const ItemStatusBadge = ({ status }: ItemStatusBadgeProps) => {
  const t = useTranslations('ItemStatus');

  const config: Record<
    ItemStatus,
    { variant: 'success' | 'warning' | 'error' | 'info' | 'neutral' }
  > = {
    Pending: { variant: 'info' },
    Bidding: { variant: 'success' },
    Sold: { variant: 'neutral' },
    Unsold: { variant: 'error' },
  };

  const { variant } = config[status] || { variant: 'neutral' };

  return <Badge variant={variant}>{t(status)}</Badge>;
};
