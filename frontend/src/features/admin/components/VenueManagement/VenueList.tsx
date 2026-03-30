import { Box, Stack, Text } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import { Venue } from '../../types';
import { VenueListCell } from './VenueListCell';

interface VenueListProps {
  venues: Venue[];
  isLoading: boolean;
  onEdit: (venue: Venue) => void;
  onDelete: (id: number) => void;
}

export const VenueList = ({ venues, isLoading, onEdit, onDelete }: VenueListProps) => {
  const t = useTranslations();

  if (isLoading) {
    return (
      <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>
        {t('Common.loading')}
      </Box>
    );
  }

  if (venues.length === 0) {
    return (
      <Box p="12" textAlign="center">
        <Text color="muted">{t('Common.no_data')}</Text>
      </Box>
    );
  }

  return (
    <Stack as="ul" spacing="0" divideY="1px" divideColor="gray.200">
      {venues.map((venue) => (
        <VenueListCell key={venue.id} venue={venue} onEdit={onEdit} onDelete={onDelete} t={t} />
      ))}
    </Stack>
  );
};
