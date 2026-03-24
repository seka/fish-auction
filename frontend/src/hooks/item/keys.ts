/**
 * アイテム関連のクエリキー
 */
export const itemKeys = {
  // Public
  publicAll: ['items'] as const,
  byAuction: (auctionId: number) => ['auctions', 'detail', auctionId, 'items'] as const,
} as const;
