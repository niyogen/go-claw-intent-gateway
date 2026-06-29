# OpenClaw — Intern Task Assignment

> **Note:** The OpenClaw project is a fully backend Go service. Neither intern will build a traditional UI, but the tasks below align each intern to what most closely matches their skills and interests.

---

## Intern Profiles

| | Intern A | Intern B |
|---|---|---|
| **Background** | CS Major + Marketing Degree | Mechanical Engineering Major |
| **Strength** | Interfaces, user flows, developer experience, documentation, configuration schemas | Structured systems thinking, deterministic pipelines, precision engineering, rules and tolerances |
| **Assigned Domain** | Consumer-facing surfaces & developer experience | Core pipeline, security mechanics & data systems |

---

## Intern A — CS + Marketing
> *"Make it feel good to use from the outside."*

Intern A's work focuses on the API surfaces, developer-facing configuration, and the integration experience — effectively, how the outside world interacts with OpenClaw. This maps to the marketing mindset of understanding what the user (developer/admin) needs.

---

### Task 4 — Admin & App Registration API
**Assignee:** Intern A
**Status:** 🔲 Not Started

**Description:**
Build `com.niyogen.mapping/internal/adapter/ingress/admin/`. Implement the REST endpoints for:
- Registering a new application (`POST /v1/apps`)
- Managing Trigger Keys (`POST /v1/apps/{id}/trigger-keys`)
- Kill Switch (`DELETE /v1/apps/{id}/enable`)

Design clean, versioned JSON request/response schemas and write OpenAPI/Swagger comments for every endpoint.

**Why Intern A?** Designing API contracts is a UX problem. The marketing background means thinking about the developer as the customer.

---

### Task 6 — Multi-Channel Ingress Adapters (Cap 1)
**Assignee:** Intern A
**Status:** 🔲 Not Started

**Description:**
Build `com.niyogen.mapping/internal/adapter/ingress/` handlers for:
- WhatsApp (`/whatsapp/`)
- Slack (`/slack/`)
- Generic REST API (`/api/`)

Parse each channel's raw payload into the `domain.UnifiedMessage` struct and handle channel-specific auth header formats.

**Why Intern A?** Each channel integration is like a product integration story — understanding the contract each platform publishes and how to translate it.

---

### Task 8 — Schema Editing & Confirmation Flows (Cap 2)
**Assignee:** Intern A
**Status:** 🔲 Not Started

**Description:**
Build `com.niyogen.mapping/internal/usecase/schema/`. Implement the natural language schema update intent flow:
- Parse user's intent to modify a field/schema
- Produce a human-readable confirmation message back to the user
- Apply changes only after confirmed

**Why Intern A?** This is about conversation design and user experience — the marketing degree adds real value in crafting clear, trustworthy confirmation prompts.

---

### Task 10 — Review-Driven Experiments (Cap 5)
**Assignee:** Intern A
**Status:** 🔲 Not Started

**Description:**
Build `com.niyogen.mapping/internal/usecase/experiment/`. Implement the RAG review ingestion pipeline and interface with the `port.FeatureAnalyzer` LLM interface to produce `ExperimentProposal` structs. Format and return experiment proposals (Top Feature, Confidence, Supporting Quotes, Hypothesis).

**Why Intern A?** This is the most "marketing analytics" flavored task in the entire RFC — deriving product insights from user reviews.

---

### README & Documentation
**Assignee:** Intern A
**Status:** 🔲 Not Started

**Description:**
Expand the root `README.md` with:
- Project overview and architecture diagram (ASCII or Mermaid)
- Quickstart guide (how to register an app, send a message)
- Channel integration setup instructions
- API endpoint reference

**Why Intern A?** Developer documentation is marketing for engineers.

---

### Test Responsibility (Intern A)
Write tests inside `com.niyogen.test/` for their own tasks:
- `test/unit/adapter/` — Unit tests for each channel ingress adapter
- `test/e2e/` — End-to-end flows (e.g., WhatsApp message → intent dispatch → webhook response)
- `test/contract/` — API contract tests for the Admin endpoints

---

## Intern B — Mechanical Engineering
> *"Build the machine that never fails."*

Intern B's work focuses on the precision pipeline, security mechanics, and data systems — the parts of OpenClaw that must be deterministic, correct, and provably safe. Mechanical engineers excel at systems that operate within strict tolerances, which maps perfectly to security and data pipeline work.

---

### Task 2 — Security Foundations (Layers 1-2 & 4-7)
**Assignee:** Intern B
**Status:** 🚧 In Progress

**Description:**
Build all packages inside `com.niyogen.mapping/internal/security/`:

| Package | What to Build | RFC Layer |
|---|---|---|
| `security/hmac/` | HMAC-SHA256 signing + verification for inbound & outbound messages | Layer 1 |
| `security/replay/` | Nonce store with TTL — reject duplicate messages within 5-minute window | Layer 1 |
| `security/ratelimit/` | Per-user and per-tenant sliding window limits + Event Circuit Breaker | Layer 4 |
| `security/permission/` | Role × Trust Tier × Intent Scope enforcer matrix | Layer 4 |
| `security/ssrf/` | Outbound URL validation — block private IP ranges, loopback, metadata endpoints | Layer 5 |
| `security/pii/` | PII scrubber applied before any data crosses the LLM boundary | Layer 5 |
| `security/firewall/` | Semantic Firewall: block prompt injection patterns, enforce result size limits | Layer 6 |
| `security/ast/` | SQL AST parser for Dynamic Intents — validate queries before execution | Layer 5 |

**Why Intern B?** Security engineering is precision engineering. HMAC math, sliding windows, and AST parsing require the same tolerance-and-constraint thinking as mechanical systems design. There is no margin for error.

---

### Task 3 — In-Memory Mock Repositories (Phase 0a)
**Assignee:** Intern B
**Status:** 🔲 Not Started

**Description:**
Build `com.niyogen.mapping/internal/adapter/repository/memory/`. Implement all port interfaces with thread-safe in-memory maps:
- `AppRepository` (store/retrieve registered apps)
- `AuditRepository` (append-only log)
- `ReplayStore` (nonce TTL cache)

**Why Intern B?** In-memory data structures with concurrent safety constraints is a precision engineering problem — like designing a component to exact tolerances.

---

### Task 5 — Orchestrator Pipeline (Layers 3-6)
**Assignee:** Intern B
**Status:** 🔲 Not Started

**Description:**
Build `com.niyogen.mapping/internal/usecase/orchestrator/`. Implement the sequential execution pipeline:

```
Receive UnifiedMessage
  → Replay Guard check
  → Rate Limit check
  → Permission check (Role × Tier × Scope)
  → Intent Lookup from Registry
  → Param Extraction + Validation
  → PII Scrub
  → SSRF-validated Outbound Webhook call (HMAC-signed)
  → Audit Log
  → Return response to channel
```

**Why Intern B?** This is the core machine — each step is a stage in a manufacturing pipeline, and the output of each feeds the next with strict type contracts. Perfect for systems thinking.

---

### Task 7 — Event Notifications Rules Engine (Cap 3)
**Assignee:** Intern B
**Status:** 🔲 Not Started

**Description:**
Build `com.niyogen.mapping/internal/usecase/notification/`. Implement:
- YAML-based rule configuration loader (`config/notification_rules.yaml`)
- Rules evaluator: trigger conditions, payload mapping, rate-limiting circuit breaker
- Inbound event webhook receiver (`adapter/ingress/events/`)

**Why Intern B?** A rules engine with circuit breakers is a control system — directly analogous to feedback loops in mechanical/electrical engineering.

---

### Task 9 — Dynamic Intent Generation (Cap 4)
**Assignee:** Intern B
**Status:** 🔲 Not Started

**Description:**
Build `com.niyogen.mapping/internal/usecase/dynamic/`. Implement:
- YAML intent file builder from user natural language description
- SQL AST parser integration (from `security/ast`) for validating generated queries
- Promotion flow: draft → review → promote to registry

**Why Intern B?** Generating structured YAML from rules and validating SQL with an AST parser is deterministic, structured engineering — exactly the kind of precision work Intern B will excel at.

---

### Test Responsibility (Intern B)
Write tests inside `com.niyogen.test/` for their own tasks:
- `test/unit/security/` — Extensive adversarial unit tests for every security package:
  - SSRF bypass attempts (IPv6, DNS rebinding, encoded URLs)
  - Replay attack simulations
  - HMAC forgery attempts
  - AST injection test cases
- `test/integration/` — PostgreSQL/Redis integration tests when repositories are promoted from mock to real

---

## Collaboration Touchpoints

| Touchpoint | Description |
|---|---|
| **Domain Models** | Both interns depend on the `domain` package (Task 1 — already complete). Neither modifies it without review. |
| **Port Interfaces** | Intern B builds implementations; Intern A calls through port interfaces. Both must agree on interface contracts before coding. |
| **Admin API ↔ Orchestrator** | Intern A's Admin API registers apps that Intern B's Orchestrator reads. They must align on the `AppRepository` schema. |
| **Channel Ingress ↔ Orchestrator** | Intern A's channel adapters produce `UnifiedMessage`; Intern B's Orchestrator consumes it. The struct is the contract — defined and frozen from Task 1. |

---

## Summary Table

| Task | Intern A (CS + Marketing) | Intern B (Mech. Eng.) |
|---|---|---|
| Task 1: Core Domain | ✅ Complete (shared) | ✅ Complete (shared) |
| Task 2: Security Foundations | | ✅ Owner |
| Task 3: Mock Repositories | | ✅ Owner |
| Task 4: Admin API | ✅ Owner | |
| Task 5: Orchestrator Pipeline | | ✅ Owner |
| Task 6: Multi-Channel Ingress | ✅ Owner | |
| Task 7: Event Notifications | | ✅ Owner |
| Task 8: Schema Editing | ✅ Owner | |
| Task 9: Dynamic Intent Gen | | ✅ Owner |
| Task 10: Review Experiments | ✅ Owner | |
| README & Docs | ✅ Owner | |
| Unit Tests (own tasks) | ✅ Owner | ✅ Owner |
| E2E & Contract Tests | ✅ Owner | |
| Security & Integration Tests | | ✅ Owner |
