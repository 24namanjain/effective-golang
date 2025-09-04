## Nakama Learning Path

Audience: Developers new to Nakama who want a practical, end-to-end path from zero to production-ready features.
Time commitment: 2–4 weeks (1–2 hours/day)
Prerequisites: Basic Go or TypeScript/JavaScript, Docker, SQL fundamentals.

---

### Phase 0 — Environment Setup (Day 0)
- Install: Docker Desktop, Go (or Node), a REST client (curl/HTTPie/Insomnia), and WebSocket client.
- Clone a starter repo
- Start Nakama + Postgres via docker-compose.
- Verify health:
  - HTTP: `curl http://localhost:7350/healthz`
  - Console: http://localhost:7351 (default admin/password — change for real use)

Milestone: You can run Nakama locally, view logs, and access the console.

---

### Phase 1 — Core Mental Model (Day 1)
- What Nakama is: a realtime game backend platform. Not a game engine.
- Key building blocks:
  - RPCs: Server functions callable via HTTP/gRPC/WebSocket.
  - Matches: Realtime sessions run by server code (authoritative or relayed).
  - Matchmaker: Finds players based on properties.
  - Storage Engine: Per-user/collection key-value JSON storage.
  - Leaderboards: Sorted/aggregated scoreboards with reset/sort windows.
  - Friends/Groups/Chat: Social graph & messaging primitives.
  - Realtime Socket: Presence, chat streams, match data, notifications.

Milestone: You can describe when to use RPC vs Match vs Storage vs Leaderboards.

---

### Phase 2 — First RPC (Day 1)
- Implement a trivial RPC (e.g., `echo` or `simple_test`).
- Compile a Go runtime module (or write JS/TS server code if using JS runtime).
- Register RPC and call it via HTTP:
  - HTTP: `POST /v2/rpc/<rpc_id>?http_key=defaultkey` with body payload.
- Add input validation and structured errors.

Milestone: You deploy and invoke a custom RPC and see logs in Nakama.

Hands-on Lab A:
- Create `submit_profile` RPC that writes a JSON profile to storage (`collection: profile`).
- Add `get_profile` RPC to read it back.

---

### Phase 3 — Storage Engine Deep Dive (Days 2–3)
- Entities: `collection`, `key`, `user_id`, `version`, `permission_read`, `permission_write`.
- Access patterns: upsert, conditional write (version check), batch ops.
- Indexing strategy: denormalize when necessary; use Postgres for analytics via export.

Hands-on Lab B:
- Build a settings system:
  - `PUT /settings` RPC → upsert JSON settings.
  - `GET /settings` RPC → read.
  - Add a server-side schema validator (simple struct checks).

Milestone: You can persisently store and retrieve player data safely.

---

### Phase 4 — Leaderboards (Days 3–4)
- Concepts: sort order (best/highest/lowest), operators (set/increment), reset schedule.
- APIs: `LeaderboardCreate`, `LeaderboardRecordWrite`, `LeaderboardRecordsList`.
- Scoping leaderboards (daily/weekly/seasonal), anti-cheat basics on server-side.

Hands-on Lab C:
- Create `daily_score` leaderboard with daily reset.
- RPC: `submit_score` → validates score server-side, writes record.
- RPC: `top10_scores` → fetches and returns top 10 with owner’s rank.

Milestone: You publish a leaderboard and integrate client reads/writes.

---

### Phase 5 — Matchmaker & Matches (Days 4–6)
- Matchmaker:
  - Tickets with properties: region, skill, mode, latency, device.
  - Filters and numeric ranges; party matchmaker basics.
- Matches:
  - Relayed vs Authoritative (server-driven state).
  - Match loop: `MatchInit`, `MatchJoin`, `MatchLoop`, `MatchLeave`, `MatchTerminate`.
  - Persistence strategy: ephemeral in-memory + write outcomes to storage/SQL.

Hands-on Lab D:
- Implement a matchmaker add/cancel RPC wrapper (if not using client matchmaker directly).
- Create an authoritative match that:
  - Tracks player joined/left.
  - Broadcasts a tick message every N seconds.
  - Ends after a fixed timer and writes a result record.

Milestone: Players can join a match, exchange messages, and see match end.

---

### Phase 6 — Authentication & Social (Days 6–7)
- Auth flows:
  - Device/Custom/Facebook/Apple/Steam/Epic/etc.
  - Link/unlink multiple identities per account.
- Friends/Groups:
  - Friend graph, blocks.
  - Group creation, membership, edges.
- Notifications & In-app mail.

Hands-on Lab E:
- Implement device auth in a minimal client.
- Add friend request flow and chat channel between friends.

Milestone: You authenticate users and support basic social/communication.

---

### Phase 7 — Events, Hooks, and Observability (Days 7–9)
- Before/After hooks: customize behavior around core operations.
- Runtime Logger; error classification.
- Metrics: Prometheus endpoint, dashboards (Grafana).
- Tracing (OpenTelemetry) and correlation IDs for RPCs/matches.

Hands-on Lab F:
- Add `AfterAuthenticate` hook to write a login streak entry and send a notification.
- Export Prometheus metrics and add a basic Grafana dashboard.

Milestone: You can trace requests across RPCs and diagnose issues quickly.

---

### Phase 8 — Deployment & Scale (Days 9–12)
- Stateless Nakama nodes behind a load balancer.
- Sticky sessions for sockets, or session handoff.
- Externalizing state: Postgres primary/replica, Redis for caching/match state.
- Migrations & schema change process.
- Blue/green or canary deploys.

Hands-on Lab G:
- Containerize your server module; deploy to a managed Kubernetes (kind/minikube OK for practice).
- Configure horizontal pod autoscaling and liveness/readiness probes.

Milestone: Your backend runs in a cluster with safe rollout and recovery.

---

### Phase 9 — Security & Hardening (Days 12–14)
- Secrets management (keys, OAuth secrets).
- Input validation everywhere; server-authoritative enforcement.
- Rate-limiting, abuse detection.
- GDPR/PII considerations, data retention.

Hands-on Lab H:
- Add rate-limiting middleware on your RPC gateway (or within RPCs).
- Rotate console/admin credentials and runtime keys.

Milestone: You’ve closed obvious security gaps and added rate-limiters.

---

### Capstone Project — "Match & Score" (2–3 days)
- Device auth → profile setup → queue in matchmaker → get paired → authoritative match → post result → leaderboard rank → friend invite & chat → notification on rank up.
- Deliverables:
  - Server module code and Docker manifest.
  - Postman/Insomnia collection for core flows.
  - Grafana dashboard JSON for metrics.
  - Runbook with troubleshooting steps.

---

### Daily 30-Minute Drills (Optional)
- Read 1 doc page (RPC, Matches, Storage, Leaderboards, etc.).
- Write 1 tiny test RPC or hook.
- Review 1 set of metrics/logs, explain anomalies.

---

### Common Pitfalls & Tips
- Go plugin compatibility: Build with the same Go version and matching deps as the Nakama image.
- Prefer server-authoritative logic for competitive modes.
- Always validate payloads server-side; never trust clients.
- Use Storage versioning for safe concurrent writes.
- Keep match state small; persist only what’s needed.

---

### Reference Map
- Official Docs: https://heroiclabs.com/docs
- Go Runtime Examples: https://github.com/heroiclabs/nakama-common
- Client SDKs: https://heroiclabs.com/docs/nakama/concepts/client-libraries/
- Console & Admin: https://heroiclabs.com/docs/nakama/console/
- Metrics/Monitoring: https://heroiclabs.com/docs/nakama/operations/monitoring/

---

### What to Learn Next
- Advanced authoritative match patterns (lockstep, rollback, ECS-driven loops).
- Anti-cheat strategies with server replay verification.
- Multi-region latency strategies and party matchmaking.
- LiveOps tooling (events, A/B tests, feature flags).
