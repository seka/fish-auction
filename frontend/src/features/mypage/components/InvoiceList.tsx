import { Text, Stack, EmptyState } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import { InvoiceListCell } from './InvoiceListCell';

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
        invoices.map((invoice, index) => <InvoiceListCell key={index} invoice={invoice} />)
      )}
    </Stack>
  );
};
