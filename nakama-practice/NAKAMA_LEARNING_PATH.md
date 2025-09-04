# ✅ Nakama Learning Path — Cricket Matchmaking System Edition

**Audience:** Developers new to Nakama implementing a cricket matchmaking system.  
**Time commitment:** 2–3 weeks (1–2 hours/day)  
**Prerequisites:** Basic Go, Docker, SQL fundamentals.  

---

## Phase 0 — Environment Setup (Day 0)
- [x] Install: Docker Desktop, Go, REST client, WebSocket client  
- [x] Clone starter repo  
- [x] Start Nakama + Postgres via docker-compose  
- [x] Verify health (`/healthz`, Console)  

**Milestone:** Nakama + Postgres run locally, logs visible, console accessible.  

---

## Phase 1 — Nakama Fundamentals & First Steps (Days 1-2)
- [ ] Understand Nakama as a backend (not a game engine)
- [ ] Learn core building blocks: RPCs, Storage, Authentication
- [ ] Explore basic Nakama concepts and architecture
- [ ] Build your first "Hello World" RPC

**Hands-on Lab A:**
- [ ] Create simple RPC (`echo` function)
- [ ] Test basic storage operations
- [ ] Explore Nakama console and logs
- [ ] Understand the runtime module concept

**Key Concepts:**
- [ ] Runtime modules and how they work
- [ ] RPC functions and HTTP endpoints
- [ ] Basic storage operations
- [ ] Nakama's event-driven architecture

**Milestone:** Comfortable with basic Nakama operations and concepts.  

---

## Phase 2 — Storage & Data Management (Days 2-3)
- [ ] Deep dive into Nakama's storage engine
- [ ] Learn collections, keys, versioning, permissions
- [ ] Design data schemas for cricket teams and players
- [ ] Practice with real storage operations

**Hands-on Lab B:**
- [ ] Create user profiles with `WriteStorageObjects`
- [ ] Read and update player data with `ReadStorageObjects`
- [ ] Design JSON schemas for flexible data
- [ ] Test storage permissions and access control

**Storage Concepts:**
- [ ] Collections and keys structure
- [ ] JSON object storage and validation
- [ ] Per-user vs. shared data strategies
- [ ] When to use Nakama vs. external databases

**Milestone:** Confident with Nakama storage and data modeling.  

---

## Phase 3 — Authentication & User Management (Days 3-4)
- [ ] Implement device login and user authentication
- [ ] Learn about user sessions and management
- [ ] Create user profiles and preferences
- [ ] Understand Nakama's user lifecycle

**Hands-on Lab C:**
- [ ] Implement device authentication flow
- [ ] Create user profile management RPCs
- [ ] Handle user preferences and settings
- [ ] Test user session management

**Authentication Features:**
- [ ] Device ID authentication
- [ ] User profile creation and updates
- [ ] Session management and validation
- [ ] User data privacy and permissions

**Milestone:** Users can authenticate and manage profiles.  

---

## Phase 4 — Real-time Communication (Days 4-5)
- [ ] Learn WebSocket connections and channels
- [ ] Implement real-time messaging between clients
- [ ] Create notification systems
- [ ] Handle presence and status updates

**Hands-on Lab D:**
- [ ] Set up WebSocket connections
- [ ] Create chat channels for teams
- [ ] Implement presence tracking
- [ ] Build real-time notification system

**Real-time Features:**
- [ ] WebSocket connection management
- [ ] Channel-based messaging
- [ ] Presence and status tracking
- [ ] Real-time notifications

**Milestone:** Real-time communication working between clients.  

---

## Phase 5 — Introduction to Matchmaking (Days 5-6)
- [ ] Learn Nakama's matchmaker concepts
- [ ] Understand how `AddMatchmaker` works
- [ ] Explore basic filtering and properties
- [ ] Build simple matchmaking flow

**Hands-on Lab E:**
- [ ] Create basic matchmaker with filters
- [ ] Test `AddMatchmaker` with simple properties
- [ ] Implement basic match proposal system
- [ ] Handle match acceptance/rejection

**Matchmaking Basics:**
- [ ] How Nakama's matchmaker works
- [ ] Filtering and property matching
- [ ] Match proposal and acceptance flow
- [ ] Basic match state management

**Milestone:** Basic matchmaking functionality working.  

---

## Phase 6 — Advanced Matchmaking & Pools (Days 6-8)
- [ ] Design pool management system
- [ ] Implement MMR-based team matching
- [ ] Create region-based filtering
- [ ] Build match acceptance workflow

**Hands-on Lab F:**
- [ ] Design pool data structures
- [ ] Implement MMR calculation algorithms
- [ ] Create region-based matching logic
- [ ] Build team acceptance flow

**Advanced Features:**
- [ ] Pool creation and management
- [ ] MMR algorithms and rating systems
- [ ] Region-based team matching
- [ ] Complex matchmaking logic

**Milestone:** Advanced matchmaking with pools and MMR working.  

---

## Phase 7 — External Integrations & Research (Days 8-9)
- [ ] Research GameLift integration possibilities
- [ ] Explore external PostgreSQL for complex queries
- [ ] Document integration requirements
- [ ] Plan architecture for external services

**Research Areas:**
- [ ] How Nakama runtime modules call external APIs
- [ ] External database integration strategies
- [ ] GameLift integration challenges
- [ ] Architecture for complex business logic

**Documentation:**
- [ ] Integration architecture diagrams
- [ ] API contract specifications
- [ ] External service requirements
- [ ] Performance and scaling considerations

**Milestone:** Clear understanding of external integration needs.  

---

## Phase 8 — Client Integration & Polish (Days 9-10)
- [ ] Design client-server communication protocols
- [ ] Implement reconnection handling
- [ ] Create user-friendly error handling
- [ ] Polish the overall user experience

**Client Features:**
- [ ] Message format standards
- [ ] Error handling and user feedback
- [ ] Reconnection and state recovery
- [ ] UI flow optimization

**User Experience:**
- [ ] Intuitive matchmaking flow
- [ ] Clear status updates and notifications
- [ ] Smooth error recovery
- [ ] Responsive real-time updates

**Milestone:** Polished client experience with robust error handling.  

---

## Phase 9 — Testing & Edge Cases (Days 10-11)
- [ ] Design comprehensive test suites
- [ ] Handle edge cases and error scenarios
- [ ] Test reconnection and failure scenarios
- [ ] Plan for production deployment

**Testing Strategy:**
- [ ] Unit tests for core logic
- [ ] Integration tests with Nakama
- [ ] Edge case testing (disconnections, failures)
- [ ] Performance testing considerations

**Edge Cases:**
- [ ] Network disconnections
- [ ] Partial team scenarios
- [ ] System failures and recovery
- [ ] Performance under load

**Milestone:** Robust system with comprehensive testing.  

---

## Phase 10 — Observability & Operations (Days 11-12)
- [ ] Explore Nakama's built-in monitoring
- [ ] Plan custom metrics and logging
- [ ] Design operational dashboards
- [ ] Document operational procedures

**Monitoring & Logging:**
- [ ] Nakama's Prometheus endpoints
- [ ] Custom metrics for cricket events
- [ ] Logging strategies and aggregation
- [ ] Operational dashboards

**Operations:**
- [ ] Health monitoring and alerts
- [ ] Performance tracking
- [ ] Error rate monitoring
- [ ] Capacity planning

**Milestone:** Clear observability strategy and operational plan.  

---

## Capstone Project — *Cricket Matchmaking System v1* (2–3 days)
- [ ] End-to-end matchmaking flow
- [ ] Team registration and pool management
- [ ] MMR-based matching with region filtering
- [ ] Real-time updates and notifications
- [ ] Error handling and edge cases

**Complete Flow:**
- [ ] User authentication and profile setup
- [ ] Team creation and pool assignment
- [ ] Matchmaking with MMR and region filtering
- [ ] Match proposal and acceptance
- [ ] Real-time status updates throughout

**Deliverables:**  
- [ ] Complete matchmaking service
- [ ] API documentation and client examples
- [ ] Test suite and deployment guides
- [ ] Architecture documentation

**Milestone:** Production-ready cricket matchmaking system.  

---

## Implementation Notes & Learning Balance

### **Nakama Fundamentals (Phases 1-4):**
- [ ] Core concepts and architecture
- [ ] Storage and data management
- [ ] Authentication and user management
- [ ] Real-time communication

### **Matchmaking Specific (Phases 5-6):**
- [ ] Nakama matchmaker concepts
- [ ] Pool and MMR systems
- [ ] Team-based matching logic

### **Advanced Topics (Phases 7-10):**
- [ ] External integrations
- [ ] Testing and edge cases
- [ ] Observability and operations

### **What You'll Learn:**
- [ ] **Nakama Basics:** RPCs, Storage, WebSockets, Authentication
- [ ] **Matchmaking:** How to use Nakama's built-in features
- [ ] **Custom Logic:** Building business-specific functionality
- [ ] **Integration:** Connecting with external services
- [ ] **Production:** Testing, monitoring, and deployment

---

## Daily 30-Minute Drills (Optional)
- [ ] Read 1 Nakama documentation page
- [ ] Test 1 new Nakama API or feature
- [ ] Write 1 small test or example
- [ ] Review 1 concept or architecture decision

---

## Common Pitfalls & Tips
- [ ] Start with simple examples before complex features
- [ ] Test Nakama features in isolation first
- [ ] Use Nakama's built-in features when possible
- [ ] Plan for external services early
- [ ] Keep match state simple and minimal
- [ ] Always validate client inputs
- [ ] Test reconnection scenarios thoroughly

---

## Reference Map
- [ ] [Official Docs](https://heroiclabs.com/docs)  
- [ ] [Getting Started](https://heroiclabs.com/docs/nakama/getting-started/)  
- [ ] [Runtime Code Basics](https://heroiclabs.com/docs/nakama/concepts/runtime-code-basics/)  
- [ ] [Storage API](https://heroiclabs.com/docs/nakama/concepts/storage/)  
- [ ] [Matches & Matchmaker](https://heroiclabs.com/docs/nakama/concepts/matches/)  
- [ ] [Go Runtime Examples](https://github.com/heroiclabs/nakama-common)  

---

## What to Learn Next (Future Versions)
- [ ] Advanced MMR algorithms and ELO systems
- [ ] Anti-cheat and security features
- [ ] Multi-region deployment strategies
- [ ] Advanced analytics and machine learning
- [ ] GameLift integration for server allocation
- [ ] Advanced pool management algorithms
