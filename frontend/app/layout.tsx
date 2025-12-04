import type { Metadata } from "next";
import { Noto_Sans_JP } from "next/font/google";
import "./globals.css";
import Providers from "./providers";
import { css } from "styled-system/css";

const notoSansJP = Noto_Sans_JP({
  subsets: ["latin"],
  variable: "--font-noto-sans-jp",
});

export const metadata: Metadata = {
  title: "漁港のせりシステム",
  description: "Fish Auction System",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <body
        className={`${notoSansJP.variable} ${css({ fontFamily: 'sans', bg: 'gray.50', color: 'gray.900', antialiased: true })}`}
      >
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
