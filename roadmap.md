# Roadmap — Go API Manager (GAM)

> **Read this first.** This document is our contract for the next 3 weeks. It defines *what* we're building, *why* we're building it this way, and *in what order*. The [README.md](README.md) explains the architecture in depth (with diagrams). Every phase below ends with one or more Markdown docs added to `docs/`, containing full working code, line-by-line explanations, and the theory behind it. We do not move to the next phase until the current one runs and is understood — not just copy-pasted.

## Working agreement (please read)

- **You know zero Go.** Every phase that introduces new Go syntax/concepts calls it out explicitly in a "New Go Concepts" box before we write code. We reuse your Java/OOP background as the anchor for analogies (e.g. "a Go `struct` + methods is like a Java class without inheritance").
- **Part-time pace.** Planned around ~2–4 focused hours/day. Each day below is a *target*, not a deadline — if a day spills into two, we absorb it into that week's buffer day rather than panicking.
- **Phase, not day, is the real unit of progress.** We only write the phase's doc(s) once the code for that phase actually runs locally. Docs describe what we *actually built*, not what we planned to build.
- **Docs contract.** Each phase produces `docs/phase-NN-<name>.md` containing: (1) the concept/theory, (2) full annotated code, (3) an explanation section in plain English, (4) a "why this way and not another way" section, (5) how to run/test it.
- **Realism check-in.** Every Sunday/end-of-week (Day 7, 14, 21) is a **buffer + review day** — catch up on slippage, consolidate docs, and re-baseline the following week if needed. Scope bends before the deadline does.
- **Definition of Done for the whole project:** all 5 services run together via one `docker compose up`, a request can flow end-to-end (publish an API → get a key → call it through the gateway → see it rate-limited → see it in analytics → see it in Grafana/Jaeger), and you can explain every box in the architecture diagram unscripted in an interview.

---

## Scope decisions (locked in before Phase 1)

These were deliberate trade-offs to make a 3-week, part-time, zero-Go-experience timeline achievable **without** making the project look like a toy. Each one is also a talking point for interviews ("I chose X over Y because...").

| Decision | Choice | Why |
|---|---|---|
| Feature scope | **Core gateway + essentials** (routing/proxying, auth, rate limiting, analytics, minimal dev portal) — not a full WSO2/Apigee-style enterprise suite (no monetization, multi-tenant orgs, GraphQL, marketplace) | Enterprise API managers took large teams years. Nailing the core control-plane/data-plane split cleanly is worth more in an interview than a shallow pass over 20 features. Enterprise features are listed as **stretch goals** at the bottom. |
| Deployment target | **Docker Compose only** (Kubernetes explicitly out of scope, noted as a stretch goal) | Kubernetes adds a second, unrelated learning curve (manifests, Helm, kubelet internals) on top of learning Go from zero. Compose still proves you understand multi-service orchestration, networking, and health checks — the actual microservices story doesn't need k8s to be told well. |
| Web framework | `net/http` (standard library) + [`chi`](https://github.com/go-chi/router) router | `chi` is a thin layer over `net/http` (no magic, no reflection-heavy binding like some frameworks). This teaches you *real* Go HTTP handling instead of a framework's abstractions — a stronger signal in interviews than "I used Gin for everything." |
| Database access | `sqlx` + raw SQL (Postgres), no ORM | Understanding your queries and schema cold is worth more than GORM magic, and this is a small enough project that an ORM buys us nothing. |
| Cache / rate-limit store | **Redis** | Industry-standard for counters/token buckets and cache invalidation; tiny operational footprint. |
| Messaging / async backbone | **NATS JetStream**, not Kafka | Kafka is heavier to operate (Zookeeper/KRaft, partitions, consumer groups) than what a 3-week solo project needs. NATS gives us the same "event-driven, decoupled services" story with a fraction of the ops overhead — and *why we didn't pick Kafka* is itself a good interview answer. |
| Internal service-to-service calls | **gRPC** (Gateway → Auth, Gateway → Rate Limiter) | External traffic (client → gateway) is REST/JSON because that's what API consumers expect. Internal traffic (service → service) is gRPC — this "north-south REST, east-west gRPC" split is a real pattern used at scale and is a strong signal of architectural maturity. |
| Auth model | JWT (signed, stateless) **and** API keys (stored hashed in Postgres), OAuth2 client-credentials flow implemented manually | Building the token issuance/validation ourselves (rather than delegating to a library like `oauth2`) is where the actual learning is. |
| Observability | Structured logs (`log/slog`, stdlib — Go 1.21+), Prometheus metrics, Grafana dashboards, OpenTelemetry + Jaeger tracing | This is the entire pitch for the **SRE** side of your CV. Non-negotiable, even though it's compressed into Week 3. |
| Repo layout | **Monorepo**, multiple Go modules tied together with a `go.work` workspace file | Lets each service be an independently versioned/built Go module (realistic microservices practice) while still letting you `go run` everything from one repo without publishing packages. `go.work` is a lesser-known but very useful Go 1.18+ feature — good to know. |
| Frontend / Dev Portal | Minimal static HTML/JS page (not a full React app) that lists published APIs and lets a developer self-register for an API key | Keeps the project "full-stack" (you'll have a browser-facing screen to demo) without spending scarce time on frontend tooling instead of Go/SRE depth. |

**Explicitly out of scope for the 3-week core project (stretch goals, see bottom of this doc):** Kubernetes/Helm, API monetization/billing, multi-tenant orgs, GraphQL/gRPC-Web for external clients, an API marketplace, full OAuth2 Authorization Code flow with a real login UI, service mesh (Istio/Linkerd).

---

## System at a glance

Five services. See [README.md](README.md) for full diagrams.

1. **Gateway** (data plane) — reverse proxy, dynamic routing, calls Auth + Rate Limiter, publishes analytics events.
2. **Auth Service** — issues/validates API keys & JWTs.
3. **Management Service** (control plane) — CRUD for registering upstream APIs/routes; also serves the public "catalog" the dev portal reads from.
4. **Rate Limiter Service** — token-bucket quota enforcement backed by Redis.
5. **Analytics Service** — consumes request events off NATS, aggregates, exposes a query API.

Infra: PostgreSQL, Redis, NATS JetStream, Prometheus, Grafana, Jaeger.

---

## Week 1 — Go Foundations + Core Domain Services

**Goal by end of week:** You can read/write basic-to-intermediate Go, and two of the five services (Auth, Management) run standalone, tested, backed by real Postgres.

### Day 1 (today) — Phase 1: Architecture & Roadmap
- Deliverables: `roadmap.md`, `README.md` (this session).
- No code yet — this is the design phase every real project starts with.

### Day 2 — Phase 2: Go Crash Course, Part 1
**New Go concepts:** toolchain (`go mod`, `go run`, `go build`), packages, variables & types, `struct`s vs Java classes, methods with receivers, pointers (Go's version of references), error handling (`error` as a value, not exceptions), zero values, slices & maps.
- Output: `docs/phase-02-go-fundamentals-part1.md` — theory + tiny runnable snippets (not part of the real services yet).

### Day 3 — Phase 3: Go Crash Course, Part 2
**New Go concepts:** interfaces (structural typing — very different from Java's explicit `implements`), goroutines & channels (concurrency primitives), `context.Context` (cancellation/timeouts threaded through calls — you'll use this *everywhere*), table-driven tests with `testing` package, Go modules vs `go.work` workspaces.
- Output: `docs/phase-03-go-fundamentals-part2.md`.

### Day 4–5 — Phase 4: Auth Service
- Build: API key generation + hashing (bcrypt), JWT issuance/validation (`golang-jwt`), Postgres schema for clients/keys (`sqlx` + migrations via `golang-migrate`), REST endpoints (`chi`) for issuing/validating credentials, unit tests.
- **New Go concepts:** structuring an HTTP service (handlers/middleware), dependency injection without a framework (plain constructor functions), working with a SQL driver, environment-based config.
- Output: `docs/phase-04-auth-service.md` (full code + explanation), service runs and is curl-able standalone.

### Day 6 — Phase 5: Management Service (Control Plane)
- Build: Postgres schema for `apis`/`routes`, admin REST API (create/update/delete a published API + its upstream + its path pattern), public "catalog" read endpoint for the future dev portal.
- **New Go concepts:** request validation, pagination, structuring a slightly larger service (routes/handlers/repository layers).
- Output: `docs/phase-05-management-service.md`.

### Day 7 — Buffer & Week 1 Review
- Catch up on anything that slipped. Re-read Phase 2–5 docs cold and confirm you can explain them without looking. Adjust Week 2 dates if needed.

---

## Week 2 — Gateway, Rate Limiting, gRPC, and Wiring It Together

**Goal by end of week:** All 5 services run together via `docker compose up`, and a request can flow end-to-end: client → Gateway → (gRPC) Auth check → (gRPC) rate-limit check → proxied to a fake upstream → response returned → event fired to NATS.

### Day 8–9 — Phase 6: Gateway Service
- Build: `httputil.ReverseProxy`-based proxy, in-memory route table loaded from Management Service (polling first, event-driven refresh comes in Phase 8), calls to Auth Service over plain REST for now (gRPC comes next phase so the concepts don't collide).
- **New Go concepts:** `net/http/httputil`, building middleware chains, graceful shutdown basics (`context` + `signal.NotifyContext`).
- Output: `docs/phase-06-gateway-service.md`.

### Day 10 — Phase 7: Rate Limiter Service
- Build: token-bucket algorithm against Redis (`go-redis`), gRPC service definition (`.proto` file), generated Go server code.
- **New Go concepts:** protobuf + gRPC code generation (`protoc`, `buf`), writing a gRPC server, Redis client basics (atomic INCR/EXPIRE, Lua scripting for atomicity).
- Output: `docs/phase-07-rate-limiter-grpc.md`.

### Day 11 — Phase 8: Wire Gateway → Auth & Rate Limiter over gRPC
- Refactor Gateway's Auth calls from REST to gRPC, add the Rate Limiter gRPC client, add route-table refresh via NATS (Management Service publishes `api.route.updated`, Gateway subscribes and invalidates its cache).
- **New Go concepts:** gRPC clients, connection pooling/keep-alives, pub/sub with NATS.
- Output: `docs/phase-08-service-to-service-grpc-and-events.md`.

### Day 12–13 — Phase 9: Docker Compose Integration
- Build: multi-stage Dockerfiles per service (small final images), `docker-compose.yml` wiring all 5 services + Postgres + Redis + NATS, health checks, `.env`-based config, a `Makefile` for common commands.
- **New Go concepts:** none new in Go itself — this is the "make it look and run like a real deployed system" phase.
- Output: `docs/phase-09-docker-compose.md`.

### Day 14 — Phase 10: Analytics Service + Buffer
- Build: NATS consumer that ingests `request.completed` events from the Gateway, aggregates into Postgres (per-API request counts, latencies, error rates), a small query API to read aggregates back.
- Output: `docs/phase-10-analytics-service.md`. Remaining time = Week 2 buffer/catch-up.

---

## Week 3 — Observability, Resilience, and Interview-Readiness

**Goal by end of week:** The system is observable (metrics/traces/logs), resilient to failure, documented, tested under load, and you have a rehearsed way to explain it in an interview.

### Day 15 — Phase 11: Structured Logging & Metrics
- Build: `log/slog` structured logging across all services, Prometheus `/metrics` endpoints, Grafana dashboards (requests/sec, error rate, latency p95, rate-limit rejections).
- **New Go concepts:** `log/slog` (stdlib), `prometheus/client_golang` instrumentation patterns (counters/histograms/middleware).
- Output: `docs/phase-11-metrics-and-logging.md`.

### Day 16 — Phase 12: Distributed Tracing
- Build: OpenTelemetry SDK wired into all 5 services, trace context propagated over both HTTP and gRPC, Jaeger UI to visualize a single request hopping across services.
- **New Go concepts:** OTel SDK/instrumentation, context propagation across process boundaries.
- Output: `docs/phase-12-distributed-tracing.md`.

### Day 17 — Phase 13: Resilience Patterns
- Build: timeouts on every outbound call, retries with exponential backoff (gRPC → Rate Limiter), a circuit breaker (`sony/gobreaker` or hand-rolled) around the Auth call, graceful shutdown everywhere, sane defaults when Redis/NATS are briefly unavailable.
- **New Go concepts:** `context.WithTimeout`, backoff/jitter patterns, the circuit breaker state machine.
- Output: `docs/phase-13-resilience.md`.

### Day 18 — Phase 14: API Docs + Minimal Dev Portal
- Build: OpenAPI/Swagger spec for the Gateway-exposed and Management admin APIs, a static HTML/JS page (Dev Portal) that lists published APIs from the catalog endpoint and lets a developer self-register for an API key.
- Output: `docs/phase-14-api-docs-and-dev-portal.md`.

### Day 19 — Phase 15: Testing & Load Testing
- Build: integration tests that spin up the Compose stack, a `k6` or `vegeta` load test script, capture results (throughput, latency under load, what breaks first).
- Output: `docs/phase-15-testing-and-load-testing.md`.

### Day 20 — Phase 16: CI/CD
- Build: GitHub Actions workflow — lint (`golangci-lint`), test, build & tag Docker images per service.
- Output: `docs/phase-16-ci-cd.md`.

### Day 21 — Phase 17: Final Review & Interview Pitch
- Final pass on README (screenshots/GIFs of Grafana + dev portal), a short demo script, and a written "how to talk about this project" doc: the problem it solves, the architecture trade-offs made (table above), what you'd do differently at 10x scale, and 3–5 STAR-format stories pulled from real snags you hit while building it.
- Output: `docs/phase-17-interview-pitch.md`.

---

## Stretch goals (only if time remains — do not let these threaten Week 3)

- Port the Compose deployment to Kubernetes (kind/minikube) with Helm charts and an HPA on the Gateway.
- Full OAuth2 Authorization Code flow with a real login page.
- API monetization (usage-based billing tiers).
- Multi-tenant organizations (teams owning sets of APIs).
- GraphQL or gRPC-Web gateway for external clients.
- Canary/blue-green routing at the Gateway.

---

## Next step

Say "let's start Day 2 / Phase 2" whenever you're ready, and we'll begin the Go crash course before writing a single line of the real services.
