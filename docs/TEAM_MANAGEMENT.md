# TAP: Team Management & Organisation

This document serves as the central operational charter for the TAP development team. Both members are learning Go on this project, so ownership is split by feature vertical (not by "hard core" vs "easy periphery") to keep workload and learning genuinely balanced. The shared concurrency foundation (hub/session) is built and maintained jointly throughout the week rather than handed to one owner, since it's the part both need to actually understand — not just the part one person defends in the eval.

## Role Assignment

| Team Member | Role | Primary Focus | File Ownership |
| :--- | :--- | :--- | :--- |
| **nadoho** | Movement, Presence & Protocol | Protocol, LOOK/MOVE/CHAT/WHO/GROUP handlers, CLI | `internal/protocol/`, `cmd/cli/` |
| **spacotto** | Items, NPCs, Combat, Quests, World Validation & Test Suite | TAKE/DROP/TALK/ATTACK/STATUS/QUEST handlers, world content, world loader validation, cross-project test suite, GUI | `data/world.yaml`, `internal/world/`, `cmd/gui/` |
| **Both** | Server Core (hub/session/dispatcher) | Joint build & maintenance — this is where the concurrency model gets learned together | `internal/server/hub.go`, `internal/server/session.go` |

## Role Descriptions & Risk Assessment

### Movement, Presence & World Infra (nadoho)

- **Core Mission**
  - Own the protocol layer and the handlers that deal with navigation and presence, plus the world loader that validates the map structure.
- **Key Deliverables**
  - `internal/protocol`: message parsing/formatting per RFC 42TAP.
  - `handleLook`, `handleMove`, `handleChat`, `handleWho`, `handleGroup`, `handleQuit`.
  - CLI client (`cmd/cli`).
- **Critical Challenges**
  - **Protocol compliance:** every error case in the RFC needs a matching response, not just the happy path — this is the piece the other clients (and other groups) depend on being exact.
  - **Go ramp-up:** same as spacotto — first real project in the language, learned in parallel via the paired hub/session work.

---

### Items, NPCs, Combat, Quests, World Validation & Testing (spacotto)

- **Core Mission**
  - Own the stateful game logic — items, NPC interaction, combat, quests — and the world content, plus the GUI.
- **Key Deliverables**
  - `handleTake`/`handleDrop` (no duplication, multi-word names), `handleTalk` (NPC dialogue).
  - Combat system: `handleAttack`/`handleStatus`, damage formulas, turn order (design synced with nadoho, implementation owned solo).
  - Quest system: `handleQuest`/`handleQuests`, progression + completion validation (same sync pattern).
  - `data/world.yaml`: 8+ interconnected rooms with a loop, 3+ NPC roles, 4+ items, 2+ quests.
  - `internal/world/loader.go`: loads YAML/JSON, validates exits/references, rejects malformed world data.
  - Cross-project test suite: malformed protocol input, TAKE/DROP edge cases, combat/quest logic — covering both members' handlers, not just her own.
  - GUI client (`cmd/gui`, Fyne): room/inventory display, action buttons, chat/log separation, player counters.
- **Critical Challenges**
  - **Design surface:** combat and quest systems are the most open-ended part of the mandatory scope — worth a short design sync before coding, but the implementation and the calls made are hers to own.
  - **Item/quest state correctness:** race-free TAKE/DROP against the shared hub.

---

### Server Core: hub.go / session.go (nadoho + spacotto, jointly)

- **Core Mission**
  - Build and maintain the TCP accept loop, one goroutine per client session, and the single-writer hub that owns world state and broadcasts events — the trickiest and most bug-prone part of the project, so it's built as a shared skill rather than one person's specialty.
- **Key Deliverables**
  - Command routing via channels, no shared-state mutexes.
  - Graceful disconnect handling and broadcast-without-interruption.
  - Structured (JSON) server logging.
- **Critical Challenges**
  - **Concurrency correctness:** races and goroutine leaks are the hardest bugs to catch under deadline — the single-writer hub pattern exists specifically to avoid this, and both need to be able to debug it under eval pressure, not just one.

## Cross-Support & Shared Ownership

> **How to read this section:** "Owner" is the person who writes the code and is responsible if it is late or broken — the single point of contact. "Support" is NOT a co-owner: it is the person who pairs briefly at a key moment (a kickoff sync, a review pass, writing test cases) so the owner is not stuck alone on a hard part. The owner still does the work; the support person just has a specific, time-boxed job on that item. On the hub/session item below, both are owners — that one is explicitly shared, not support.

### Protocol & Server Skeleton (Day 1)
*   **Owners:** `nadoho` + `spacotto` (jointly)
*   **Objective:** Design the RFC message format and the hub/session/dispatcher skeleton together. This is the shared foundation everything else depends on, and it's built by both so neither is locked out of the part of the codebase most likely to need debugging under eval pressure.

### Command Dispatcher Contract
*   **Owners:** `nadoho` + `spacotto` (jointly)
*   **Objective:** Agree on the `Handler` function signature and the `Result` (direct reply + broadcasts) shape before either starts writing handlers against it.

### Combat & Quest Design
*   **Owner:** `spacotto`
*   **Support:** `nadoho`
*   **Objective:** Short sync to settle damage formulas, turn order, and quest validation rules together before spacotto implements them — this is the most design-heavy part of the mandatory scope, and a joint decision even though spacotto writes the code.

### World Loader Validation
*   **Owner:** `spacotto`
*   **Support:** `nadoho`
*   **Objective:** nadoho reviews the loader once it exists, since it consumes his protocol/error-format conventions.

### Cross-Project Test Suite
*   **Owner:** `spacotto`
*   **Support:** `nadoho`
*   **Objective:** nadoho flags the trickiest edge cases in his own modules (protocol parsing, movement handlers) for her to write tests against — she owns writing and maintaining the suite.

### GUI ↔ Client Connection Layer
*   **Owner:** `spacotto` (GUI itself)
*   **Support:** `nadoho` (protocol/CLI experience feeds into what the connection layer needs to expose)
*   **Objective:** Quick sync once the CLI's connection handling is stable, so the GUI doesn't reinvent message parsing from scratch.

### README.md
*   **Owner:** Each member writes their own section
*   **Support:** the other member
*   **Objective:** Mutual proofread and consistency pass across all required sections.

## Tasks

> Update as you go:
> - [ ] To do
> - [x] Done

### Shared: Server Core (nadoho + spacotto)

- [ ] Define protocol message format together (Day 1).
- [ ] Design hub.go: single-writer world state owner, command channel, broadcast logic.
- [ ] Design session.go: per-client goroutine, read loop, reply/event channels.
- [ ] Agree on dispatcher `Handler` signature and `Result` shape.
- [ ] Implement graceful disconnect + leave-event broadcast.
- [ ] Implement structured (JSON) server logging.

### Movement, Presence & World Infra (nadoho)

- [ ] Implement `internal/protocol` parsing/formatting.
- [ ] Implement `handleLook`, `handleMove`, `handleChat`, `handleWho`, `handleGroup`, `handleQuit`.
- [ ] Implement `internal/world/loader.go` (validate exits/references).
- [ ] Build CLI client (`cmd/cli`) with shared connection logic.

### Items, NPCs, Combat, Quests, World Validation & Testing (spacotto)

- [ ] Implement `handleTake`/`handleDrop` (no duplication, multi-word names).
- [ ] Implement `handleTalk` (NPC dialogue).
- [ ] Sync on combat/quest design with nadoho.
- [ ] Implement `handleAttack`/`handleStatus` (combat system).
- [ ] Implement `handleQuest`/`handleQuests` (quest system).
- [ ] Write `data/world.yaml`: 8+ rooms with a loop, 3+ NPCs, 4+ items, 2+ quests.
- [ ] Build GUI client (`cmd/gui`, Fyne): room view, inventory, action buttons, chat/log split, player counters.
