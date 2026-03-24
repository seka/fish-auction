/**
 * オークション関連のクエリキー
 */
export const auctionKeys = {
  all: ['auctions'] as const,
  lists: () => [...auctionKeys.all, 'list'] as const,
  list: (filters: unknown) => [...auctionKeys.lists(), { filters }] as const,
  details: () => [...auctionKeys.all, 'detail'] as const,
  detail: (id: number) => [...auctionKeys.details(), id] as const,
  items: (id: number) => [...auctionKeys.detail(id), 'items'] as const,
} as const;
