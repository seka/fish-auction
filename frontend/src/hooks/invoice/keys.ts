/**
 * 請求書関連のクエリキー
 */
export const invoiceKeys = {
  // Public
  publicAll: ['invoices'] as const,

  // Buyer
  meAll: () => ['buyer', 'me', 'invoices'] as const,
} as const;
