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
        <div style={{ padding: '20px' }}>
            <div
                style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '16px',
                    marginBottom: '16px',
                }}
            >
                <button
                    onClick={() => router.push('/village')}
                    style={{ cursor: 'pointer' }}
                >
                    ← Back to Village
                </button>

                <h1>Troops</h1>

                <div style={{ marginLeft: 'auto' }}>
                    Gold: {village.gold} | Elixir: {village.elixir}
                </div>
            </div>

            <div
                style={{
                    display: 'flex',
                    flexWrap: 'wrap',
                    gap: '12px',
                }}
            >
                {troopConfigs.map((troop) => {
                    const count = getTrainedCount(troop.name)

                    return (
                        <div
                            key={troop.name}
                            style={{
                                border: '1px solid #ccc',
                                padding: '12px',
                                width: '180px',
                                textAlign: 'center',
                            }}
                        >
                            <img
                                src={`/sprites/troops/${troop.name.toLowerCase()}.png`}
                                alt={troop.name}
                                style={{
                                    width: '64px',
                                    height: '64px',
                                    imageRendering: 'pixelated',
                                }}
                            />

                            <div style={{ marginTop: '8px' }}>
                                <strong>{troop.name}</strong>
                            </div>

                            <div style={{ fontSize: '12px', marginTop: '4px' }}>
                                DPS: {troop.dps} | HP: {troop.health} | Range:{' '}
                                {troop.range}
                            </div>

                            <div style={{ fontSize: '12px' }}>
                                Cost: {troop.training_cost} elixir
                            </div>

                            <div style={{ fontSize: '12px' }}>
                                Space: {troop.housing_space}
                            </div>

                            <div
                                style={{
                                    marginTop: '8px',
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'center',
                                    gap: '8px',
                                }}
                            >
                                <button
                                    onClick={() => handleTrain(troop.name)}
                                    style={{ cursor: 'pointer' }}
                                >
                                    Train
                                </button>

                                {count > 0 && (
                                    <>
                                        <span>x{count}</span>
                                        <button
                                            onClick={() =>
                                                handleDelete(troop.name)
                                            }
                                            style={{
                                                cursor: 'pointer',
                                                color: 'red',
                                            }}
                                        >
                                            Delete
                                        </button>
                                    </>
                                )}
                            </div>
                        </div>
                    )
                })}
            </div>
        </div>
    )
}
