'use client';

import { Box, Button, HStack, Text } from '@atoms';
import { css } from 'styled-system/css';
import { Venue } from '../../types';

interface VenueListCellProps {
  venue: Venue;
  onEdit: (venue: Venue) => void;
  onDelete: (id: number) => void;
  t: (key: string) => string;
}

export const VenueListCell = ({ venue, onEdit, onDelete, t }: VenueListCellProps) => {
  return (
    <Box as="li" key={venue.id} p="6" _hover={{ bg: 'gray.50' }} transition="colors">
      <HStack justify="between" align="start">
        <Box>
          <Text as="h3" fontSize="lg" fontWeight="bold" className={css({ color: 'indigo.900' })}>
            {venue.name}
          </Text>
          {venue.location && (
            <Text
              fontSize="sm"
              className={css({ color: 'gray.700' })}
              mt="1"
              display="flex"
              alignItems="center"
            >
              <span className={css({ mr: '2' })}>📍</span>
              {venue.location}
            </Text>
          )}
          {venue.description && (
            <Text fontSize="sm" className={css({ color: 'gray.600' })} mt="2">
              {venue.description}
            </Text>
          )}
        </Box>
        <HStack spacing="2">
          <Button size="sm" variant="outline" onClick={() => onEdit(venue)}>
            {t('Common.edit')}
          </Button>
          <Button
            size="sm"
            className={css({
              bg: 'red.50',
              color: 'red.600',
              _hover: { bg: 'red.100' },
            })}
            onClick={() => onDelete(venue.id)}
          >
            {t('Common.delete')}
          </Button>
        </HStack>
      </HStack>
    </Box>
  );
};
