import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import { cookieHeader } from '@/src/core/api/client';
import { checkAdminSession, checkBuyerSession } from '@/src/data/queries/auth/useQuery';

export async function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const cookie = cookieHeader(request.headers.get('cookie'));

  if (pathname.startsWith('/admin') || pathname.startsWith('/invoice')) {
    if (!request.cookies.get('admin_session')) {
      return NextResponse.redirect(new URL('/login', request.url));
    }
    if (!await checkAdminSession(cookie)) {
      return NextResponse.redirect(new URL('/login', request.url));
    }
  }

  if (pathname.startsWith('/auction')) {
    if (!request.cookies.get('buyer_session')) {
      return NextResponse.redirect(new URL('/login/buyer', request.url));
    }
    if (!await checkBuyerSession(cookie)) {
      return NextResponse.redirect(new URL('/login/buyer', request.url));
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/admin/:path*', '/invoice/:path*', '/auction/:path*'],
};
