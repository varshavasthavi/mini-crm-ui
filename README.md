# Mini CRM Automation Project

## Project Overview
This project is a simulation of an event-driven CRM system. It demonstrates a complete data flow where player behavior (deposits) is ingested via an API, evaluated against business rules, and persisted across two different database technologies for analytical and operational purposes.

## Project Structure
```text
mini-crm/
├── scripts/
│   └── schema.sql         # ClickHouse Schema
├── backend/
│   ├── main.go            # API & Rule Logic
│   ├── go.mod             # Dependencies
│   └── go.sum             # Checksums
└── frontend/
    ├── src/
    │   └── App.vue        # Vue.js UI
    └── package.json       # JS Dependencies

### Core Business Logic
* **Ingestion**: Backend receives a deposit event containing `player_id` and `amount`.
* **Rule Evaluation**: The system evaluates a hardcoded rule: `deposit_amount >= 1000`.
* **Action**: If the rule matches, a `BONUS_MESSAGE` action is triggered (mocked).
* **Persistence**:
    * **ClickHouse**: Analytical storage of every incoming event.
    * **MongoDB**: Operational log of triggered automation actions (`campaign_logs`).

---

## Setup & Execution Instructions

### 1. Prerequisites
* **Go**: Version 1.21+
* **Node.js & NPM**: For the Vue.js frontend.
* **MongoDB**: Running on `localhost:27017`.
* **ClickHouse**: Running on `localhost:9000`.

### 2. Database Initialization (ClickHouse)
Ensure the correct authentication settings (no password for the `default` user) and per the requirement to keep database lifecycle management simple, please execute the following SQL to create the `events` table before starting the backend:

CREATE TABLE IF NOT EXISTS events (
    event_id UUID,
    player_id String,
    event_type String,
    amount Float64,
    timestamp DateTime64(3)
) ENGINE = MergeTree()
ORDER BY timestamp;

### 3. Backend Execution (Go)
1. Navigate to the `/backend` directory.
2. Install dependencies: `go mod tidy`.
3. Run the application: `go run main.go`
* The server will be active at `http://localhost:8080`.

### 4. Frontend Execution (Vue.js)
1. Navigate to the `/frontend` directory.
2. Install dependencies: `npm install`.
3. Start the development server: `npm run serve`
* The UI will be accessible at `http://localhost:8081`.

---

## Assumptions
* **Database Connectivity**: Standard local instances of MongoDB and ClickHouse are used for development/testing without authentication.
* **Mock Action**: The `BONUS_MESSAGE` is represented by a specific entry in the `campaign_logs` collection to demonstrate the rule was triggered.
* **Simple Fetch**: To keep logic straightforward, the UI fetches logs from the backend immediately after a successful event submission rather than using WebSockets.
* **Schema Management**: The backend assumes the ClickHouse table already exists; schema is provided in the repository for manual initialization.

---

