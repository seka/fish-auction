// このファイルをモジュールとして認識させ、declare global を有効にする
export {};

declare const CookieHeaderBrand: unique symbol;

declare global {
  type CookieHeader = string & { readonly [CookieHeaderBrand]: true };
  // Next.js Edge Runtime global — present only when running on the edge, undefined otherwise
  const EdgeRuntime: string | undefined;
}
