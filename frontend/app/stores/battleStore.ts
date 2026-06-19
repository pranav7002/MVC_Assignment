import { create } from 'zustand'

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

interface BattleStore {
    defenderID: string | null
    buildings: Building[]
    setBattle: (defenderID: string, buildings: Building[]) => void
    clear: () => void
}

export const useBattleStore = create<BattleStore>((set) => ({
    defenderID: null,
    buildings: [],
    setBattle: (defenderID: string, buildings: Building[]) =>
        set({ defenderID, buildings }),
    clear: () => set({ defenderID: null, buildings: [] }),
}))
