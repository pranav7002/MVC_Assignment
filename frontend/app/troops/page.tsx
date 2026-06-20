'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '../stores/authStore'
import { protectedFetch } from '../utils/api'

interface TroopConfig {
    name: string
    dps: number
    health: number
    range: number
    housing_space: number
    training_cost: number
}

interface TrainedTroop {
    troop_name: string
    quantity: number
}

export default function TroopsPage() {
    const token = useAuthStore((state) => state.token)
    const router = useRouter()

    const [troopConfigs, setTroopConfigs] = useState<TroopConfig[]>([])
    const [trainedTroops, setTrainedTroops] = useState<TrainedTroop[]>([])
    const [village, setVillage] = useState({ gold: 0, elixir: 0 })

    async function loadTroopConfigs() {
        try {
            const res = await protectedFetch('/api/shop/troops', 'GET')
            const data = await res.json()
            setTroopConfigs(data.data)
        } catch (err) {
            console.error(err)
        }
    }

    async function loadTrainedTroops() {
        try {
            const res = await protectedFetch('/api/troops', 'GET')
            const data = await res.json()
            setTrainedTroops(data.data || [])
        } catch (err) {
            console.error(err)
        }
    }

    async function loadVillage() {
        if (!token) return
        try {
            const res = await protectedFetch('/api/village', 'GET')
            if (res.ok) {
                const data = await res.json()
                setVillage(data.data)
            }
        } catch (err) {
            console.error(err)
        }
    }

    useEffect(() => {
        if (token) {
            loadTroopConfigs()
            loadTrainedTroops()
            loadVillage()
        }
    }, [token])

    async function handleTrain(troopName: string) {
        try {
            const res = await protectedFetch('/api/troops/train', 'POST', {
                troop_name: troopName,
                quantity: 1,
            })
            if (!res.ok) throw new Error((await res.json()).error)
            await loadTrainedTroops()
            await loadVillage()
        } catch (error: any) {
            alert(error.message || 'Training failed')
        }
    }

    async function handleDelete(troopName: string) {
        try {
            const res = await protectedFetch(
                `/api/troops/${troopName}`,
                'DELETE',
            )
            if (!res.ok) throw new Error((await res.json()).error)
            await loadTrainedTroops()
            await loadVillage()
        } catch (error: any) {
            alert(error.message || 'Delete failed')
        }
    }

    function getTrainedCount(troopName: string): number {
        const found = trainedTroops.find((t) => t.troop_name === troopName)
        return found ? found.quantity : 0
    }

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
            {/* ── Top Bar ── */}
            <div className="topbar">
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                    <button className="btn" onClick={() => router.push('/village')}>← Village</button>
                    <span className="topbar-title">Troops</span>
                </div>
                <div className="topbar-nav">
                    <span className="resource-pill gold">⛏ {village.gold}</span>
                    <span className="resource-pill elixir">🧪 {village.elixir}</span>
                </div>
            </div>

            {/* ── Content ── */}
            <div style={{ flex: 1, padding: '0 12px 12px 12px', paddingTop: '96px', overflow: 'auto', display: 'flex' }}>
                <div style={{ display: 'flex', gap: '18px', flex: 1 }}>
                    {troopConfigs.map((troop) => {
                        const count = getTrainedCount(troop.name)

                        return (
                            <div key={troop.name} className="troop-card" style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'flex-start', padding: '220px 16px 24px 16px' }}>
                                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', width: '100%', height: '40%' }}>
                                    <img
                                        src={`/sprites/troops/${troop.name.toLowerCase()}.png`}
                                        alt={troop.name}
                                        style={{ width: '80%', height: '80%', objectFit: 'contain', imageRendering: 'pixelated' }}
                                    />
                                </div>
                                <div style={{ fontSize: 'clamp(16px, 1.5vw, 22px)', fontWeight: 'bold', marginBottom: '10px', marginTop: '12px' }}>{troop.name}</div>
                                <div className="troop-stat" style={{ fontSize: 'clamp(12px, 1.1vw, 15px)' }}>DPS: {troop.dps}  ·  HP: {troop.health}</div>
                                <div className="troop-stat" style={{ fontSize: 'clamp(12px, 1.1vw, 15px)' }}>Range: {troop.range}  ·  Space: {troop.housing_space}</div>
                                <div className="troop-stat" style={{ fontSize: 'clamp(12px, 1.1vw, 15px)', color: 'var(--elixir)' }}>{troop.training_cost} elixir</div>

                                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '8px', paddingTop: '14px' }}>
                                    <button className="btn btn-green" onClick={() => handleTrain(troop.name)}>
                                        + Train
                                    </button>
                                    {count > 0 && (
                                        <>
                                            <span style={{ color: 'var(--text-secondary)' }}>x{count}</span>
                                            <button
                                                className="btn btn-danger"
                                                onClick={() => handleDelete(troop.name)}
                                            >
                                                −
                                            </button>
                                        </>
                                    )}
                                </div>
                            </div>
                        )
                    })}
                </div>
            </div>
        </div>
    )
}
