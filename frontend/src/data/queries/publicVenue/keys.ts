/**
 * 会場関連のクエリキー
 */
export const venueKeys = {
  // Public
  publicAll: ['venues'] as const,
  publicDetail: (id: number) => ['venues', 'detail', id] as const,
} as const;
