'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useSelector } from 'react-redux'
import type { RootState } from '@/store'

export default function AuthGuard({ children }: { children: React.ReactNode }) {
  const router = useRouter()
  const token = useSelector((s: RootState) => s.auth.token)

  useEffect(() => {
    if (!token) router.replace('/login')
  }, [token, router])

  if (!token) return null
  return <>{children}</>
}
