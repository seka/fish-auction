'use client';

import { Box, Button, HStack, Stack, Text } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import { Venue } from '@/src/models';

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
        <Box as="li" key={venue.id} p="6" _hover={{ bg: 'gray.50' }} transition="colors">
          <HStack justify="between" align="start">
            <Box>
              <Text
                as="h3"
                fontSize="lg"
                fontWeight="bold"
                className={css({ color: 'indigo.900' })}
              >
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
      ))}
    </Stack>
  );
};
