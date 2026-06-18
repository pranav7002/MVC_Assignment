'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '../stores/authStore'
import { protectedFetch } from '../utils/api'

export default function LoginPage() {
    const setToken = useAuthStore((state) => state.setToken)

    const router = useRouter()

    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')
    const [loading, setLoading] = useState(false)

    return (
        <>
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
                        setToken(data.data)
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
            >
                <input
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder="Username"
                />

                <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Password"
                />

                <button type="submit" disabled={loading}>
                    Login
                </button>
            </form>
            <div> {error ? <p>{error}</p> : null} </div>
        </>
    )
}
