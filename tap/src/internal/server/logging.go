package server

import "strings"

// reply sends msg to the client and logs it, satisfying the requirement to
// log every server response and error code sent to clients. ERR replies
// are logged at Warn level, everything else at Info.
func (c *Client) reply(msg []byte) {
	c.send <- msg

	text := strings.TrimSuffix(string(msg), "\n")
	if strings.HasPrefix(text, "ERR") {
		c.hub.logger.Warn("response_sent", "username", c.username, "response", text)
		return
	}
	c.hub.logger.Info("response_sent", "username", c.username, "response", text)
}
