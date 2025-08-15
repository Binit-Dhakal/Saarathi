import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  const hostname = request.headers.get("host");
  const pathname = request.nextUrl.pathname;
  const accessToken = request.cookies.get("accessToken");

  const publicPaths = [
    '/driver/sign-in',
    '/driver/sign-up',
    '/rider/sign-in',
    '/rider/sign-up'
  ];

  const isPublicPath = publicPaths.includes(pathname);

  // Condition 1: User is NOT logged in and is trying to access a protected page.
  if (!accessToken && !isPublicPath) {
    if (hostname?.startsWith('driver.saarathi.com')) {
      const url = request.nextUrl.clone();
      url.pathname = '/driver/sign-in';
      return NextResponse.redirect(url);
    }
    // Assume all other hosts are for riders by default
    const url = request.nextUrl.clone();
    url.pathname = '/rider/sign-in';
    return NextResponse.redirect(url);
  }

  // Condition 2: User IS logged in and is trying to access a public sign-in/sign-up page.
  if (accessToken && isPublicPath) {
    const url = request.nextUrl.clone();
    url.pathname = '/'; // Redirect to the protected root path
    return NextResponse.redirect(url);
  }

  // fallback
  return NextResponse.next();
}

export const config = {
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico|.*\\.png$|.*\\.svg$).*)',
  ],
};
