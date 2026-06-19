'use client'

import { useEffect, useRef, useState } from 'react'
import { useRouter, useParams } from 'next/navigation'
import { useAuthStore } from '../../stores/authStore'
import { useBattleStore } from '../../stores/battleStore'
import { protectedFetch } from '../../utils/api'
import sprites from '../../village/spriteLoader'

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

interface TroopState {
    id: number
    name: string
    pos: { x: number; y: number }
    hp: number
    dead: boolean
}

interface BuildingHP {
    hp: number
    max_hp: number
    dead: boolean
}

interface BattleState {
    destruction_pct: number
    stars: number
    troops: TroopState[]
    buildings: { id: number; name: string; hp: number; max_hp: number; dead: boolean }[]
}

interface TrainedTroop {
    troop_name: string
    quantity: number
}

const GRID_SIZE = 20
const CELL_SIZE = 45

function getUserIDFromToken(token: string): string | null {
    try {
        const payload = token.split('.')[1]
        const decoded = JSON.parse(atob(payload))
        return decoded.user_id || null
    } catch {
        return null
    }
}

export default function BattlePage() {
    const token = useAuthStore((state) => state.token)
    const router = useRouter()
    const params = useParams()
    const defenderID = params.defenderID as string

    const storeBuildings = useBattleStore((state) => state.buildings)

    const [buildings, setBuildings] = useState<Building[]>(storeBuildings)
    const [buildingHPs, setBuildingHPs] = useState<Map<number, BuildingHP>>(
        new Map(),
    )
    const [troops, setTroops] = useState<TroopState[]>([])
    const [destructionPct, setDestructionPct] = useState(0)
    const [starCount, setStarCount] = useState(0)
    const [battleOver, setBattleOver] = useState(false)

    const [trainedTroops, setTrainedTroops] = useState<TrainedTroop[]>([])
    const [selectedTroop, setSelectedTroop] = useState<string | null>(null)

    const canvasRef = useRef<HTMLCanvasElement>(null)
    const grassTileRef = useRef<HTMLImageElement | null>(null)
    const wsRef = useRef<WebSocket | null>(null)

    async function loadTrainedTroops() {
        try {
            const res = await protectedFetch('/api/troops', 'GET')
            const data = await res.json()
            setTrainedTroops(data.data || [])
        } catch (err) {
            console.error(err)
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
        if (!token) return

        loadTrainedTroops()

        const userID = getUserIDFromToken(token)
        if (!userID) return

        let wasOpen = false

        const ws = new WebSocket(
            `ws://localhost:8080/api/battle/ws/${defenderID}?user_id=${userID}`,
        )
        wsRef.current = ws

        ws.onopen = () => {
            wasOpen = true
        }

        ws.onmessage = (event) => {
            const state: BattleState = JSON.parse(event.data)
            console.log('ws tick', { troops: state.troops.length, buildings: state.buildings.length, destruction: state.destruction_pct })

            setTroops(state.troops)
            setDestructionPct(state.destruction_pct)
            setStarCount(state.stars)

            const hpMap = new Map<number, BuildingHP>()
            for (const b of state.buildings) {
                hpMap.set(b.id, {
                    hp: b.hp,
                    max_hp: b.max_hp,
                    dead: b.dead,
                })
            }
            setBuildingHPs(hpMap)
        }

        ws.onclose = () => {
            if (wasOpen) {
                setBattleOver(true)
            }
        }

        return () => {
            ws.close()
        }
    }, [token])

    // canvas draw
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

        // draw buildings
        ;[...buildings]
            .sort((a, b) => a.pos_y + a.size - (b.pos_y + b.size))
            .forEach((b) => {
                const dx = b.pos_x * CELL_SIZE
                const footprint = b.size * CELL_SIZE
                const bottomY = (b.pos_y + b.size) * CELL_SIZE

                const hpInfo = buildingHPs.get(b.id)
                const isDead = hpInfo ? hpInfo.dead : false

                if (isDead) {
                    ctx.globalAlpha = 0.3
                }

                const sprite = sprites[`${b.building_name}_${b.level}`]
                if (sprite && sprite.complete && sprite.naturalWidth > 0) {
                    const spriteW = footprint
                    const spriteH =
                        (sprite.naturalHeight / sprite.naturalWidth) * spriteW
                    const dy = bottomY - spriteH
                    ctx.drawImage(sprite, dx, dy, spriteW, spriteH)
                } else {
                    ctx.fillStyle = 'rgba(100, 100, 100, 0.9)'
                    ctx.fillRect(
                        dx,
                        b.pos_y * CELL_SIZE,
                        footprint,
                        footprint,
                    )
                    ctx.fillStyle = 'white'
                    ctx.font = '12px Arial'
                    ctx.fillText(
                        b.building_name,
                        dx + 4,
                        b.pos_y * CELL_SIZE + 18,
                    )
                }

                ctx.globalAlpha = 1

                // HP bar
                if (hpInfo && !isDead) {
                    const barY = b.pos_y * CELL_SIZE - 6
                    const barW = footprint
                    const barH = 4
                    const hpRatio = hpInfo.hp / hpInfo.max_hp

                    // background
                    ctx.fillStyle = '#333'
                    ctx.fillRect(dx, barY, barW, barH)

                    // fill
                    if (hpRatio > 0.6) {
                        ctx.fillStyle = '#4caf50'
                    } else if (hpRatio > 0.3) {
                        ctx.fillStyle = '#ffeb3b'
                    } else {
                        ctx.fillStyle = '#f44336'
                    }
                    ctx.fillRect(dx, barY, barW * hpRatio, barH)
                }
            })

        // draw troops
        for (const t of troops) {
            if (t.dead) continue

            const tx = t.pos.x * CELL_SIZE
            const ty = t.pos.y * CELL_SIZE

            const sprite = sprites[t.name]
            if (sprite && sprite.complete && sprite.naturalWidth > 0) {
                ctx.drawImage(sprite, tx, ty, CELL_SIZE, CELL_SIZE)
            } else {
                ctx.fillStyle = 'blue'
                ctx.fillRect(tx + 8, ty + 8, CELL_SIZE - 16, CELL_SIZE - 16)
            }
        }
    }, [buildings, buildingHPs, troops])

    function handleCanvasClick(e: React.MouseEvent<HTMLCanvasElement>) {
        console.log('click', { selectedTroop, wsOpen: wsRef.current?.readyState, trainedTroops })
        if (!selectedTroop || !wsRef.current) return
        if (wsRef.current.readyState !== WebSocket.OPEN) return

        const scaleX =
            (GRID_SIZE * CELL_SIZE) / e.currentTarget.clientWidth
        const scaleY =
            (GRID_SIZE * CELL_SIZE) / e.currentTarget.clientHeight

        const cellX = Math.floor(
            (e.nativeEvent.offsetX * scaleX) / CELL_SIZE,
        )
        const cellY = Math.floor(
            (e.nativeEvent.offsetY * scaleY) / CELL_SIZE,
        )

        console.log('deploying', { cellX, cellY, selectedTroop })

        const troopData = trainedTroops.find(
            (t) => t.troop_name === selectedTroop,
        )
        if (!troopData || troopData.quantity <= 0) return

        wsRef.current.send(
            JSON.stringify({
                name: selectedTroop,
                x: cellX,
                y: cellY,
            }),
        )

        setTrainedTroops((prev) =>
            prev.map((t) =>
                t.troop_name === selectedTroop
                    ? { ...t, quantity: t.quantity - 1 }
                    : t,
            ),
        )
    }

    return (
        <div style={{ padding: '20px' }}>
            {/* Header */}
            <div
                style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '16px',
                    marginBottom: '16px',
                }}
            >
                <h1>Battle</h1>
                <div>Destruction: {destructionPct}%</div>
                <div>Stars: {'⭐'.repeat(starCount)}</div>
            </div>

            {/* Canvas */}
            <div style={{ position: 'relative', display: 'inline-block' }}>
                <canvas
                    ref={canvasRef}
                    width={GRID_SIZE * CELL_SIZE}
                    height={GRID_SIZE * CELL_SIZE}
                    onClick={handleCanvasClick}
                    style={{
                        border: '1px solid black',
                        cursor: selectedTroop ? 'crosshair' : 'default',
                    }}
                />

                {/* Result Overlay */}
                {battleOver && (
                    <div
                        style={{
                            position: 'absolute',
                            top: 0,
                            left: 0,
                            width: '100%',
                            height: '100%',
                            backgroundColor: 'rgba(0, 0, 0, 0.7)',
                            display: 'flex',
                            flexDirection: 'column',
                            alignItems: 'center',
                            justifyContent: 'center',
                            color: 'white',
                        }}
                    >
                        <h2
                            style={{
                                fontSize: '48px',
                                marginBottom: '16px',
                            }}
                        >
                            {'⭐'.repeat(starCount)}
                        </h2>
                        <p style={{ fontSize: '24px' }}>
                            Destruction: {destructionPct}%
                        </p>
                        <button
                            onClick={() => {
                                useBattleStore.getState().clear()
                                router.push('/village')
                            }}
                            style={{
                                marginTop: '20px',
                                padding: '10px 20px',
                                cursor: 'pointer',
                                fontSize: '16px',
                            }}
                        >
                            Back to Village
                        </button>
                    </div>
                )}
            </div>

            {/* Deploy Bar */}
            {!battleOver && (
                <div
                    style={{
                        display: 'flex',
                        gap: '8px',
                        marginTop: '12px',
                        flexWrap: 'wrap',
                    }}
                >
                    {trainedTroops.map((t) => (
                        <button
                            key={t.troop_name}
                            onClick={() => setSelectedTroop(t.troop_name)}
                            disabled={t.quantity <= 0}
                            style={{
                                padding: '8px 12px',
                                cursor:
                                    t.quantity > 0
                                        ? 'pointer'
                                        : 'not-allowed',
                                border:
                                    selectedTroop === t.troop_name
                                        ? '2px solid blue'
                                        : '1px solid #ccc',
                            }}
                        >
                            <img
                                src={`/sprites/troops/${t.troop_name.toLowerCase()}.png`}
                                alt={t.troop_name}
                                style={{
                                    width: '32px',
                                    height: '32px',
                                    imageRendering: 'pixelated',
                                    display: 'block',
                                    margin: '0 auto 4px',
                                }}
                            />
                            {t.troop_name} x{t.quantity}
                        </button>
                    ))}
                </div>
            )}
        </div>
    )
}
