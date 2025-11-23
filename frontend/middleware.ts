import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
    // Check if the path starts with /admin
    if (request.nextUrl.pathname.startsWith('/admin') || request.nextUrl.pathname.startsWith('/invoice')) {
        const adminSession = request.cookies.get('admin_session');

        if (!adminSession || adminSession.value !== 'authenticated') {
            // Redirect to login page if not authenticated
            return NextResponse.redirect(new URL('/login', request.url));
        }
    }

    return NextResponse.next();
}

export const config = {
    matcher: ['/admin/:path*', '/invoice/:path*'],
};
