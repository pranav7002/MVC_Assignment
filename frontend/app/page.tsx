'use client'

import { useRouter } from 'next/navigation'

export default function Home() {
    const router = useRouter()

    return (
        <div className="auth-page">
            <div className="auth-card" style={{ textAlign: 'center' }}>
                <h1 style={{ fontSize: '24px', marginBottom: '4px', letterSpacing: '4px' }}>Vanguard</h1>
                <p className="subtitle">village combat</p>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
                    <button className="btn btn-green" style={{ width: '100%', justifyContent: 'center' }} onClick={() => router.push('/login')}>
                        Login
                    </button>
                    <button className="btn btn-blue" style={{ width: '100%', justifyContent: 'center' }} onClick={() => router.push('/register')}>
                        Register
                    </button>
                </div>
            </div>
        </div>
    )
}
