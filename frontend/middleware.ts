import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import { cookieHeader } from '@/src/core/api/client';
import {
  checkAdminSession,
  checkBuyerSession,
  adminSessionCookie,
  buyerSessionCookie,
  SessionResult,
} from '@/src/middleware/session';

export async function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const cookie = cookieHeader(request.headers.get('cookie'));

  if (pathname.startsWith('/admin') || pathname.startsWith('/invoice')) {
    if (!request.cookies.get(adminSessionCookie)) {
      return NextResponse.redirect(new URL('/login', request.url));
    }
    const adminResult: SessionResult = await checkAdminSession(cookie);
    if (adminResult === 'error') return NextResponse.redirect(new URL('/error', request.url));
    if (adminResult === 'invalid') return NextResponse.redirect(new URL('/login', request.url));
  }

  if (pathname.startsWith('/auction')) {
    if (!request.cookies.get(buyerSessionCookie)) {
      return NextResponse.redirect(new URL('/login/buyer', request.url));
    }
    const buyerResult: SessionResult = await checkBuyerSession(cookie);
    if (buyerResult === 'error') return NextResponse.redirect(new URL('/error', request.url));
    if (buyerResult === 'invalid')
      return NextResponse.redirect(new URL('/login/buyer', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/admin/:path*', '/invoice/:path*', '/auction/:path*'],
};
