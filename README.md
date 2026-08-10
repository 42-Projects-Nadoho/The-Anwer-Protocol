*This project has been created as part of the 42 curriculum by nadoho, spacotto.*


![TAP_KH_Images](/docs/images/image_for_tap.jpg)
## Description

TAP (The Answer Protocol) is a small shared-world multiplayer text
adventure: a TCP server hosting a persistent-feeling world, with a CLI
client and a GUI client speaking the same line-based protocol (RFC 42TAP).

TODO: expand with goal + brief overview.

## Instructions

See "Building and Running" below.

## Resources

TODO: classic references (MUD design, Go concurrency patterns, etc.)
and a description of how AI was used, specifying which tasks/parts.

## Architecture

```
tap/
├── go.mod
├── Makefile
├── cmd/
│   ├── server/main.go
│   ├── cli/main.go
│   └── gui/main.go
├── internal/
│   ├── protocol/
│   │   ├── message.go     # Command/Response types, OK/ERR/EVT formatting
│   │   ├── parse.go       # line -> Command
│   │   └── errors.go      # RFC error codes + local extensions
│   ├── world/
│   │   └── loader.go      # YAML loading + validation (exits, duplicate ids)
│   └── server/
│       ├── hub.go         # single-writer state owner (clients, groups, world)
│       ├── session.go     # per-connection goroutines, command dispatch
│       └── handlers.go    # one handleX per RFC command
├── data/
│   └── world.yaml
└── README.md
```

**Concurrency model.** Each client connection runs two goroutines:
`readPump` (parses incoming lines, dispatches to handlers) and `writePump`
(drains that client's outgoing channel to the socket). Command dispatch is
done inline in `readPump`'s switch statement rather than through a separate
dispatcher — the subject explicitly allows either approach, and the command
set is small enough that a switch stays readable.

All state shared across clients (connected usernames, room membership,
group membership, next-group-id counter) lives on `Hub` and is only ever
touched from inside `Hub.Run()`'s goroutine. Rather than guarding that state
with a mutex, other goroutines submit a closure through `Hub.do(fn func())`,
which is executed synchronously inside `Run()` and blocks the caller until
done. This keeps the "single writer" invariant (no two goroutines ever read
or write `Hub`'s maps concurrently) without scattering `sync.Mutex` calls
through the handler code. Broadcasting to a client's own `send` channel from
inside a `Hub.do` closure always uses a non-blocking `select { case: default:
}` — a blocking send there would stall the entire hub if that one client's
buffer were ever full.

## Protocol Implementation

TAP follows RFC 42TAP for all documented commands, responses, and events.
Where the RFC is silent or ambiguous, the choices below were made and are
documented here as required.

**Local error codes (not in RFC 42TAP):**

| Code | Message | Meaning |
|---|---|---|
| 902 | `NOT_AUTHENTICATED` | Any command other than `CONNECT`/`QUIT` sent before authentication completes. |
| 903 | `UNKNOWN_COMMAND` | Command name not recognized, or malformed usage (e.g. missing required argument). |

**GROUP: no privileged "leader" role.** The RFC's `GROUP JOIN <leader-name>`
syntax and `EVT GROUP INVITE <leader>` event both name their argument
"leader", but the RFC never defines what a leader is or what privileges it
has — checked directly against the RFC source, this is the only place the
word appears. We deliberately did not build a leader role: any current
member of a group may invite others, and a group persists as long as it has
at least one member, regardless of who created it (there is no special
"leader leaves -> group dissolves" behavior).

**GROUP JOIN takes a group id, not a username — a deliberate RFC deviation.**
The RFC specifies `GROUP JOIN <leader-name>` (a player's username). Because
our groups have no leader role and can outlive their creator, resolving
"join by naming a member" breaks as soon as the named member has since left
or disconnected, even though the group itself may still exist. We instead
require the server-generated opaque group id returned by `GROUP CREATE`
(e.g. `group-3`). **This breaks wire compatibility with a strictly
RFC-literal client from another group for this one command** — a conscious
tradeoff, not an oversight.

**`EVT GROUP INVITE` carries the group id in addition to the inviter.**
Format: `EVT GROUP INVITE <inviter-username> <group-id>` instead of the
RFC's `EVT GROUP INVITE <leader>`. This is a direct consequence of the
previous point: since `JOIN` needs a group id, the invitee has to learn it
from somewhere, and the invite event is the only place that can carry it.

**Invitations are single-use.** A pending invitation (tracked per-client in
`isInvited []string`) is consumed the moment `GROUP JOIN` succeeds against
it; re-joining after leaving requires a fresh invite.

**`LOOK`'s `items`/`npcs` fields are always empty arrays for now.** The JSON
shape matches the RFC exactly; the arrays will be populated once the
item/NPC system (spacotto's scope) lands.

## World Design

TODO (spacotto): room layout, NPC roles, item distribution. Current
`data/world.yaml` is a 2-room placeholder well below the required minimum
(8+ interconnected rooms with a loop, 3+ NPC roles, 4+ items, 2+ quests).

## Combat System

TODO (spacotto): turn-based mechanics, damage formulas, initiative order,
additional commands (`DEFEND`, `FLEE`, ...), justification for the design
choices above.

## Quest System

TODO (spacotto): quest progression mechanics, completion validation, reward
distribution, quest dependencies.

## Server Logging

TODO: structured (JSON recommended) logging of connections/disconnections,
commands received, responses sent, world state changes, quest events, and
abuse-pattern monitoring (command flooding, rapid connections). Not yet
implemented.

## Group Contributions

| Member | Responsibilities |
|---|---|
| nadoho | Protocol (`internal/protocol/`), server core (`internal/server/hub.go`, `session.go`), command handlers for `CONNECT`/`LOOK`/`MOVE`/`CHAT`/`WHO`/`GROUP`/`QUIT`, world loader validation, CLI client, GUI client, Makefile. |
| spacotto | Items, NPCs, combat, quests, world content (`data/world.yaml`), cross-project test suite. |

See `docs/TEAM_MANAGEMENT.md` for the full role breakdown and rationale.

## Building and Running

Requires Go 1.22+.

```
make deps
make run-server        # in one terminal
make run-client         # in another
make run-client-gui     # or this, instead of/alongside the CLI
```

## Testing

```
make test
```

TODO: document how to exercise multiplayer (two `run-client` sessions),
combat, and quest flows manually or via the test suite.
