# ChatApp — Real-time Messaging

A full-stack real-time chat application built as a monorepo.

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | React 19 + Vite 8 |
| Backend | Go (Gorilla WebSocket + Mux) |
| Database | PostgreSQL 16 |
| Auth | JWT + bcrypt |
| Real-time | WebSocket |

## Project Structure

```
chat-app/
├── frontend/     # React + Vite app
├── server/       # Go backend
├── docker-compose.yml
└── .env
```

## Getting Started

### Prerequisites
- Node.js 18+
- Go 1.22+
- PostgreSQL 16 (or Docker)

### 1. Start Database

```bash
# Using Docker
docker-compose up -d

# Or create the database manually
createdb chatapp
```

### 2. Start Backend

```bash
cd server
go run ./cmd/server
```

The server starts on `http://localhost:8080`.

### 3. Start Frontend

```bash
cd frontend
npm install
npm run dev
```

The app opens at `http://localhost:5173`.

## Features

- ✅ JWT Authentication (Register/Login)
- ✅ Real-time messaging via WebSocket
- ✅ User search from database (no contact creation needed)
- ✅ Online/Offline status indicators
- ✅ Typing indicators
- ✅ Read receipts
- ✅ Message history with pagination
- ✅ Premium dark theme UI
- ✅ Auto-reconnect WebSocket
