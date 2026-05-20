declare const CookieHeaderBrand: unique symbol;

declare global {
  type CookieHeader = string & { readonly [CookieHeaderBrand]: true };
}
