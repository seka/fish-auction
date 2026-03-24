import { useState } from 'react';
import { useInvoiceQuery } from '@/src/hooks/invoice/useQuery';
import { InvoiceItem } from '@/src/models';

export const useInvoicePage = () => {
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
  };
};
