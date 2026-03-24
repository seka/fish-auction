/**
 * アイテム関連のクエリキー
 */
export const itemKeys = {
  all: ['items'] as const,
  byAuction: (auctionId: number) => [...itemKeys.all, 'auction', auctionId] as const,
} as const;
