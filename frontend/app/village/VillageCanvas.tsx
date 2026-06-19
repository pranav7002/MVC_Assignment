'use client'

import { useEffect, useRef, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '../stores/authStore'
import { protectedFetch } from '../utils/api'
import sprites from './spriteLoader'

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

interface ShopBuilding {
    building_type: string
    building_name: string
    max_level: number
    max_built: number
    cost: number
    cost_type: string
    size: number
}

const GRID_SIZE = 20
const CELL_SIZE = 45

export default function VillageCanvas() {
    const token = useAuthStore((state) => state.token)
    const router = useRouter()

    const [buildings, setBuildings] = useState<Building[]>([])
    const [shopBuildings, setShopBuildings] = useState<ShopBuilding[]>([])

    // building not on the grid, due to buying/moving
    const [selectedBuilding, setSelectedBuilding] = useState<{
        Name: string
        Level: number
        Size: number
    } | null>(null)

    const [hoverCell, setHoverCell] = useState<{
        x: number
        y: number
    } | null>(null)

    // clicked on a building, but it is still at its place
    const [activeBuilding, setActiveBuilding] = useState<Building | null>(null)

    const [upgradeInfo, setUpgradeInfo] = useState<{
        is_max_level: boolean
        upgrade_cost: number
        upgrade_cost_type: string
        next_max_hp: number
        upgrade_duration_sec: number
    } | null>(null)

    // i have kept this to store the building id for the api request, rest similar 
    // to selected building, just only for buildings that are already placed and being moved 
    const [movingBuildingId, setMovingBuildingId] = useState<number | null>(
        null,
    )

    const [village, setVillage] = useState({
        gold: 0,
        elixir: 0,
    })

    const canvasRef = useRef<HTMLCanvasElement>(null)
    const grassTileRef = useRef<HTMLImageElement | null>(null)

    async function loadBuildings() {
        try {
            const res = await protectedFetch('/api/buildings', 'GET')
            const data = await res.json()
            setBuildings(data.data)
        } catch (err) {
            console.error(err)
        }
    }

    async function loadShopBuildings() {
        try {
            const res = await protectedFetch('/api/shop/buildings', 'GET')
            const data = await res.json()
            setShopBuildings(data.data)
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
            loadBuildings()
            loadVillage()
            loadShopBuildings()
        }
    }, [token])

    useEffect(() => {
        const img = new Image()
        img.src = '/sprites/grass.png'
        img.onload = () => {
            grassTileRef.current = img
            setBuildings((prev) => [...prev])
        }
    }, [])

    function isValidPlacement(
        posX: number,
        posY: number,
        size: number,
        movingId?: number,
    ) {
        if (
            posX < 0 ||
            posY < 0 ||
            posX + size > GRID_SIZE ||
            posY + size > GRID_SIZE
        ) {
            return false
        }

        for (const b of buildings) {
            if (movingId && b.id === movingId) continue

            if (
                posX < b.pos_x + b.size &&
                posX + size > b.pos_x &&
                posY < b.pos_y + b.size &&
                posY + size > b.pos_y
            ) {
                return false
            }
        }

        return true
    }

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

        [...buildings]
            .sort((a, b) => (a.pos_y + a.size) - (b.pos_y + b.size))
            .forEach((b) => {
            if (movingBuildingId === b.id) return

            const dx = b.pos_x * CELL_SIZE
            const footprint = b.size * CELL_SIZE
            const bottomY = (b.pos_y + b.size) * CELL_SIZE

            if (activeBuilding?.id === b.id) {
                ctx.fillStyle = 'rgba(255, 255, 0, 0.4)'
                ctx.fillRect(dx, b.pos_y * CELL_SIZE, footprint, footprint)
            }

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

        if (selectedBuilding && hoverCell) {
            const valid = isValidPlacement(
                hoverCell.x,
                hoverCell.y,
                selectedBuilding.Size,
                movingBuildingId || undefined,
            )
            ctx.globalAlpha = 0.5
            ctx.fillStyle = valid ? 'green' : 'red'
            ctx.fillRect(
                hoverCell.x * CELL_SIZE,
                hoverCell.y * CELL_SIZE,
                selectedBuilding.Size * CELL_SIZE,
                selectedBuilding.Size * CELL_SIZE,
            )
            ctx.globalAlpha = 1
        }
    }, [
        buildings,
        selectedBuilding,
        hoverCell,
        activeBuilding,
        movingBuildingId,
    ])

    function handleMouseMove(e: React.MouseEvent<HTMLCanvasElement>) {
        const scaleX = (GRID_SIZE * CELL_SIZE) / e.currentTarget.clientWidth
        const scaleY = (GRID_SIZE * CELL_SIZE) / e.currentTarget.clientHeight

        const x = Math.floor((e.nativeEvent.offsetX * scaleX) / CELL_SIZE)
        const y = Math.floor((e.nativeEvent.offsetY * scaleY) / CELL_SIZE)

        setHoverCell({ x, y })
    }

    async function handleCanvasClick() {
        if (!hoverCell) return

        if (selectedBuilding) {
            if (
                !isValidPlacement(
                    hoverCell.x,
                    hoverCell.y,
                    selectedBuilding.Size,
                    movingBuildingId || undefined,
                )
            )
                return

            try {
                if (movingBuildingId) {
                    const res = await protectedFetch(
                        `/api/buildings/${movingBuildingId}/move`,
                        'PUT',
                        {
                            pos_x: hoverCell.x,
                            pos_y: hoverCell.y,
                        },
                    )
                    if (!res.ok) {
                        throw new Error((await res.json()).error)
                    }
                } else {
                    const info = shopBuildings.find(
                        (b) => b.building_name === selectedBuilding.Name,
                    )
                    if (!info) return

                    const res = await protectedFetch('/api/buildings', 'POST', {
                        building_type: info.building_type,
                        building_name: selectedBuilding.Name,
                        pos_x: hoverCell.x,
                        pos_y: hoverCell.y,
                    })
                    if (!res.ok) {
                        throw new Error((await res.json()).error)
                    }
                }

                await loadBuildings()
                await loadVillage()
                await loadShopBuildings()
                deselectEverything()
            } catch (error: any) {
                alert(error.message || 'Action failed')
            }
            return
        }

        const clicked = buildings.find(
            (b) =>
                hoverCell.x >= b.pos_x &&
                hoverCell.x < b.pos_x + b.size &&
                hoverCell.y >= b.pos_y &&
                hoverCell.y < b.pos_y + b.size,
        )
        setActiveBuilding(clicked || null)

        if (clicked) {
            try {
                const res = await protectedFetch(`/api/buildings/${clicked.id}/upgrade-info`, 'GET')
                if (res.ok) {
                    const rawText = await res.text()
                    console.log("Raw successful response:", rawText)
                    const data = JSON.parse(rawText)
                    setUpgradeInfo(data.data)
                } else {
                    const rawText = await res.text()
                    console.error("Backend returned error:", res.status, rawText)
                    setUpgradeInfo(null)
                }
            } catch (err) {
                console.error("Network or parsing error:", err)
                setUpgradeInfo(null)
            }
        } else {
            setUpgradeInfo(null)
        }
    }

    async function handleUpgrade() {
        if (!activeBuilding) return
        try {
            if (activeBuilding.building_name === 'Town Hall') {
                const res = await protectedFetch(
                    `/api/village/upgrade-th`,
                    'PUT',
                )
                if (!res.ok) throw new Error((await res.json()).error)
                await loadVillage()
                await loadBuildings()
                await loadShopBuildings()
                setActiveBuilding(null)
            } else {
                const res = await protectedFetch(
                    `/api/buildings/${activeBuilding.id}/upgrade`,
                    'PUT',
                )
                if (!res.ok) throw new Error((await res.json()).error)
                await loadVillage()
                await loadBuildings()
                setActiveBuilding(null)
            }
        } catch (error: any) {
            alert(error.message || 'Upgrade failed')
        }
    }

    function deselectEverything() {
        setSelectedBuilding(null)
        setHoverCell(null)
        setMovingBuildingId(null)
        setActiveBuilding(null)
        setUpgradeInfo(null)
    }

    async function handleCollect(resourceType: string) {
        try {
            const res = await protectedFetch('/api/economy/collect', 'POST', {
                resource_type: `${resourceType}`
            })
            if (!res.ok) throw new Error((await res.json()).error)
            await loadVillage()
        } catch (error: any) {
            alert(error.message || 'Collection failed')
        }
    }

    return (
        <div>
            {/* Top Bar */}
            <div
                style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '16px',
                    padding: '12px',
                    marginBottom: '16px',
                }}
            >
                <h1>My Village</h1>

                <div>Gold: {village?.gold ?? 0}</div>

                <div>Elixir: {village?.elixir ?? 0}</div>

                <button
                    onClick={async () => handleCollect('gold')}
                    style={{
                        marginLeft: 'auto',
                        cursor: 'pointer',
                    }}
                >
                    Collect Gold
                </button>
                <button
                    onClick={async () => handleCollect('elixir')}
                    style={{
                        marginLeft: 'auto',
                        cursor: 'pointer',
                    }}
                >
                    Collect Elixir
                </button>
                <button
                    onClick={() => router.push('/troops')}
                    style={{
                        marginLeft: 'auto',
                        cursor: 'pointer',
                    }}
                >
                    Troops
                </button>
                <button
                    onClick={() => router.push('/matchmaking')}
                    style={{
                        cursor: 'pointer',
                    }}
                >
                    Find Match
                </button>
            </div>

            <div style={{ display: 'flex' }}>
                {/* Sidebar */}
                <div
                    style={{
                        width: '250px',
                        padding: '12px',
                        marginRight: '20px',
                    }}
                >
                    {activeBuilding && !selectedBuilding && (
                        <div style={{ marginBottom: '16px' }}>
                            <h3>{activeBuilding.building_name}</h3>

                            <p>Level: {activeBuilding.level}</p>

                            <p>HP: {activeBuilding.hp}</p>

                            <button
                                style={{
                                    marginRight: '8px',
                                    cursor: 'pointer',
                                }}
                                onClick={() => {
                                    setMovingBuildingId(activeBuilding.id)
                                    setSelectedBuilding({
                                        Name: activeBuilding.building_name,
                                        Level: activeBuilding.level,
                                        Size: activeBuilding.size,
                                    })
                                }}
                            >
                                Move
                            </button>

                            {upgradeInfo ? (
                                upgradeInfo.is_max_level ? (
                                    <p style={{ color: 'red', marginTop: '8px' }}>Max Level for current Town Hall</p>
                                ) : (
                                    <>
                                        <p style={{ marginTop: '8px' }}>Upgrade Cost: {upgradeInfo.upgrade_cost} {upgradeInfo.upgrade_cost_type}</p>
                                        <button
                                            style={{
                                                cursor: (upgradeInfo.upgrade_cost_type === 'gold' ? village.gold < upgradeInfo.upgrade_cost : village.elixir < upgradeInfo.upgrade_cost) ? 'not-allowed' : 'pointer',
                                            }}
                                            disabled={upgradeInfo.upgrade_cost_type === 'gold' ? village.gold < upgradeInfo.upgrade_cost : village.elixir < upgradeInfo.upgrade_cost}
                                            onClick={handleUpgrade}
                                        >
                                            Upgrade 
                                        </button>
                                    </>
                                )
                            ) : (
                                <p style={{ marginTop: '8px' }}>Loading upgrade info...</p>
                            )}
                        </div>
                    )}

                    <h2>Shop</h2>

                    {shopBuildings.map((info) => {
                        const countBuilt = buildings.filter(
                            (b) => b.building_name === info.building_name,
                        ).length

                        const isMaxedOut = countBuilt >= info.max_built

                        return (
                            <button
                                key={info.building_name}
                                disabled={isMaxedOut}
                                onClick={() => {
                                    setSelectedBuilding({
                                        Name: info.building_name,
                                        Level: 1,
                                        Size: info.size,
                                    })
                                    setMovingBuildingId(null)
                                    setActiveBuilding(null)
                                }}
                                style={{
                                    display: 'block',
                                    width: '100%',
                                    marginBottom: '8px',
                                    padding: '6px',
                                    cursor: isMaxedOut
                                        ? 'not-allowed'
                                        : 'pointer',
                                }}
                            >
                                {`${info.building_name} ${info.cost_type} ${info.cost}`}
                            </button>
                        )
                    })}

                    {selectedBuilding && (
                        <button
                            onClick={deselectEverything}
                            style={{
                                marginTop: '12px',
                                cursor: 'pointer',
                            }}
                        >
                            {movingBuildingId
                                ? 'Cancel Move'
                                : 'Cancel Placement'}
                        </button>
                    )}
                </div>

                {/* Canvas */}
                <div>
                    <canvas
                        ref={canvasRef}
                        width={GRID_SIZE * CELL_SIZE}
                        height={GRID_SIZE * CELL_SIZE}
                        onMouseMove={handleMouseMove}
                        onClick={handleCanvasClick}
                        onMouseLeave={() => setHoverCell(null)}
                        style={{
                            border: '1px solid black',
                            cursor: selectedBuilding ? 'crosshair' : 'default',
                        }}
                    />
                </div>
            </div>
        </div>
    )
}
