/**
 * 認証関連のクエリキー
 */
export const authKeys = {
  me: () => ['buyer', 'me'] as const,
} as const;
