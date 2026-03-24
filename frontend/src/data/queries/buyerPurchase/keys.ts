/**
 * バイヤー用購入履歴関連のクエリキー
 */
export const buyerPurchaseKeys = {
  meAll: () => ['buyer', 'me', 'purchases'] as const,
} as const;
