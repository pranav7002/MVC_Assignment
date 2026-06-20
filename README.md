# Vanguard — Clash of Clans Style Game

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)

Everything else runs inside containers

## Setup

```bash
docker compose up --build
```

This command:
1. Starts a **PostgreSQL** database
2. Runs **schema migrations**
3. Seeds **16 test accounts** (15 bots + 1 rich test user)
4. Starts the **Go backend** on `http://localhost:8080`
5. Starts the **Next.js frontend** on `http://localhost:3000`

Open **http://localhost:3000** in your browser.

## Test Accounts

All accounts use password: `password`

| Username | Town Hall | Gold | Elixir | Trophies |
|---|---|---|---|---|
| `kai` | 1 | 1,000,000 | 1,000,000 | 0 |
| `player_1` – `player_5` | 1 | 500 | 500 | 30–150 |
| `player_6` – `player_10` | 2 | 1,200 | 1,200 | 180–300 |
| `player_11` – `player_15` | 3 | 3,000 | 3,000 | 330–450 |

Use `kai` for unrestricted building/upgrading/training. Use the `player_*` accounts to test matchmaking from different trophy/TH ranges. Try to test everything before battling using the account with unlimited money as it might bring the resources down to their actual cap according to the storage buildings. 

## Testing the Battle Flow

1. **Login** as `kai` (or register a new account)
2. **Village** — buy buildings from the shop, upgrade them, move them around
3. **Troops** — navigate to the troops page and train an army
4. **Matchmaking** — click the ⚔ icon to find an opponent, press **Skip** for a different one
5. **Battle** — click **Attack** to start a WebSocket battle, drop troops on the canvas

## Battle Features

The battle simulation runs server side. Most of these mechanics are not visible in the UI due to the absence of animations, but they are a part of the engine:

1. **Troop-specific targeting** — Goblins prioritize resource buildings, Giants prioritize defenses. Both fall back to any building if preferred targets are destroyed.
2. **AOE damage** — Mortars deal splash damage to all troops within their AOE range.
3. **Mortar dead zone** — Mortars have a minimum range, so troops inside the dead zone are safe from mortar fire.
4. **BFS pathfinding** — 8-directional BFS with diagonal corner-cutting prevention (skips diagonal if either adjacent cardinal cell is blocked).
5. **Defense cooldown** — Defenses have a 10-tick cooldown between attacks, preventing instant-kill spam.
6. **Variable troop speed** — Each troop type moves at a different fractional speed per tick (Giants are slow, Goblins are fast).
7. **Dynamic grid** — Destroyed buildings free their grid cells, allowing troops to pathfind through the rubble.

## Useful Commands

```bash
# Start everything
docker compose up --build

# Start fresh (wipe database)
docker compose down -v && docker compose up --build

# Stop everything
docker compose down
