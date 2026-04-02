import type { Metadata } from 'next';
import { Noto_Sans_JP } from 'next/font/google';
import '@/src/core/styles/globals.css';

import { AuthorizablePublicNavbar } from '@/src/features/login';
import { MainLayoutTemplate } from '@templates';
import { PushInitializer } from '@bootstraps';

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
import { QueryClientProvider, ToastProvider } from '@bootstraps';

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
          <QueryClientProvider>
            <ToastProvider>
              <MainLayoutTemplate navbar={<AuthorizablePublicNavbar />}>
                <PushInitializer />
                {children}
              </MainLayoutTemplate>
            </ToastProvider>
          </QueryClientProvider>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
