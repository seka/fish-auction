import { Box, Stack, Text } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import { BuyerListItem } from './BuyerListItem';

interface Buyer {
  id?: number;
  name: string;
}

interface BuyerListProps {
  buyers: Buyer[];
  isLoading: boolean;
  isDeleting: boolean;
  onDelete: (id: number) => void;
}

export const BuyerList = ({ buyers, isLoading, isDeleting, onDelete }: BuyerListProps) => {
  const t = useTranslations();

  if (isLoading) {
    return (
      <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>
        {t('Common.loading')}
      </Box>
    );
  }

  if (buyers.length === 0) {
    return (
      <Box p="12" textAlign="center">
        <Text color="muted">{t('Common.no_data')}</Text>
      </Box>
    );
  }

  return (
    <Stack as="ul" spacing="0" divideY="1px" divideColor="gray.200">
      {buyers.map((buyer) => (
        <BuyerListItem
          key={buyer.id}
          buyer={buyer}
          isDeleting={isDeleting}
          onDelete={onDelete}
          t={t}
        />
      ))}
    </Stack>
  );
};
