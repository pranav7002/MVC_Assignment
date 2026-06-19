const sprites: Record<string, HTMLImageElement> = {}

if (typeof window !== 'undefined') {
    const buildingSprites: Record<string, number> = {
        'Town Hall': 4,
        'Cannon': 4,
        'Archer Tower': 4,
        'Gold Mine': 3,
        'Elixir Collector': 3,
        'Gold Storage': 3,
        'Elixir Storage': 3,
        'Mortar': 2,
        'Training Grounds': 2,
    }

    const toFileName = (name: string) =>
        name.toLowerCase().replace(/ /g, '_')

    for (const [name, maxLevel] of Object.entries(buildingSprites)) {
        for (let level = 1; level <= maxLevel; level++) {
            const key = `${name}_${level}`
            sprites[key] = new Image()
            sprites[key].src = `/sprites/buildings/${toFileName(name)}_${level}.png`
        }
    }

    const troopNames = ['Barbarian', 'Archer', 'Goblin', 'Giant', 'Wizard']
    for (const name of troopNames) {
        sprites[name] = new Image()
        sprites[name].src = `/sprites/troops/${name.toLowerCase()}.png`
    }
}

export default sprites