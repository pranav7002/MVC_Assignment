'use client'

import { useEffect, useRef, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '../stores/authStore'
import { useBattleStore } from '../stores/battleStore'
import { protectedFetch } from '../utils/api'
import { useRequireAuth } from '../utils/authGuard'
import sprites from '../village/spriteLoader'

interface Building {
    id: number
    building_name: string
    level: number
    pos_x: number
    pos_y: number
    size: number
    hp: number
    is_upgrading: boolean
}

const GRID_SIZE = 20
const CELL_SIZE = 45

export default function MatchmakingPage() {
    useRequireAuth()
    const token = useAuthStore((state) => state.token)
    const router = useRouter()

    const [defenderID, setDefenderID] = useState<string | null>(null)
    const [buildings, setBuildings] = useState<Building[]>([])
    const [loading, setLoading] = useState(false)
    const [troopCount, setTroopCount] = useState<number | null>(null)

    const canvasRef = useRef<HTMLCanvasElement>(null)
    const grassTileRef = useRef<HTMLImageElement | null>(null)

    async function findMatch() {
        setLoading(true)
        try {
            const res = await protectedFetch('/api/battle/match', 'GET')
            if (!res.ok) throw new Error((await res.json()).error)
            const data = await res.json()
            setDefenderID(data.data.defenders_id)
            setBuildings(data.data.buildings)
        } catch (error: any) {
            alert(error.message || 'No opponents found')
        }
        setLoading(false)
    }

    async function checkTroops() {
        try {
            const res = await protectedFetch('/api/troops', 'GET')
            const data = await res.json()
            const troops = data.data || []
            const total = troops.reduce((sum: number, t: any) => sum + t.quantity, 0)
            setTroopCount(total)
        } catch {
            setTroopCount(0)
        }
    }

    useEffect(() => {
        const img = new Image()
        img.src = '/sprites/grass.png'
        img.onload = () => {
            grassTileRef.current = img
            setBuildings((prev) => [...prev])
        }
    }, [])

    useEffect(() => {
        if (token) {
            findMatch()
            checkTroops()
        }
    }, [token])

    useEffect(() => {
        const canvas = canvasRef.current
        if (!canvas) return

        const ctx = canvas.getContext('2d')
        if (!ctx) return

        ctx.imageSmoothingEnabled = false
        ctx.clearRect(0, 0, canvas.width, canvas.height)

        const grassTile = grassTileRef.current
        if (grassTile) {
            for (let x = 0; x < GRID_SIZE; x++) {
                for (let y = 0; y < GRID_SIZE; y++) {
                    ctx.drawImage(
                        grassTile,
                        x * CELL_SIZE,
                        y * CELL_SIZE,
                        CELL_SIZE,
                        CELL_SIZE,
                    )
                }
            }
        }

        ;[...buildings]
            .sort((a, b) => (a.pos_y + a.size) - (b.pos_y + b.size))
            .forEach((b) => {
                const dx = b.pos_x * CELL_SIZE
                const footprint = b.size * CELL_SIZE
                const bottomY = (b.pos_y + b.size) * CELL_SIZE

                const sprite = sprites[`${b.building_name}_${b.level}`]
                if (sprite && sprite.complete && sprite.naturalWidth > 0) {
                    const spriteW = footprint
                    const spriteH = (sprite.naturalHeight / sprite.naturalWidth) * spriteW
                    const dy = bottomY - spriteH
                    ctx.drawImage(sprite, dx, dy, spriteW, spriteH)
                } else {
                    ctx.fillStyle = 'rgba(100, 100, 100, 0.9)'
                    ctx.fillRect(dx, b.pos_y * CELL_SIZE, footprint, footprint)
                    ctx.fillStyle = 'white'
                    ctx.font = '12px Arial'
                    ctx.fillText(b.building_name, dx + 4, b.pos_y * CELL_SIZE + 18)
                }
            })
    }, [buildings])

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
            {/* ── Top Bar ── */}
            <div className="topbar">
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                    <button className="btn" onClick={() => router.push('/village')}>← Village</button>
                    <span className="topbar-title">Matchmaking</span>
                </div>
                <div className="topbar-nav">
                    <button className="btn" onClick={findMatch} disabled={loading}>
                        {loading ? '...' : '↻ Skip'}
                    </button>
                    <button
                        className="btn btn-green"
                        onClick={() => {
                            if (defenderID && troopCount && troopCount > 0) {
                                useBattleStore.getState().setBattle(defenderID, buildings)
                                router.push(`/battle/${defenderID}`)
                            }
                        }}
                        disabled={!defenderID || !troopCount || troopCount <= 0}
                    >
                        {troopCount === 0 ? 'No Troops!' : '⚔ Attack'}
                    </button>
                    <button className="btn btn-compact" suppressHydrationWarning onClick={() => { useAuthStore.getState().logout(); router.push('/login') }}>{useAuthStore.getState().username} · Logout</button>
                </div>
            </div>

            {/* ── Canvas ── */}
            <div className="canvas-wrap" style={{ flex: 1 }}>
                <canvas
                    ref={canvasRef}
                    width={GRID_SIZE * CELL_SIZE}
                    height={GRID_SIZE * CELL_SIZE}
                />
            </div>
        </div>
    )
}
