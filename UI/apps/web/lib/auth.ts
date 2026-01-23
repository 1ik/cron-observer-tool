/**
 * Server-side authentication utilities for NextAuth.js
 * 
 * IMPORTANT: These functions are for SERVER COMPONENTS and SERVER ACTIONS only!
 * 
 * For CLIENT COMPONENTS, use:
 * - import { signIn, signOut, useSession } from 'next-auth/react'
 * 
 * For SERVER COMPONENTS, use:
 * - import { auth, signIn, signOut, getSession } from '@/lib/auth'
 */

import { authConfig } from '@/auth.config'
import NextAuth from 'next-auth'

const nextAuth = NextAuth(authConfig)

export async function auth() {
  return await nextAuth.auth()
}

export async function signIn(...args: Parameters<typeof nextAuth.signIn>) {
  return await nextAuth.signIn(...args)
}

export async function signOut(...args: Parameters<typeof nextAuth.signOut>) {
  return await nextAuth.signOut(...args)
}

export async function getSession() {
  return await auth()
}

