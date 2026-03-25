'use client';

import { Box, Button, HStack, Stack, Text } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';

interface Fisherman {
  id: number;
  name: string;
}

interface FishermanListProps {
  fishermen: Fisherman[];
  isLoading: boolean;
  isDeleting: boolean;
  onDelete: (id: number) => void;
}

export const FishermanList = ({
  fishermen,
  isLoading,
  isDeleting,
  onDelete,
}: FishermanListProps) => {
  const t = useTranslations();

  if (isLoading) {
    return (
      <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>
        {t('Common.loading')}
      </Box>
    );
  }

  if (fishermen.length === 0) {
    return (
      <Box p="12" textAlign="center">
        <Text color="gray.500">{t('Common.no_data')}</Text>
      </Box>
    );
  }

  return (
    <Stack as="ul" spacing="0" divideY="1px" divideColor="gray.200">
      {fishermen.map((fisherman) => (
        <Box
          as="li"
          key={fisherman.id}
          p="6"
          _hover={{ bg: 'gray.50' }}
          transition="colors"
        >
          <HStack justify="between" align="center">
            <Box>
              <Text
                as="h3"
                fontSize="lg"
                fontWeight="bold"
                className={css({ color: 'indigo.900' })}
              >
                {fisherman.name}
              </Text>
              <Text fontSize="sm" className={css({ color: 'gray.600' })} mt="1">
                ID: {fisherman.id}
              </Text>
            </Box>
            <Button
              variant="outline"
              size="sm"
              className={css({
                color: 'red.500',
                borderColor: 'red.200',
                _hover: { bg: 'red.50', borderColor: 'red.500' },
              })}
              onClick={() => fisherman.id && onDelete(fisherman.id)}
              disabled={isDeleting}
            >
              {t('Common.delete')}
            </Button>
          </HStack>
        </Box>
      ))}
    </Stack>
  );
};
