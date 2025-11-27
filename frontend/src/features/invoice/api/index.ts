import { apiClient } from '@/src/shared/api/client';
import { InvoiceItem } from '@/src/shared/models';

export const getInvoices = async (): Promise<InvoiceItem[]> => {
    return apiClient.get<InvoiceItem[]>('/api/invoices');
};
