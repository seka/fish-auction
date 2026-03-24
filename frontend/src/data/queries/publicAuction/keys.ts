/**
 * オークション関連のクエリキー
 */
export const auctionKeys = {
  // Public (API: /api/auctions)
  publicAll: ['auctions'] as const,
  publicList: (filters: unknown) => ['auctions', 'list', { filters }] as const,
  publicDetail: (id: number) => ['auctions', 'detail', id] as const,
  publicItems: (id: number) => ['auctions', 'detail', id, 'items'] as const,
} as const;
