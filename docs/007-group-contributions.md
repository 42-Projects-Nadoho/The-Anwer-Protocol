# TAP: Team Management & Organisation

This document serves as the central operational charter for the TAP development team. Both members are learning Go on this project, so ownership is split by feature vertical (not by "hard core" vs "easy periphery") to keep workload and learning genuinely balanced. The shared concurrency foundation (hub/session) is built and maintained jointly throughout the week rather than handed to one owner, since it's the part both need to actually understand — not just the part one person defends in the eval.

**Default design rule:** for any lot owned by spacotto, design decisions rest with her by default — nadoho supports (short sync, no veto). This applies per-lot, not project-wide: on lots nadoho owns (protocol, CLI, world loader, and now the GUI — see below), design decisions rest with him.

## Role Assignment

| Team Member | Role | Primary Focus | File Ownership |
| :--- | :--- | :--- | :--- |
| **nadoho** | Movement, Presence & Protocol, GUI | Protocol, LOOK/MOVE/CHAT/WHO/GROUP handlers, World Loader, CLI, GUI client | `internal/protocol/`, `cmd/cli/`, `cmd/gui/`, `internal/world/loader.go` |
| **spacotto** | Items, NPCs, Combat, Quests, World Content & Test Suite | TAKE/DROP/TALK/ATTACK/STATUS/QUEST handlers, world content, world loader validation, cross-project test suite | `data/world.yaml`, `internal/world/` |
| **Both** | Server Core (hub/session/dispatcher) | Joint build & maintenance. This is where the concurrency model gets learned together | `internal/server/hub.go`, `internal/server/session.go` |

> **Rebalance note:** the GUI client moved from spacotto to nadoho. spacotto's original lot (items, combat, quests, world content, GUI, test suite) was too much to carry in one week alongside her test-suite ownership — moving the GUI, which has no real overlap with game logic, keeps her scope coherent (items/NPCs/combat/quests/world content, all "game logic + validation") without cutting into the part she's strongest at and wants to keep: testing.

## Role Descriptions & Risk Assessment

### Movement, Presence, World Infra & GUI (nadoho)

- **Core Mission**
  - Own the protocol layer and the handlers that deal with navigation and presence, the world loader that validates the map structure, and the GUI client.

- **Key Deliverables**
  - `internal/protocol`: message parsing/formatting per RFC 42TAP.
  - `handleLook`, `handleMove`, `handleChat`, `handleWho`, `handleGroup`, `handleQuit`.
  - `internal/world/loader.go`: loads YAML/JSON, validates exits/references, rejects malformed world data.
  - CLI client (`cmd/cli`).
  - GUI client (`cmd/gui`, toolkit of choice): room/inventory display, action buttons, chat/log separation, player counters.

- **Critical Challenges**
  - **Protocol compliance:** every error case in the RFC needs a matching response, not just the happy path. This is the piece the other clients (and other groups) depend on being exact.
  - **Go ramp-up:** First real project in the language, learned in parallel via the paired hub/session work.
  - **GUI scope:** async event handling while staying responsive is new ground — the CLI's connection handling (already his) should be reused rather than rebuilt.

---

### Items, NPCs, Combat, Quests, World Content & Testing (spacotto)

- **Core Mission**
  - Own the stateful game logic — items, NPC interaction, combat, quests — and the world content, plus the cross-project test suite. Design decisions on this lot are hers by default; nadoho supports.

- **Key Deliverables**
  - `handleTake`/`handleDrop` (no duplication, multi-word names), `handleTalk` (NPC dialogue).
  - Combat system: `handleAttack`/`handleStatus`, damage formulas, turn order (design synced with nadoho, implementation and calls owned solo).
  - Quest system: `handleQuest`/`handleQuests`, progression + completion validation (same sync pattern).
  - `data/world.yaml`: 8+ interconnected rooms with a loop, 3+ NPC roles, 4+ items, 2+ quests.
  - Cross-project test suite: malformed protocol input, TAKE/DROP edge cases, combat/quest logic — covering both members' handlers, not just her own.

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

> **How to read this section:** "Owner" is the person who writes the code and is responsible if it is late or broken — the single point of contact. "Support" is NOT a co-owner: it is the person who pairs briefly at a key moment (a kickoff sync, a review pass, writing test cases) so the owner is not stuck alone on a hard part. The owner still does the work; the support person just has a specific, time-boxed job on that item. On the hub/session item below, both are owners — that one is explicitly shared, not support. Design ownership defaults to whoever owns the underlying lot (spacotto for game logic, nadoho for protocol/CLI/GUI).

### Protocol & Server Skeleton (Day 1)
*   **Owners:** `nadoho` + `spacotto` (jointly)
*   **Objective:** Design the RFC message format and the hub/session/dispatcher skeleton together. This is the shared foundation everything else depends on, and it's built by both so neither is locked out of the part of the codebase most likely to need debugging under eval pressure.

### Command Dispatcher Contract
*   **Owners:** `nadoho` + `spacotto` (jointly)
*   **Objective:** Agree on the `Handler` function signature and the `Result` (direct reply + broadcasts) shape before either starts writing handlers against it.

### Combat & Quest Design
*   **Owner:** `spacotto`
*   **Support:** `nadoho`
*   **Objective:** Short sync to settle damage formulas, turn order, and quest validation rules together before spacotto implements them — this is the most design-heavy part of the mandatory scope. Design calls default to her; nadoho's role is to pressure-test the ideas, not to co-decide.

### World Loader Validation
*   **Owner:** `nadoho`
*   **Support:** `spacotto`
*   **Objective:** nadoho reviews the loader once it exists, since it consumes his protocol/error-format conventions.

### Cross-Project Test Suite
*   **Owner:** `spacotto`
*   **Support:** `nadoho`
*   **Objective:** nadoho flags the trickiest edge cases in his own modules (protocol parsing, movement handlers, GUI event handling) for her to write tests against — she owns writing and maintaining the suite.

### GUI Client
*   **Owner:** `nadoho`
*   **Support:** `spacotto` (game-state/items/combat/quest data feeds into what the GUI needs to display and expose via buttons)
*   **Objective:** Quick sync once spacotto's handlers are stable enough to expose their data shapes, so the GUI displays room/inventory/combat/quest state correctly from the start. Design calls (layout, UX) default to nadoho as owner.

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

### Movement, Presence, World Infra & GUI (nadoho)

- [ ] Implement `internal/protocol` parsing/formatting.
- [ ] Implement `handleLook`, `handleMove`, `handleChat`, `handleWho`, `handleGroup`, `handleQuit`.
- [ ] Implement `internal/world/loader.go` (validate exits/references).
- [ ] Build CLI client (`cmd/cli`) with shared connection logic.
- [ ] Build GUI client (`cmd/gui`): room view, inventory, action buttons, chat/log split, player counters.

### Items, NPCs, Combat, Quests, World Content & Testing (spacotto)

- [ ] Implement `handleTake`/`handleDrop` (no duplication, multi-word names).
- [ ] Implement `handleTalk` (NPC dialogue).
- [ ] Sync on combat/quest design with nadoho (spacotto decides).
- [ ] Implement `handleAttack`/`handleStatus` (combat system).
- [ ] Implement `handleQuest`/`handleQuests` (quest system).
- [ ] Write `data/world.yaml`: 8+ rooms with a loop, 3+ NPCs, 4+ items, 2+ quests.
- [ ] Write and maintain cross-project test suite (with edge cases flagged by nadoho).
