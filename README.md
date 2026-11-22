# NetWatch

> A local-first home network monitor with an offline AI assistant.

NetWatch is a desktop app that keeps an eye on the devices on your home network (routers, Raspberry Pis, NAS boxes, smart plugs, self-hosted services) and tells you when something goes down. Everything runs on your own machine. The Electron app launches a containerized Go backend, and an optional AI assistant powered by a local small language model can answer questions and take actions in plain English through an MCP server.

> **Status:** 🚧 Early development. The architecture is defined; the features below are the target for v1.

---

## Features

- **Device monitoring:** ICMP ping, TCP port, HTTP(S) and DNS checks on configurable intervals
- **Real-time status:** live updates over GraphQL subscriptions; no refreshing
- **Smart alerting:** failure/recovery thresholds to avoid flapping, plus incident tracking
- **History & uptime:** latency charts, uptime percentages and hourly rollups with automatic retention
- **Desktop native:** tray icon with overall status, native notifications and LAN device discovery
- **Offline AI assistant (optional):** ask *"How often did the Pi drop off this week?"* or *"Pause checks on the NAS for 30 minutes"*, answered by a local model with no cloud calls
- **MCP server (optional):** expose NetWatch as tools for any MCP-compatible client
- **One-click backend:** the desktop app starts, health-checks and stops its own containers

---

## Architecture

```
┌──────────────────────── Electron ────────────────────────┐
│  Renderer: dashboard · devices · incidents · chat panel  │
│      │ GraphQL (queries, mutations, subscriptions)       │
│      │                         ▲ IPC                     │
│  Main process: container manager · tray · agent loop     │
└──────┼─────────────────────────┼───────────┼─────────────┘
       │                         │ MCP       │ LLM (OpenAI-compatible)
       ▼                         ▼           ▼
┌─────────────── Docker Compose (localhost only) ──────────────┐
│                                                              │
│   api ◄──GraphQL── mcp            ollama ◄── ollama-pull     │
│    │                [mcp]          [ai]        [ai]          │
│    ▼                                                         │
│   Postgres ◄────── poller          migrate (one-shot)        │
│        LISTEN/NOTIFY                                         │
└──────────────────────────────────────────────────────────────┘
```

The **GraphQL API is the single backend contract.** The Electron UI and the MCP server are both clients of it, so business logic lives in one place and every change, whether made by you or by the assistant, streams to the UI in real time.

### Services

| Service       | Role                                                                                         |
|---------------|----------------------------------------------------------------------------------------------|
| `db`          | Postgres (optionally TimescaleDB): devices, checks, results, status, incidents, rollups     |
| `migrate`     | One-shot setup: schema migrations, DB roles, retention policies, idempotent seed data       |
| `api`         | Go GraphQL API (gqlgen): queries, mutations, live subscriptions via Postgres `LISTEN/NOTIFY` |
| `poller`      | Go check engine: scheduler, worker pool, checkers, up/down state machine, rollups           |
| `mcp`         | *(profile `mcp`)* MCP server that maps curated tools onto GraphQL operations                 |
| `ollama`      | *(profile `ai`)* local model runtime with an OpenAI-compatible endpoint                     |
| `ollama-pull` | *(profile `ai`)* one-shot step that makes sure the configured model is downloaded           |

`migrate`, `api` and `poller` are subcommands of a single Go binary built into one image.

---

## Tech Stack

- **Backend:** Go · gqlgen · sqlc · goose · pro-bing
- **Database:** PostgreSQL (TimescaleDB optional)
- **Desktop:** Electron · TypeScript · GraphQL client with `graphql-ws`
- **AI:** Model Context Protocol · Ollama · a small tool-calling model (configurable)
- **Infrastructure:** Docker Compose with profiles, healthchecks and one-shot init containers

---

## Getting Started

### Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (or Docker Engine + Compose v2)
- Node.js LTS and Go (for development)
- *(Optional)* an NVIDIA GPU with the NVIDIA Container Toolkit for faster local inference. On macOS, install [Ollama](https://ollama.com) on the host for GPU acceleration; NetWatch uses it automatically if it's running.

### Run the backend only

```bash
cp .env.example .env
docker compose up -d                          # core: db, migrate, api, poller
docker compose --profile mcp up -d            # + MCP server
docker compose --profile mcp --profile ai up -d   # + local AI assistant
```

### Run the desktop app (development)

```bash
cd desktop
npm install
npm run dev
```

The desktop app detects running containers and attaches to them, or starts the stack itself if nothing is running.

---

## Configuration

| Variable      | Description                                   | Default    |
|---------------|-----------------------------------------------|------------|
| `DB_USER`     | Postgres application user                     | `netwatch` |
| `DB_PASSWORD` | Generated on first launch by the desktop app  | —          |
| `DB_NAME`     | Database name                                 | `netwatch` |
| `AI_MODEL`    | Ollama model tag used by the assistant        | *(TBD)*    |
| `APP_VERSION` | Backend image tag; set by the desktop app     | app version |

All published ports bind to `127.0.0.1`. Nothing is exposed to your LAN.

---

## Using the MCP Server with Other Clients

With the `mcp` profile enabled, NetWatch can be added to any MCP client over streamable HTTP:

```
http://127.0.0.1:<mcp-port>/mcp
```

The desktop app's settings page shows the current URL and can copy a ready-made client config.

---

## Project Structure

```
.
├── backend/            # Go: api, poller, migrate subcommands
│   ├── cmd/
│   ├── graph/          # GraphQL schema + resolvers
│   ├── internal/       # service layer, checkers, scheduler
│   ├── db/             # migrations + sqlc queries
│   └── Dockerfile
├── mcp/                # MCP server config + curated GraphQL operations
├── desktop/            # Electron app (main, preload, renderer)
├── docker-compose.yml
└── .env.example
```

---

## Roadmap

- [ ] Schema, migrations and seed data
- [ ] Poller with ICMP + TCP checks
- [ ] GraphQL queries and mutations
- [ ] Status state machine + live subscriptions
- [ ] Electron dashboard, tray icon and notifications
- [ ] Container lifecycle management from the desktop app
- [ ] HTTP/TLS and DNS checks, rollups and retention
- [ ] MCP server
- [ ] Local AI assistant with confirm-before-write actions
- [ ] Host-side LAN discovery
- [ ] Packaged installers (macOS / Windows / Linux)

---

## License

TBD
