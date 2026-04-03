import { useTranslations } from 'next-intl';
import { Badge } from '@atoms';
import { ItemStatus } from '../types/item';

interface ItemStatusBadgeProps {
  status: ItemStatus;
}

export const ItemStatusBadge = ({ status }: ItemStatusBadgeProps) => {
  const t = useTranslations('ItemStatus');

  return <Badge variant={status.variant}>{t(status.value)}</Badge>;
};
