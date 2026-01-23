import type { NextAuthConfig } from 'next-auth'
import Google from 'next-auth/providers/google'

export const authConfig: NextAuthConfig = {
  providers: [
    Google({
      clientId: process.env.GOOGLE_CLIENT_ID!,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
    }),
  ],
  callbacks: {
    async jwt({ token, user, account }) {
      // Initial sign in - store user info in token
      if (account && user) {
        token.sub = user.id || user.email || ''
        token.email = user.email || ''
        token.name = user.name || ''
        token.picture = user.image || ''
        // Store Google access token for API calls if needed
        token.accessToken = account.access_token
      }

      return token
    },
    async session({ session, token }) {
      // Add user info to session
      if (token) {
        session.user = {
          ...session.user,
          id: token.sub || '',
          email: (token.email as string) || '',
          name: (token.name as string) || '',
          image: (token.picture as string) || '',
        }
        // The JWT token itself will be used for backend auth
        // We'll get it from NextAuth's internal token
        ;(session as any).accessToken = token.accessToken as string
      }
      return session
    },
  },
  pages: {
    signIn: '/auth/signin',
  },
  session: {
    strategy: 'jwt',
  },
  secret: process.env.NEXTAUTH_SECRET,
} satisfies NextAuthConfig

