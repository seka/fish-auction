'use client';

import { Box, Button, HStack, Text } from '@atoms';
import { css } from 'styled-system/css';

interface Buyer {
  id?: number;
  name: string;
}

interface BuyerListItemProps {
  buyer: Buyer;
  isDeleting: boolean;
  onDelete: (id: number) => void;
  t: (key: string) => string;
}

export const BuyerListItem = ({ buyer, isDeleting, onDelete, t }: BuyerListItemProps) => {
  return (
    <Box as="li" p="6" _hover={{ bg: 'gray.50' }} transition="colors">
      <HStack justify="between" align="center">
        <Box>
          <Text
            as="h3"
            fontSize="lg"
            fontWeight="bold"
            className={css({ color: 'green.900' })}
          >
            {buyer.name}
          </Text>
          <Text fontSize="sm" className={css({ color: 'gray.600' })} mt="1">
            ID: {buyer.id}
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
          onClick={() => buyer.id && onDelete(buyer.id)}
          disabled={isDeleting}
        >
          {t('Common.delete')}
        </Button>
      </HStack>
    </Box>
  );
};
