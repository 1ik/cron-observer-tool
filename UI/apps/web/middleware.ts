import { authConfig } from '@/auth.config'
import NextAuth from 'next-auth'
import type { NextRequest } from 'next/server'
import { NextResponse } from 'next/server'

const { auth } = NextAuth(authConfig)

export async function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl

  // Public routes that don't require authentication
  const publicRoutes = ['/auth/signin', '/api/auth']
  const isPublicRoute = publicRoutes.some((route) => pathname.startsWith(route))

  if (isPublicRoute) {
    return NextResponse.next()
  }

  // Check authentication for protected routes
  const session = await auth()

  if (!session) {
    // Redirect to sign in if not authenticated
    const signInUrl = new URL('/auth/signin', request.url)
    signInUrl.searchParams.set('callbackUrl', pathname)
    return NextResponse.redirect(signInUrl)
  }

  return NextResponse.next()
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     */
    '/((?!_next/static|_next/image|favicon.ico).*)',
  ],
}

