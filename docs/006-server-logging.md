# Server Logging
The Answer Protocol server leverages Go's native `log/slog` package for highly structured, JSON-formatted logging. This approach ensures that logs are not only human-readable but also easily ingestible by modern log aggregation tools for monitoring server health, tracking player activity, and preventing abuse.

## Implementation Details
The core of our response logging is handled by a centralized `reply` wrapper in `tap/src/internal/server/logging.go`:

```go
func (c *Client) reply(msg []byte) {
	c.send <- msg

	text := strings.TrimSuffix(string(msg), "\n")
	if strings.HasPrefix(text, "ERR") {
		c.hub.logger.Warn("response_sent", "username", c.username, "response", text)
		return
	}
	c.hub.logger.Info("response_sent", "username", c.username, "response", text)
}
```

### Log Levels
We employ strict log-level separation to facilitate debugging:

| Level | Usage |
| :--- | :--- |
| **`Info`** | Used for standard operational events such as successful client connections (`CONNECT`), player movement, standard game events (`ROOM CHAT`, `ROOM COMBAT`), and successful server responses (`OK`). |
| **`Warn`** | Used whenever the server dispatches an `ERR` response to a client (e.g., malformed syntax, unknown commands). Spikes in `Warn` logs associated with a specific IP or username clearly indicate abuse attempts or brute-force flooding. |
| **`Error`** | Reserved for critical subsystem failures, such as the inability to bind the TCP port or fatal YAML parsing errors during world initialization. |

### Structured Metadata
Every log entry is enriched with key-value metadata. Rather than parsing raw text strings, administrators can filter logs by fields such as:

| Field | Description |
| :--- | :--- |
| `username` | The acting player. |
| `room_id` | Where the event occurred. |
| `npc_id` | Which NPC was targeted during combat. |
| `response` | The raw ABNF protocol string dispatched. |
