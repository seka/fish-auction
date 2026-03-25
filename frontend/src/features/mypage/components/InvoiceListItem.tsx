'use client';

import { Box, Card, Text } from '@atoms';
import { css } from 'styled-system/css';

interface Invoice {
  buyerName: string;
  buyerId: number;
  totalAmount: number;
}

interface InvoiceListItemProps {
  invoice: Invoice;
  t: (key: string) => string;
}

export const InvoiceListItem = ({ invoice, t }: InvoiceListItemProps) => {
  return (
    <Card
      padding="lg"
      className={css({
        _hover: { shadow: 'md' },
        transition: 'all 0.2s',
        borderWidth: '1px',
        borderColor: 'gray.200',
        bg: 'white',
      })}
    >
      <Box display="flex" justifyContent="space-between" alignItems="center">
        <Box>
          <Text fontSize="lg" fontWeight="bold">
            {invoice.buyerName}
          </Text>
          <Text fontSize="sm" className={css({ color: 'gray.500' })}>
            Buyer ID: {invoice.buyerId}
          </Text>
        </Box>
        <Box textAlign="right">
          <Text fontSize="2xl" fontWeight="bold" className={css({ color: 'green.600' })}>
            ¥{invoice.totalAmount.toLocaleString()}
          </Text>
        </Box>
      </Box>
    </Card>
  );
};
