// このファイルをモジュールとして認識させ、declare global を有効にする
export {};

declare const CookieHeaderBrand: unique symbol;

declare global {
  type CookieHeader = string & { readonly [CookieHeaderBrand]: true };
}
