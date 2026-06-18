'use client'

import { useEffect, useRef, useState } from 'react'
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

const GRID_SIZE = 20
const CELL_SIZE = 40

const BUILDING_INFO: Record<string, { type: string; size: number }> = {
    'Town Hall': { type: 'town_hall', size: 4 },
    'Cannon': { type: 'defense', size: 3 },
    'Archer Tower': { type: 'defense', size: 3 },
    'Mortar': { type: 'defense', size: 3 },
    'Gold Mine': { type: 'resource', size: 3 },
    'Elixir Collector': { type: 'resource', size: 3 },
    'Gold Storage': { type: 'storage', size: 3 },
    'Elixir Storage': { type: 'storage', size: 3 },
    'Training Grounds': { type: 'training_grounds', size: 3 },
}

export default function VillageCanvas() {
    const token = useAuthStore((state) => state.token)

    const [buildings, setBuildings] = useState<Building[]>([])
    const [selectedBuilding, setSelectedBuilding] = useState<string | null>(null)
    const [hoverCell, setHoverCell] = useState<{ x: number; y: number } | null>(null)
    const [activeBuilding, setActiveBuilding] = useState<Building | null>(null)
    const [movingBuildingId, setMovingBuildingId] = useState<number | null>(null)

    const canvasRef = useRef<HTMLCanvasElement>(null)
    const grassTileRef = useRef<HTMLImageElement | null>(null)

    // --- Data Fetching ---

    async function loadBuildings() {
        try {
            const res = await protectedFetch('/api/buildings', 'GET')
            const data = await res.json()
            setBuildings(data.data || [])
        } catch (err) {
            console.error(err)
        }
    }

    useEffect(() => {
        if (token) loadBuildings()
    }, [token])

    // --- Grass Tile ---

    useEffect(() => {
        const img = new Image()
        img.src = '/sprites/grass.png'
        img.onload = () => {
            grassTileRef.current = img
            setBuildings((prev) => [...prev])
        }
    }, [])

    // --- Validation ---

    function isValidPlacement(posX: number, posY: number, size: number, ignoreId?: number) {
        if (posX < 0 || posY < 0 || posX + size > GRID_SIZE || posY + size > GRID_SIZE) {
            return false
        }

        for (const b of buildings) {
            if (ignoreId && b.id === ignoreId) continue

            if (posX < b.pos_x + b.size && posX + size > b.pos_x &&
                posY < b.pos_y + b.size && posY + size > b.pos_y) {
                return false
            }
        }

        return true
    }

    // --- Canvas Drawing ---

    useEffect(() => {
        const canvas = canvasRef.current
        if (!canvas) return

        const ctx = canvas.getContext('2d')
        if (!ctx) return

        ctx.clearRect(0, 0, canvas.width, canvas.height)

        // Grass
        const grassTile = grassTileRef.current
        if (grassTile) {
            for (let x = 0; x < GRID_SIZE; x++) {
                for (let y = 0; y < GRID_SIZE; y++) {
                    ctx.drawImage(grassTile, x * CELL_SIZE, y * CELL_SIZE, CELL_SIZE, CELL_SIZE)
                }
            }
        }

        // Buildings
        buildings.forEach((b) => {
            if (movingBuildingId === b.id) return

            const px = b.pos_x * CELL_SIZE
            const py = b.pos_y * CELL_SIZE
            const pw = b.size * CELL_SIZE

            // Selection highlight
            if (activeBuilding?.id === b.id) {
                ctx.fillStyle = 'rgba(255, 255, 0, 0.4)'
                ctx.fillRect(px, py, pw, pw)
            }

            // Sprite or fallback
            const sprite = sprites[b.building_name]
            if (sprite && sprite.complete) {
                ctx.drawImage(sprite, px, py, pw, pw)
            } else {
                ctx.fillStyle = 'rgba(100, 100, 100, 0.9)'
                ctx.fillRect(px, py, pw, pw)
                ctx.fillStyle = 'white'
                ctx.font = '12px Arial'
                ctx.fillText(b.building_name, px + 4, py + 18)
            }
        })

        // Hover preview (placing or moving)
        if (selectedBuilding && hoverCell) {
            const info = BUILDING_INFO[selectedBuilding]
            if (info) {
                const valid = isValidPlacement(hoverCell.x, hoverCell.y, info.size, movingBuildingId || undefined)
                ctx.globalAlpha = 0.5
                ctx.fillStyle = valid ? 'green' : 'red'
                ctx.fillRect(hoverCell.x * CELL_SIZE, hoverCell.y * CELL_SIZE, info.size * CELL_SIZE, info.size * CELL_SIZE)
                ctx.globalAlpha = 1
            }
        }
    }, [buildings, selectedBuilding, hoverCell, activeBuilding, movingBuildingId])

    // --- Event Handlers ---

    function handleMouseMove(e: React.MouseEvent<HTMLCanvasElement>) {
        const rect = canvasRef.current!.getBoundingClientRect()
        const scaleX = (GRID_SIZE * CELL_SIZE) / rect.width
        const scaleY = (GRID_SIZE * CELL_SIZE) / rect.height
        const x = Math.floor(((e.clientX - rect.left) * scaleX) / CELL_SIZE)
        const y = Math.floor(((e.clientY - rect.top) * scaleY) / CELL_SIZE)
        setHoverCell({ x, y })
    }

    async function handleCanvasClick() {
        if (!hoverCell) return

        // Placing or moving
        if (selectedBuilding) {
            const info = BUILDING_INFO[selectedBuilding]
            if (!info) return

            if (!isValidPlacement(hoverCell.x, hoverCell.y, info.size, movingBuildingId || undefined)) return

            try {
                if (movingBuildingId) {
                    const res = await protectedFetch(`/api/buildings/${movingBuildingId}/move`, 'PUT', {
                        pos_x: hoverCell.x,
                        pos_y: hoverCell.y,
                    })
                    if (!res.ok) throw new Error((await res.json()).error)
                } else {
                    const res = await protectedFetch('/api/buildings', 'POST', {
                        building_type: info.type,
                        building_name: selectedBuilding,
                        pos_x: hoverCell.x,
                        pos_y: hoverCell.y,
                    })
                    if (!res.ok) throw new Error((await res.json()).error)
                }

                await loadBuildings()
                cancelAction()
            } catch (error: any) {
                alert(error.message || 'Action failed')
            }
            return
        }

        // Selecting an existing building
        const clicked = buildings.find((b) =>
            hoverCell.x >= b.pos_x && hoverCell.x < b.pos_x + b.size &&
            hoverCell.y >= b.pos_y && hoverCell.y < b.pos_y + b.size,
        )
        setActiveBuilding(clicked || null)
    }

    async function handleUpgrade() {
        if (!activeBuilding) return
        try {
            const res = await protectedFetch(`/api/buildings/${activeBuilding.id}/upgrade`, 'PUT')
            if (!res.ok) throw new Error((await res.json()).error)
            await loadBuildings()
            setActiveBuilding(null)
        } catch (error: any) {
            alert(error.message || 'Upgrade failed')
        }
    }

    function cancelAction() {
        setSelectedBuilding(null)
        setHoverCell(null)
        setMovingBuildingId(null)
        setActiveBuilding(null)
    }

    // --- Render ---

    return (
        <div style={{ display: 'flex' }}>
            {/* Sidebar */}
            <div style={{ width: '280px', padding: '16px', display: 'flex', flexDirection: 'column', gap: '8px' }}>

                {/* Info Panel */}
                {activeBuilding && !selectedBuilding && (
                    <div>
                        <h3>{activeBuilding.building_name}</h3>
                        <p>Level: {activeBuilding.level}</p>
                        <p>HP: {activeBuilding.hp}</p>

                        <div style={{ display: 'flex', gap: '8px', marginTop: '12px' }}>
                            <button style={{ flex: 1 }} onClick={() => {
                                setMovingBuildingId(activeBuilding.id)
                                setSelectedBuilding(activeBuilding.building_name)
                            }}>
                                Move
                            </button>
                            <button style={{ flex: 1 }} onClick={handleUpgrade}>
                                Upgrade
                            </button>
                        </div>
                    </div>
                )}

                {/* Shop */}
                <h2>Shop</h2>

                <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                    {Object.entries(BUILDING_INFO).map(([name, info]) => (
                        <button
                            key={name}
                            onClick={() => {
                                setSelectedBuilding(name)
                                setMovingBuildingId(null)
                                setActiveBuilding(null)
                            }}
                        >
                            {name} ({info.size}x{info.size})
                        </button>
                    ))}
                </div>

                {selectedBuilding && (
                    <button onClick={cancelAction} style={{ marginTop: '16px' }}>
                        {movingBuildingId ? 'Cancel Move' : 'Cancel Placement'}
                    </button>
                )}
            </div>

            {/* Canvas */}
            <div style={{ flex: 1, display: 'flex', justifyContent: 'center' }}>
                <canvas
                    ref={canvasRef}
                    width={GRID_SIZE * CELL_SIZE}
                    height={GRID_SIZE * CELL_SIZE}
                    onMouseMove={handleMouseMove}
                    onClick={handleCanvasClick}
                    onMouseLeave={() => setHoverCell(null)}
                    style={{ cursor: selectedBuilding ? 'crosshair' : 'default' }}
                />
            </div>
        </div>
    )
}
