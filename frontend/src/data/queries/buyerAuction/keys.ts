/**
 * バイヤー用オークション関連のクエリキー
 */
export const buyerAuctionKeys = {
  meAll: () => ['buyer', 'me', 'auctions'] as const,
} as const;
