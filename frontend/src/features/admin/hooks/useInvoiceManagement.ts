import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useInvoiceQuery } from '@/src/data/queries/adminInvoice/useQuery';
import { InvoiceItem } from '@/src/models';

export const useInvoiceManagement = () => {
  const t = useTranslations();
  const { invoices, isLoading } = useInvoiceQuery();
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
