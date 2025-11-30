import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
    const adminSession = request.cookies.get('admin_session');
    const buyerSession = request.cookies.get('buyer_session');
    const { pathname } = request.nextUrl;

    // Protect admin routes
    if (pathname.startsWith('/admin') || pathname.startsWith('/invoice')) {
        if (!adminSession || adminSession.value !== 'authenticated') {
            return NextResponse.redirect(new URL('/login', request.url));
        }
    }

    // Protect auction routes
    if (pathname.startsWith('/auction')) {
        if (!buyerSession || buyerSession.value !== 'authenticated') {
            return NextResponse.redirect(new URL('/login/buyer', request.url));
        }
    }

    return NextResponse.next();
}

export const config = {
    matcher: ['/admin/:path*', '/invoice/:path*', '/auction/:path*'],
};
