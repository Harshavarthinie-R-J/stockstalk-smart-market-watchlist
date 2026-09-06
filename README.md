# StockStalk Backend

Backend service for **StockStalk**, a market watchlist that identifies meaningful changes in the stocks a user follows and helps explain what deserves attention.

## Tech Stack

* Go
* GraphQL
* PostgreSQL
* Redis
* JWT
* Docker

## Current Backend Responsibilities

* User registration and authentication
* JWT-based authorization
* Stock/instrument search and lookup
* Watchlist creation and management
* Market-data abstraction
* Market quotes and snapshots
* User checkpoints
* Meaningful-change detection
* "Since Last Visit" comparison
* GraphQL queries and mutations

## Core Flow

```text
User
→ Watchlist
→ Market Data
→ Snapshot
→ Checkpoint
→ Change Detection
→ GraphQL
→ Frontend
```

## Backend Structure

```text
cmd/server
    Application entry point and dependency wiring

internal/auth
    Registration, login, JWT and authentication middleware

internal/user
    User model, service and repository

internal/instrument
    Stock identity, metadata and search

internal/watchlist
    Watchlist and stock membership management

internal/marketdata
    Market models, provider abstraction and snapshots

internal/checkpoint
    User watchlist checkpoints

internal/change
    Meaningful market-change detection

internal/graphql
    GraphQL schema, resolvers and server

internal/db
    PostgreSQL connection, migrations and seed logic

pkg/logger
    Shared logging

pkg/response
    Shared response structures

mock
    Development market-data source
```

## Product Logic

StockStalk does not simply display current prices.

It compares the current market state with the user's previous checkpoint and identifies changes that deserve attention.

The initial change signals include:

* Significant price movement
* 52-week high/low
* Notable market changes

The advanced intelligence layer will later add:

* Event Fusion
* Personal Baselines
* Salience Scoring
* Context and Attribution
* Event Memory

## API

GraphQL is the primary client-facing API.

### Queries

```text
me
market
quote
searchInstruments
watchlists
watchlist
changesSinceLastVisit
```

### Mutations

```text
register
login
createWatchlist
renameWatchlist
deleteWatchlist
addStock
removeStock
createCheckpoint
```

### Future Subscriptions

```text
marketUpdate
meaningfulChange
watchlistEvent
```

## Engineering Principles

* Keep GraphQL resolvers thin.
* Keep business logic inside services.
* Keep persistence behind repository interfaces.
* Keep market-data providers replaceable.
* Keep processing deterministic and testable.
* Enforce user ownership in the backend.
* Handle stale and unavailable data explicitly.
* Prefer modular architecture over premature microservices.
* Design for horizontal scalability without unnecessary infrastructure.

## Development

Run tests:

```bash
go test ./...
```

Run the backend:

```bash
go run ./cmd/server
```

Health endpoint:

```text
GET /health
```

GraphQL endpoint:

```text
/graphql
```

## Development Data

The current development environment uses:

```text
mock/market_data.json
```

The data source is isolated behind the market-data layer so it can later be replaced by a live provider without changing the rest of the application.

## Architecture Direction

The backend is designed to evolve from:

```text
Go
→ Services
→ In-memory repositories
→ Mock market data
```

to:

```text
Go
→ GraphQL
→ Services
→ PostgreSQL + Redis
→ Live market data
→ Event processing
```

The long-term goal is to provide a scalable backend that converts raw market activity into meaningful, explainable watchlist events.
