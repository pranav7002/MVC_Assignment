'use client'

import VillageCanvas from './VillageCanvas'
import { useRequireAuth } from '../utils/authGuard'

export default function VillagePage() {
    useRequireAuth()
    return <VillageCanvas />
}
