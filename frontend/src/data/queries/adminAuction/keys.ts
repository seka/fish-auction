/**
 * 管理用オークション関連のクエリキー
 */
export const adminAuctionKeys = {
  all: ['admin', 'auctions'] as const,
  list: (filters: unknown) => ['admin', 'auctions', 'list', { filters }] as const,
  detail: (id: number) => ['admin', 'auctions', 'detail', id] as const,
  items: (id: number) => ['admin', 'auctions', 'detail', id, 'items'] as const,
} as const;
