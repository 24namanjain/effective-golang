## **1. Golang Fundamentals & Best Practices**

### **Core Concepts**

- **Project Structure**
    - Keep a clear and modular package structure.
    - Avoid circular dependencies.
    - Group by functionality, not just type (e.g., `auth`, `game`, `leaderboard`).
- **Modules & Dependencies**
    - Use `go mod` for dependency management.
    - Avoid unnecessary external libraries; prefer stdlib where possible.
- **Data Types**
    - Understand `struct`, `map`, `slice`, `array`, `interface`.
    - Use `struct` for strongly typed data; keep JSON tags consistent.
- **Error Handling**
    - Always check for `err` explicitly.
    - Wrap errors with context using `fmt.Errorf("...: %w", err)`.
- **Concurrency & Goroutines**
    - Use `context.Context` for cancelation and deadlines.
    - Avoid goroutine leaks — always ensure they exit on cancellation.
    - Protect shared resources with `sync.Mutex` or channels.
    - Use worker pools for high concurrency workloads.
- **Channels**
    - Buffered channels for rate-limiting and queueing.
    - Select statements for handling multiple channel inputs.
- **Testing**
    - Write table-driven tests.
    - Use `testing` package and benchmarks with `go test -bench`.

---

### **Naming Conventions**

- **Packages** → lowercase, no underscores (`leaderboard`, `authservice`).
- **Structs/Interfaces** → PascalCase (`GameManager`, `UserStore`).
- **Variables/Functions** → camelCase (`maxScore`, `initCache`).
- **Constants** → ALL_CAPS with underscores (`MAX_RETRY_COUNT`).
- **Receivers** → short and meaningful (`func (gm *GameManager)`).

---

## **2. Nakama Backend Concepts**

### **Core Components**

- **RPC (Remote Procedure Call)**
    - Custom server functions callable from the client or other services.
    - Use for authenticated actions (e.g., `JoinGame`, `SubmitScore`).
    - Register with `nk.RegisterRpc("id", handlerFunc)`.
- **Events**
    - Use Nakama’s event system for reacting to player activity.
    - Examples:
        - `AfterAuthenticate` → track login streaks.
        - `AfterMatchEnd` → award rewards.
- **Leaderboards**
    - Create with `nk.LeaderboardCreate` — define sort order & reset schedule.
    - Submit scores via `nk.LeaderboardRecordWrite`.
    - Query with `nk.LeaderboardRecordsList`.
- **Matches & Matchmaker**
    - Use authoritative matches for server-controlled logic.
    - Store match state in memory or Redis for persistence.
- **Storage Engine**
    - Key-value store per collection per user.
    - Use JSON objects for flexible schema.
- **Realtime Events**
    - Handle chat, match data streams, and presence events.

---

## **3. Database & Caching**

### **PostgreSQL**

- **Core Concepts**
    - ACID compliance, relational schema design.
    - Use indexes to speed up queries; avoid over-indexing.
    - Joins vs subqueries — prefer joins when possible.
    - Transactions for atomic multi-step changes.
- **Best Practices**
    - Normalize for data integrity; denormalize for performance if needed.
    - Always parameterize queries to avoid SQL injection.
    - Use `EXPLAIN ANALYZE` for query optimization.

---

### **Redis**

- **Core Concepts**
    - In-memory data store for caching and fast lookups.
    - Data types: String, Hash, List, Set, Sorted Set.
    - TTL (Time-To-Live) for ephemeral data.
- **Use Cases in Backend**
    - Session caching.
    - Leaderboard storage (Sorted Sets).
    - Rate limiting (using counters with TTL).
    - Pub/Sub for event broadcasting.
- **Best Practices**
    - Avoid storing large payloads in Redis.
    - Use hash slots correctly in Redis Cluster.
    - Always handle connection pool limits.

1. Git work flow with tag on each commit , branch pattern /dev/mukul , release candidates.