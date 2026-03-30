import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useInvoices } from '../queries/useInvoices';
import { InvoiceItem } from '../types';

export const useInvoiceManagement = () => {
  const t = useTranslations();
  const { invoices, isLoading } = useInvoices();
  const [selectedInvoice, setSelectedInvoice] = useState<InvoiceItem | null>(null);

  return {
    state: {
      invoices,
      isLoading,
      selectedInvoice,
    },
    actions: {
      setSelectedInvoice,
    },
    t,
  };
};
