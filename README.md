*This project has been created as part of the 42 curriculum by nadoho, spacotto.*

## Description

TAP (The Answer Protocol) is a small shared-world multiplayer text
adventure: a TCP server hosting a persistent-feeling world, with a CLI
client and a GUI client speaking the same line-based protocol.

TODO: expand with goal + brief overview.

## Instructions

See "Building and Running" below.

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
│   │   ├── message.go
│   │   ├── parse.go
│   │   └── errors.go
│   ├── world/
│   │   ├── room.go
│   │   ├── item.go
│   │   ├── npc.go
│   │   ├── quest.go
│   │   └── loader.go
│   ├── server/
│   │   ├── server.go
│   │   ├── session.go
│   │   ├── hub.go
│   │   ├── dispatcher.go
│   │   ├── handlers.go
│   │   ├── combat.go
│   │   └── logging.go
│   └── client/
│       └── conn.go
├── data/
│   └── world.yaml
└── README.md
```

## Group Contributions


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

## Resources

TODO: classic references (MUD design, Go concurrency patterns, etc.)
and a description of how AI was used, specifying which tasks/parts.
