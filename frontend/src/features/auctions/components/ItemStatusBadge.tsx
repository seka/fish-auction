import { useTranslations } from 'next-intl';
import { Badge } from '@atoms';

interface ItemStatusBadgeProps {
  status: {
    labelKey: string;
    variant: 'success' | 'warning' | 'error' | 'info' | 'neutral';
  };
}

export const ItemStatusBadge = ({ status }: ItemStatusBadgeProps) => {
  const t = useTranslations('ItemStatus');

  return <Badge variant={status.variant}>{t(status.labelKey)}</Badge>;
};
