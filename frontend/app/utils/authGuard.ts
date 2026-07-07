'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '../stores/authStore'

// Protected pages — redirect to /login if no token
export function useRequireAuth() {
    const token = useAuthStore((state) => state.token)
    const router = useRouter()

    useEffect(() => {
        if (!token) router.replace('/login')
    }, [token, router])
}

// Guest pages — redirect to /village if token exists
export function useRedirectIfAuth() {
    const token = useAuthStore((state) => state.token)
    const router = useRouter()

    useEffect(() => {
        if (token) router.replace('/village')
    }, [token, router])
}
