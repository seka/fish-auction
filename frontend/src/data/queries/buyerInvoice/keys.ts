export const buyerInvoiceKeys = {
  all: ['buyerInvoice'] as const,
  meAll: () => [...buyerInvoiceKeys.all, 'me'] as const,
};
