# StockStalk

### Smart Market Watchlist

> **Your stocks. What changed. What matters.**

StockStalk is a smart market watchlist built to help users understand which changes in their tracked stocks actually deserve attention.

Instead of only showing price, percentage change, and volume, StockStalk compares the current market state with previous observations and historical behaviour to identify meaningful and unusual market events.

---

## Problem

Traditional stock watchlists show market data, but users still have to manually figure out:

- What changed since they last checked?
- Was the movement unusual?
- Is the trading volume abnormal?
- Which stock deserves attention first?
- What evidence supports the event?

StockStalk addresses this information overload by turning market observations into prioritized and explainable events.

---

## Key Features

### Smart Watchlists

- Create and manage watchlists
- Add, remove, and reorder stocks
- Search instruments
- View current market information

### Since You Last Checked

Compares the current market state with the user's previous checkpoint and highlights meaningful changes.

### Personal Baseline

Compares a stock's current behaviour with its own historical behaviour instead of using the same threshold for every stock.

### Event Fusion

Combines related signals such as:

- Price movement
- Volume anomaly
- Price anomaly
- 52-week conditions
- Relative movement

into a single meaningful event.

### Attention Score

Ranks events using measurable factors such as:

- Movement magnitude
- Historical unusualness
- Volume anomaly
- Relative movement
- Event importance
- Recency

### Context & Confidence

Provides supporting context while clearly indicating attribution confidence:

- HIGH
- MEDIUM
- LOW
- NO CLEAR CONTEXT

StockStalk does not automatically claim that a particular event caused a price movement without sufficient evidence.

### Event Memory

Important events are stored so users can review what happened previously.

### Historical Charts

View persisted price history for:

- 1 Day
- 1 Week
- 1 Month

### Stocks to Consider

Highlights stocks outside the current watchlist that show unusual observed activity.

This is for market awareness and discovery, **not a buy/sell recommendation**.

### Data Reliability

Market observations can be marked as:

- LIVE
- DELAYED
- STALE
- UNAVAILABLE
- CONFLICTING

---

## How It Works

Market Data
     ↓
Market Snapshot
     ↓
Historical Baseline
     ↓
Change Detection
     ↓
Signal Detection
     ↓
Event Fusion
     ↓
Context & Confidence
     ↓
Attention Score
     ↓
Event Memory
     ↓
Since You Last Checked


Technology Stack
Layer	Technology
Frontend	React + TypeScript
Build Tool	Vite
API Client	Apollo Client
API	GraphQL
Backend	Go
Database	PostgreSQL
Authentication	JWT
Password Security	bcrypt
Charts	Recharts
Local Database	Docker
Deployment	Render
Version Control	Git + GitHub


Architecture
                    User
                      |
                      ↓
             React + TypeScript
                      |
                   GraphQL
                      |
                      ↓
                 Go Backend
                      |
          +-----------+-----------+
          |                       |
          ↓                       ↓
   Event Intelligence       Watchlist / Auth
          |                       |
          +-----------+-----------+
                      |
                      ↓
                 PostgreSQL

                 
Backend Structure
backend/
├── cmd/
│   └── server/
├── internal/
│   ├── auth/
│   ├── change/
│   ├── checkpoint/
│   ├── config/
│   ├── db/
│   ├── graphql/
│   ├── instrument/
│   ├── marketdata/
│   ├── metrics/
│   ├── processor/
│   ├── user/
│   └── watchlist/
├── migrations/
│   ├── 001_init.sql
│   ├── 002_event_intelligence.sql
│   ├── 003_event_context.sql
│   └── 004_snapshot_history.sql
└── mock/
    └── market_data.json
    
Frontend Structure

frontend/
├── public/
├── src/
│   ├── components/
│   ├── pages/
│   └── lib/
├── package.json
└── vite.config.ts
Local Setup
Prerequisites

Install:

Node.js LTS
Go
Docker Desktop
Git
1. Clone the Repository
git clone https://github.com/Harshavarthinie-R-J/stockstalk-smart-market-watchlist.git
cd stockstalk-smart-market-watchlist
2. Start PostgreSQL
cd backend
docker compose up -d
3. Start Backend

In PowerShell:

$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/stockstalk?sslmode=disable"
$env:JWT_SECRET="change-me"

go run ./cmd/server

Backend:

http://localhost:8080

Health check:

http://localhost:8080/health

GraphQL:

http://localhost:8080/graphql

4. Start Frontend

Open another terminal:

cd frontend
npm install
npm run dev

Frontend:

http://localhost:5173

Environment Variables
Backend
APP_PORT=8080
APP_ENV=development
MARKET_DATA_FILE=mock/market_data.json
DATABASE_URL=postgres://postgres:postgres@localhost:5432/stockstalk?sslmode=disable
JWT_SECRET=change-me
Frontend

For production:

VITE_GRAPHQL_URL=https://stockstalk-smart-market-watchlist-1.onrender.com/graphql
Production Deployment

StockStalk is deployed using Render.

Frontend

StockStalk Frontend

Backend

StockStalk Backend

GraphQL API

GraphQL Endpoint

Deployment Architecture
Browser
   ↓
Render Frontend
   ↓
HTTPS / GraphQL
   ↓
Render Go Backend
   ↓
Render PostgreSQL
Database Migrations

The project uses PostgreSQL migrations:

001_init.sql
002_event_intelligence.sql
003_event_context.sql
004_snapshot_history.sql

The snapshot-history migration adds indexes for efficient historical data retrieval and snapshot deduplication.

Authentication

StockStalk uses:

JWT for authentication
bcrypt for secure password hashing

Authenticated GraphQL requests use:

Authorization: Bearer <JWT>
Design Principle

StockStalk is designed around one simple idea:

Don't watch every movement. Know what matters.

The system focuses on:

Market awareness
Event prioritization
Explainability
Historical comparison
Data reliability

rather than predicting the future price of a stock.

What StockStalk Does NOT Do

StockStalk intentionally does not provide:

Buy/sell recommendations
Automated trading
Portfolio optimization
Stock-price prediction
Brokerage execution

It is an intelligent market monitoring and awareness tool.

Future Enhancements
Live market-data providers
Real-time WebSocket events
News and company-event integration
Advanced volatility and sector baselines
Redis-based caching
Background event-processing workers
AI-assisted event summaries
Mobile application
