import { authConfig } from '@/auth.config'
import { SignJWT } from 'jose'
import NextAuth from 'next-auth'
import { NextResponse } from 'next/server'

const { auth } = NextAuth(authConfig)

export async function GET() {
  try {
    const session = await auth()

    if (!session || !session.user) {
      return NextResponse.json(
        { error: 'Unauthorized' },
        { status: 401 }
      )
    }

    // Create a JWT token from session data that matches what the backend expects
    const secret = new TextEncoder().encode(process.env.NEXTAUTH_SECRET!)
    
    const token = await new SignJWT({
      sub: session.user.id || session.user.email || '',
      email: session.user.email || '',
      name: session.user.name || '',
    })
      .setProtectedHeader({ alg: 'HS256' })
      .setIssuedAt()
      .setExpirationTime('24h')
      .sign(secret)

    return NextResponse.json({
      token,
    })
  } catch (error) {
    console.error('Error getting token:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}

