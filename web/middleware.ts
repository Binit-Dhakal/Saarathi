// middleware.ts
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  var hostname = request.headers.get("host");
  const pathname = request.nextUrl.pathname;


  // Handle redirects for the 'driver' subdomain
  if (hostname?.startsWith('driver.saarathi.com')) {
    if (pathname === '/sign-in') {
      const redirectUrl = new URL('/driver/sign-in', request.url);
      return NextResponse.redirect(redirectUrl);
    }
    if (pathname === '/sign-up') {
      const redirectUrl = new URL('/driver/sign-up', request.url);
      return NextResponse.redirect(redirectUrl);
    }
  }

  // Handle redirects for the 'saarathi' subdomain (riders)
  if (hostname?.startsWith('saarathi.com')) {
    if (pathname === '/sign-in') {
      const redirectUrl = new URL('/rider/sign-in', request.url);
      return NextResponse.redirect(redirectUrl);
    }
    if (pathname === '/sign-up') {
      const redirectUrl = new URL('/rider/sign-up', request.url);
      return NextResponse.redirect(redirectUrl);
    }
  }

  // For all other requests, let the request proceed as normal without a redirect
  return NextResponse.next();
}

export const config = {
  matcher: [
    // Apply middleware to specific paths that might need a redirect
    '/sign-in',
    '/sign-up',
    // You can also include other paths as needed
  ],
};

