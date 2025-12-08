import type { Metadata } from "next";
import { Noto_Sans_JP } from "next/font/google";
import "./globals.css";
import Providers from "./providers";
import { css } from "styled-system/css";

import { PublicNavbar } from "./_components/PublicNavbar";

const notoSansJP = Noto_Sans_JP({
  subsets: ["latin"],
  variable: "--font-noto-sans-jp",
});

export const metadata: Metadata = {
  title: "漁港のせりシステム",
  description: "Fish Auction System",
};

import { NextIntlClientProvider } from 'next-intl';
import { getLocale, getMessages } from 'next-intl/server';

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const locale = await getLocale();
  const messages = await getMessages();

  return (
    <html lang={locale}>
      <body
        className={`${notoSansJP.variable} ${css({ fontFamily: 'sans', bg: 'gray.50', color: 'gray.900' })}`}
      >
        <NextIntlClientProvider messages={messages}>
          <Providers>
            <PublicNavbar />
            {children}
          </Providers>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
