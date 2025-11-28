import { apiClient } from '@/src/core/api/client';
import { InvoiceItem } from '@/src/models';

export const getInvoices = async (): Promise<InvoiceItem[]> => {
    return apiClient.get<InvoiceItem[]>('/api/invoices');
};
