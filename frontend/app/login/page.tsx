'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '../stores/authStore'
import { protectedFetch } from '../utils/api'
import { useRedirectIfAuth } from '../utils/authGuard'

export default function LoginPage() {
    useRedirectIfAuth()
    const setAuth = useAuthStore((state) => state.setAuth)

    const router = useRouter()

    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')
    const [loading, setLoading] = useState(false)

    return (
        <div className="auth-page">
            <div className="auth-card">
                <h1>Vanguard</h1>
                <p className="subtitle">sign in to your village</p>
                <form
                    onSubmit={async (e) => {
                        e.preventDefault()
                        setLoading(true)
                        try {
                            const res = await protectedFetch('/api/auth/login', 'POST', {
                                username,
                                password,
                            })
                            if (!res.ok) {
                                const err = await res.json()
                                setError(err.error)
                                return
                            }
                            const data = await res.json()
                            setAuth(data.data, username)
                            router.push('/village')
                        } catch (err) {
                            setError('Unable to connect to server')
                            setUsername('')
                            setPassword('')
                            console.error(err)
                        } finally {
                            setLoading(false)
                        }
                    }}
                    style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}
                >
                    <input
                        className="input"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        placeholder="Username"
                    />
                    <input
                        className="input"
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        placeholder="Password"
                    />
                    <button className="btn btn-green" type="submit" disabled={loading} style={{ width: '100%', justifyContent: 'center' }}>
                        {loading ? 'Signing in...' : 'Login'}
                    </button>
                </form>
                {error && <p style={{ color: 'var(--danger)', fontSize: 'clamp(11px, 1vw, 13px)', marginTop: '12px', textAlign: 'center' }}>{error}</p>}
                <p style={{ fontSize: 'clamp(11px, 1vw, 13px)', color: 'var(--text-muted)', marginTop: '16px', textAlign: 'center' }}>
                    No account?{' '}
                    <span onClick={() => router.push('/register')} style={{ color: 'var(--accent-blue)', cursor: 'pointer' }}>Register</span>
                </p>
            </div>
        </div>
    )
}
