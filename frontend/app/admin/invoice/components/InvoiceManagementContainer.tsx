'use client';

import { useInvoiceManagement } from '@/src/features/admin/states/useInvoiceManagement';
import { InvoiceList, InvoiceDetailModal } from '@/src/features/admin/components/InvoiceManagement';
import { Box, Card, Text } from '@atoms';
import { css } from 'styled-system/css';

export const InvoiceManagementContainer = () => {
  const { state, actions, t } = useInvoiceManagement();

  return (
    <Box maxW="5xl" mx="auto" p="6">
      <Text
        as="h1"
        variant="h2"
        className={css({ color: 'gray.800' })}
        mb="8"
        pb="4"
        borderBottom="1px solid"
        borderColor="gray.200"
      >
        {t('Admin.Invoice.title')}
      </Text>

      <Card padding="none" overflow="hidden">
        <InvoiceList
          invoices={state.invoices}
          isLoading={state.isLoading}
          onSelect={actions.setSelectedInvoice}
        />
      </Card>

      {/* Detail Modal */}
      {state.selectedInvoice && (
        <InvoiceDetailModal
          invoice={state.selectedInvoice}
          onClose={() => actions.setSelectedInvoice(null)}
        />
      )}
    </Box>
  );
};
