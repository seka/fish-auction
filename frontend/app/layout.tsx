import type { Metadata } from 'next';
import { Noto_Sans_JP } from 'next/font/google';
import './globals.css';
import Providers from './providers';

import { PublicNavbar } from '@organisms';
import { MainLayoutTemplate } from '@templates';
import { PushInitializer } from '@functionals';

const notoSansJP = Noto_Sans_JP({
  subsets: ['latin'],
  variable: '--font-noto-sans-jp',
});

export const metadata: Metadata = {
  title: 'FISHING AUCTION',
  description: 'FISHING AUCTION System',
  icons: {
    icon: '/favicon.ico',
  },
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
      <body className={notoSansJP.className}>
        <NextIntlClientProvider messages={messages}>
          <Providers>
            <MainLayoutTemplate navbar={<PublicNavbar />}>
              <PushInitializer />
              {children}
            </MainLayoutTemplate>
          </Providers>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
