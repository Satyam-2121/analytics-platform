# Analytics Platform

A scalable event-driven analytics backend built with Go, PostgreSQL, ClickHouse, Redis, and Docker.

## Overview

Analytics Platform is a backend system designed to collect, store, process, and query application events at scale. The platform provides a REST API that allows client applications to send events such as page views, user actions, and product interactions. These events are ingested through a Go API, validated using API-key authentication, stored in ClickHouse for high-performance analytics, and exposed through analytics endpoints.

The project demonstrates modern backend engineering concepts including multi-database architecture, repository pattern implementation, middleware-based authentication, event ingestion pipelines, analytics aggregation, and Redis-based caching.

The primary goal of this project is to simulate the architecture used by modern analytics products such as Google Analytics, Mixpanel, and PostHog while maintaining simplicity and developer-friendly deployment.

---

## Key Features

### Event Ingestion API

Applications can send analytics events through a REST endpoint.

Example:

```json
{
  "event": "page_view",
  "page": "/pricing",
  "user_id": "u123",
  "project_id": 1
}
```

Supported event categories include:

* Page Views
* User Actions
* Product Interactions
* Custom Business Events

---

### Project Management

Projects are managed in PostgreSQL.

Each project receives a unique UUID API key used to authenticate analytics requests.

Example:

```text
64c28ccf-36c5-4785-af7b-5d1bffb9af64
```

---

### Analytics Query Layer

The platform exposes aggregated analytics through REST endpoints.

Current metrics:

* Total Events
* Top Pages
* Page View Counts

Example query:

```http
GET /analytics/events?project_id=1
```

Response:

```json
{
  "total_events": 1200
}
```

---

### Redis Caching

Frequently requested analytics results are cached in Redis.

Benefits:

* Reduced database load
* Lower latency
* Faster dashboard responses
* Improved scalability

Cache Flow:

```text
Request
   ↓
Redis Lookup
   ↓
Cache Hit → Return Response
   ↓
Cache Miss
   ↓
ClickHouse Query
   ↓
Store Result in Redis
   ↓
Return Response
```

---

## System Architecture

```text
                    ┌─────────────────┐
                    │ Client / Browser│
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │     Go API      │
                    └────────┬────────┘
                             │
         ┌───────────────────┼───────────────────┐
         │                   │                   │
         ▼                   ▼                   ▼
 ┌──────────────┐   ┌──────────────┐   ┌──────────────┐
 │ PostgreSQL   │   │ ClickHouse   │   │ Redis Cache  │
 └──────────────┘   └──────────────┘   └──────────────┘
         │                   │                   │
         │                   │                   │
 Users & Projects      Analytics Events      Cached Results
```

---

## Technology Stack

### Backend

* Go 1.24+
* Standard Library HTTP Server
* Context-based Database Operations

### Databases

#### PostgreSQL

Used for transactional data:

* Users
* Projects
* API Keys

#### ClickHouse

Used for analytical workloads:

* Event Storage
* Aggregations
* Analytics Queries

Why ClickHouse?

* Columnar storage engine
* High compression ratio
* Extremely fast analytical queries
* Designed for large-scale event data

#### Redis

Used for:

* Analytics caching
* Fast lookups
* Reduced database pressure

---

## Design Patterns

### Repository Pattern

Database logic is separated from business logic.

Repositories:

```text
UserRepository
ProjectRepository
EventRepository
AnalyticsRepository
```

Benefits:

* Cleaner architecture
* Easier testing
* Database abstraction
* Improved maintainability

---

### Middleware Pattern

Authentication is implemented using HTTP middleware.

```text
Request
   ↓
API Key Middleware
   ↓
Route Handler
   ↓
Database
```

Responsibilities:

* Extract API Key
* Validate Project
* Reject Unauthorized Requests

---

## API Endpoints

### Create Event

```http
POST /events
```

Headers:

```http
X-API-Key: PROJECT_API_KEY
```

Request:

```json
{
  "event": "page_view",
  "page": "/pricing",
  "user_id": "u123",
  "project_id": 1
}
```

---

### Get Total Events

```http
GET /analytics/events?project_id=1
```

Response:

```json
{
  "total_events": 1200
}
```

---

### Get Top Pages

```http
GET /analytics/pages?project_id=1
```

Response:

```json
[
  {
    "page": "/pricing",
    "count": 420
  }
]
```

---

## Project Structure

```text
analytics-platform
├── cmd
│   └── api
├── internal
│   ├── database
│   ├── handlers
│   ├── middleware
│   ├── models
│   ├── repositories
│   ├── server
│   └── services
├── migrations
├── docker-compose.yml
├── go.mod
└── README.md
```
Updated project documentation.
Analytics API secured with project API keys.
