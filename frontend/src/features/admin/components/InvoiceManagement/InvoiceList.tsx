'use client';

import { Table, Thead, Tbody, Tr, Th } from '@molecules';
import { Box, EmptyState } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import { InvoiceItem } from '@/src/models';
import { InvoiceListCell } from './InvoiceListCell';

interface InvoiceListProps {
  invoices: InvoiceItem[];
  isLoading: boolean;
  onSelect: (invoice: InvoiceItem) => void;
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
        <Tr className={css({ cursor: 'default' })}>
          <Th>{t('Admin.Invoice.buyer_id')}</Th>
          <Th>{t('Admin.Invoice.buyer_name')}</Th>
          <Th className={css({ textAlign: 'right' })}>{t('Admin.Invoice.total_amount')}</Th>
        </Tr>
      </Thead>
      <Tbody>
        {invoices.map((invoice) => (
          <InvoiceListCell key={invoice.buyerId} invoice={invoice} onSelect={onSelect} />
        ))}
      </Tbody>
    </Table>
  );
};
