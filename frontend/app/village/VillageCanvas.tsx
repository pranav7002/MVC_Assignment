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
}

const GRID_SIZE = 20
const CELL_SIZE = 40

export default function VillageCanvas() {
    const token = useAuthStore((state) => state.token)

    const [buildings, setBuildings] = useState<Building[]>([])

    const canvasRef = useRef<HTMLCanvasElement>(null)
    const grassTileRef = useRef<HTMLImageElement | null>(null)

    useEffect(() => {
        const img = new Image()
        img.src = '/sprites/grass.png'

        img.onload = () => {
            grassTileRef.current = img
            setBuildings((prev) => [...prev])
        }
    }, [])

    useEffect(() => {
        async function loadBuildings() {
            try {
                const res = await protectedFetch('/api/buildings', 'GET')

                const data = await res.json()

                setBuildings(data.data)
            } catch (err) {
                console.error(err)
            }
        }

        if (token) {
            loadBuildings()
        }
    }, [token])

    useEffect(() => {
        const canvas = canvasRef.current
        if (!canvas) return

        const ctx = canvas.getContext('2d')
        if (!ctx) return

        ctx.clearRect(0, 0, canvas.width, canvas.height)

        const grassTile = grassTileRef.current
        if (grassTile) {
            drawMap(ctx, grassTile, CELL_SIZE)
        }

        // Draw buildings
        buildings.forEach((building) => {
            const sprite = sprites[building.building_name]

            if (sprite && sprite.complete) {
                ctx.drawImage(
                    sprite,
                    building.pos_x * CELL_SIZE,
                    building.pos_y * CELL_SIZE,
                    building.size * CELL_SIZE,
                    building.size * CELL_SIZE,
                )
            }
        })
    }, [buildings])

    return (
        <div className="flex justify-center items-center min-h-screen">
            <canvas
                ref={canvasRef}
                // kept at 800 for some weird resolution reason
                width={GRID_SIZE * 40}
                height={GRID_SIZE * 40}

                style={{
                    width: `${GRID_SIZE * 44}px`,
                    height: `${GRID_SIZE * 44}px`,
                }}
            />
        </div>
    )
}

function drawMap(
    ctx: CanvasRenderingContext2D,
    tile: HTMLImageElement,
    cellSize: number,
) {
    for (let x = 0; x < GRID_SIZE; x++) {
        for (let y = 0; y < GRID_SIZE; y++) {
            ctx.drawImage(tile, x * cellSize, y * cellSize, cellSize, cellSize)
        }
    }
}
