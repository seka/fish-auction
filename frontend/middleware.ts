import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

const TIMEOUT_MS = 3000;

async function isSessionValid(endpoint: string, cookieHeader: string): Promise<boolean> {
  const baseUrl = process.env.API_BASE_URL;
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), TIMEOUT_MS);

  try {
    const res = await fetch(`${baseUrl}${endpoint}`, {
      headers: { Cookie: cookieHeader },
      signal: controller.signal,
    });
    return res.status !== 401;
  } catch {
    // バックエンド障害やタイムアウト時はフロント全断を防ぐため fail-open とする
    return true;
  } finally {
    clearTimeout(timer);
  }
}

export async function middleware(request: NextRequest) {
  const adminSession = request.cookies.get('admin_session');
  const buyerSession = request.cookies.get('buyer_session');
  const { pathname } = request.nextUrl;

  if (pathname.startsWith('/admin') || pathname.startsWith('/invoice')) {
    if (!adminSession) {
      return NextResponse.redirect(new URL('/login', request.url));
    }
    const valid = await isSessionValid('/api/admin/me', `admin_session=${adminSession.value}`);
    if (!valid) {
      return NextResponse.redirect(new URL('/login', request.url));
    }
  }

  if (pathname.startsWith('/auction')) {
    if (!buyerSession) {
      return NextResponse.redirect(new URL('/login/buyer', request.url));
    }
    const valid = await isSessionValid('/api/buyer/me', `buyer_session=${buyerSession.value}`);
    if (!valid) {
      return NextResponse.redirect(new URL('/login/buyer', request.url));
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/admin/:path*', '/invoice/:path*', '/auction/:path*'],
};
