const sprites: Record<string, HTMLImageElement> = {}

if (typeof window !== 'undefined') {
    sprites['Town Hall'] = new Image()
    sprites['Town Hall'].src = '/sprites/red.png'

    sprites['Cannon'] = new Image()
    sprites['Cannon'].src = '/sprites/green.png'

    sprites['Gold Mine'] = new Image()
    sprites['Gold Mine'].src = '/sprites/blue.png'
}

export default sprites