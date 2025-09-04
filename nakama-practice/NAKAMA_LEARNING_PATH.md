# ✅ Nakama Learning Path — Mini Battle Arena Edition

**Audience:** Developers new to Nakama who want a practical, end-to-end path from zero to production-ready features.  
**Time commitment:** 2–4 weeks (1–2 hours/day)  
**Prerequisites:** Basic Go or TypeScript/JavaScript, Docker, SQL fundamentals.  

---

## Phase 0 — Environment Setup (Day 0)
- [ ] Install: Docker Desktop, Go (or Node), REST client, WebSocket client  
- [ ] Clone starter repo  
- [ ] Start Nakama + Postgres via docker-compose  
- [ ] Verify health (`/healthz`, Console)  

**Milestone:** Nakama + Postgres run locally, logs visible, console accessible.  

---

## Phase 1 — Core Mental Model (Day 1)
- [ ] Understand Nakama as a backend (not a game engine)  
- [ ] Learn building blocks: RPCs, Matches, Matchmaker, Storage, Leaderboards, Social, Realtime Socket  

**Milestone:** Map Mini Battle Arena features → Nakama primitives.  

---

## Phase 2 — First RPC (Day 1)
- [ ] Implement trivial RPC (`echo`)  
- [ ] Compile runtime module (Go or JS)  
- [ ] Register and call RPC via HTTP  
- [ ] Add validation + errors  

**Hands-on Lab A:**  
- [ ] `submit_profile` RPC → stores player profile  
- [ ] `get_profile` RPC → reads player profile  

**Milestone:** First custom RPC working.  

---

## Phase 3 — Storage Engine Deep Dive (Days 2–3)
- [ ] Learn collections, keys, versioning, permissions  
- [ ] Explore upserts, conditional writes, batch ops  
- [ ] Design storage for profiles, settings, battle history  

**Hands-on Lab B:**  
- [ ] `PUT /settings` RPC → save preferences  
- [ ] `GET /settings` RPC → fetch preferences  
- [ ] Add schema validation  

**Milestone:** Player data safely stored/retrieved.  

---

## Phase 4 — Leaderboards (Days 3–4)
- [ ] Learn leaderboard concepts: sort orders, reset schedules  
- [ ] Practice with `LeaderboardCreate`, `LeaderboardRecordWrite`, `LeaderboardRecordsList`  

**Hands-on Lab C:**  
- [ ] Create `daily_wins` leaderboard (daily reset)  
- [ ] `submit_score` RPC → validates & writes record  
- [ ] `top10_scores` RPC → fetches leaderboard  

**Milestone:** Battle winners appear on leaderboard.  

---

## Phase 5 — Matchmaker & Matches (Days 4–6)
- [ ] Learn matchmaker basics (tickets, properties)  
- [ ] Explore authoritative match loop (init, join, loop, leave, terminate)  
- [ ] Persistence: save results, not full state  

**Hands-on Lab D:**  
- [ ] Matchmaker RPC wrapper (add/cancel queue)  
- [ ] Authoritative match:  
  - [ ] Track joins/leaves  
  - [ ] Broadcast tick every 5s  
  - [ ] Handle “attack/defend” actions  
  - [ ] End after 30s → decide winner → save result  

**Milestone:** Players complete Mini Battle Arena matches.  

---

## Phase 6 — Authentication & Social (Days 6–7)
- [ ] Implement Device login  
- [ ] Understand friend graph: add, accept, block  
- [ ] Explore chat & groups  

**Hands-on Lab E:**  
- [ ] Device login flow  
- [ ] Friend request → chat channel  

**Milestone:** Players authenticate, add friends, and chat.  

---

## Phase 7 — Events, Hooks, and Observability (Days 7–9)
- [ ] Learn before/after hooks  
- [ ] Setup Prometheus metrics + Grafana dashboards  
- [ ] Enable tracing with OpenTelemetry  

**Hands-on Lab F:**  
- [ ] `AfterAuthenticate` hook → log streak + send notification  
- [ ] Export Prometheus metrics  
- [ ] Create Grafana dashboard (players online, matches played)  

**Milestone:** Observe and debug Mini Battle Arena.  

---

## Phase 8 — Deployment & Scale (Days 9–12)
- [ ] Run Nakama stateless behind LB  
- [ ] Setup sticky sessions / session handoff  
- [ ] Use Postgres replica, Redis for caching  
- [ ] Practice migrations + rollout strategies  

**Hands-on Lab G:**  
- [ ] Containerize server module  
- [ ] Deploy to kind/minikube  
- [ ] Setup autoscaling + liveness/readiness probes  

**Milestone:** Backend scales safely in cluster.  

---

## Phase 9 — Security & Hardening (Days 12–14)
- [ ] Manage secrets + rotate keys  
- [ ] Validate inputs server-side  
- [ ] Add rate-limiting + abuse detection  
- [ ] Consider GDPR/PII + data retention  

**Hands-on Lab H:**  
- [ ] Add rate-limit to `submit_score` RPC  
- [ ] Rotate admin/console credentials  

**Milestone:** Mini Battle Arena backend hardened.  

---

## Capstone Project — *Mini Battle Arena v1* (2–3 days)
- [ ] Device login → profile setup  
- [ ] Queue in matchmaking → play battle  
- [ ] Match ends → store result  
- [ ] Update leaderboard  
- [ ] Invite friend → chat  
- [ ] Notify rank-up  

**Deliverables:**  
- [ ] Server module code + Docker manifest  
- [ ] Postman/Insomnia collection  
- [ ] Grafana dashboard JSON  
- [ ] Runbook for troubleshooting  

**Milestone:** End-to-end Mini Battle Arena backend complete.  

---

## Daily 30-Minute Drills (Optional)
- [ ] Read 1 Nakama doc page  
- [ ] Write 1 tiny test RPC or hook  
- [ ] Review 1 log/metric anomaly  

---

## Common Pitfalls & Tips
- [ ] Match Nakama + Go build versions  
- [ ] Always validate client payloads  
- [ ] Use Storage versioning  
- [ ] Keep match state minimal  
- [ ] Prefer server-authoritative logic  

---

## Reference Map
- [ ] [Official Docs](https://heroiclabs.com/docs)  
- [ ] [Go Runtime Examples](https://github.com/heroiclabs/nakama-common)  
- [ ] [Client SDKs](https://heroiclabs.com/docs/nakama/concepts/client-libraries/)  
- [ ] [Console & Admin](https://heroiclabs.com/docs/nakama/console/)  
- [ ] [Monitoring](https://heroiclabs.com/docs/nakama/operations/monitoring/)  

---

## What to Learn Next
- [ ] Advanced match loops: lockstep, rollback, ECS  
- [ ] Anti-cheat with server replay verification  
- [ ] Multi-region matchmaking + latency strategies  
- [ ] LiveOps: events, A/B tests, feature flags  
