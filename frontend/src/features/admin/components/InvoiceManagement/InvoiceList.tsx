'use client';

import { Table, Thead, Tbody, Tr, Th, Td } from '@molecules';
import { Box, Text, EmptyState } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';

interface InvoiceListProps {
  invoices: any[];
  isLoading: boolean;
  onSelect: (invoice: any) => void;
}

export const InvoiceList = ({ invoices, isLoading, onSelect }: InvoiceListProps) => {
  const t = useTranslations();

  if (isLoading) {
    return (
      <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>
        {t('Common.loading')}
      </Box>
    );
  }

  if (invoices.length === 0) {
    return (
      <EmptyState
        message={t('Admin.Invoice.no_data')}
        icon={
          <span role="img" aria-label="invoice">
            🧾
          </span>
        }
      />
    );
  }

  return (
    <Table>
      <Thead>
        <Tr className={css({ cursor: 'default', _hover: { bg: 'gray.50' } })}>
          <Th>{t('Admin.Invoice.buyer_id')}</Th>
          <Th>{t('Admin.Invoice.buyer_name')}</Th>
          <Th className={css({ textAlign: 'right' })}>{t('Admin.Invoice.total_amount')}</Th>
        </Tr>
      </Thead>
      <Tbody>
        {invoices.map((invoice) => (
          <Tr key={invoice.buyerId} onClick={() => onSelect(invoice)}>
            <Td className={css({ fontSize: 'sm', color: 'gray.500', fontFamily: 'mono' })}>
              {invoice.buyerId}
            </Td>
            <Td className={css({ fontSize: 'sm', fontWeight: 'bold', color: 'gray.900' })}>
              {invoice.buyerName}
            </Td>
            <Td
              className={css({
                textAlign: 'right',
                fontWeight: 'bold',
                color: 'indigo.700',
                fontSize: 'lg',
              })}
            >
              ¥{invoice.totalAmount.toLocaleString()}
            </Td>
          </Tr>
        ))}
      </Tbody>
    </Table>
  );
};
