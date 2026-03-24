'use client';

import { Box, Card, Text, Stack, EmptyState } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';

interface Invoice {
  buyerName: string;
  buyerId: number;
  totalAmount: number;
}

interface InvoiceListProps {
  invoices: Invoice[];
}

export const InvoiceList = ({ invoices }: InvoiceListProps) => {
  const t = useTranslations();

  return (
    <Stack spacing="4">
      <Text fontSize="xl" fontWeight="bold" className={css({ color: 'gray.800' })}>
        {t('Public.MyPage.invoices')}
      </Text>
      {invoices.length === 0 ? (
        <EmptyState
          message={t('Public.MyPage.no_invoices')}
          icon={
            <span role="img" aria-label="invoice">
              🧾
            </span>
          }
        />
      ) : (
        invoices.map((invoice, index) => (
          <Card
            key={index}
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
        ))
      )}
    </Stack>
  );
};
