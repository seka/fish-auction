/**
 * 中買人関連のクエリキー
 */
export const buyerKeys = {
  // Buyer My Page (API: /api/buyer/me/...)
  mePurchases: () => ['buyer', 'me', 'purchases'] as const,
  meAuctions: () => ['buyer', 'me', 'auctions'] as const,
} as const;
