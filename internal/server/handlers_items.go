package server

import (
	"encoding/json"
	"strings"

	"the_answer_protocol/internal/protocol"
)

func (c *Client) handleTake(args []string) {
	if len(args) == 0 {
		c.reply([]byte(protocol.FormatErr(protocol.ErrUnknownCommand, "USAGE: TAKE <item>")))
		return
	}
	targetItem := strings.ToLower(strings.Join(args, " "))

	// Resolve the item by ID or name
	var itemID string
	var found bool
	c.hub.do(func() {
		for id, present := range c.hub.roomItems[c.currentRoomID] {
			if !present {
				continue
			}
			itemData, exists := c.hub.worldMap.Items[id]
			if !exists {
				continue
			}
			if strings.ToLower(id) == targetItem || strings.ToLower(itemData.Name) == targetItem {
				itemID = id
				found = true
				break
			}
		}

		if found {
			c.hub.roomItems[c.currentRoomID][itemID] = false
			c.inventory = append(c.inventory, itemID)

			c.hub.logger.Info("item_taken",
				"username", c.username,
				"room_id", c.currentRoomID,
				"item_id", itemID,
			)
		}
	})

	if !found {
		c.reply([]byte(protocol.FormatErr(protocol.ErrItemNotFound, "ITEM_NOT_FOUND")))
		return
	}

	c.reply([]byte(protocol.FormatOK("taken=" + itemID)))

	// Check if this fulfills a fetch quest
	c.checkQuestCompletion()
}

func (c *Client) handleDrop(args []string) {
	if len(args) == 0 {
		c.reply([]byte(protocol.FormatErr(protocol.ErrUnknownCommand, "USAGE: DROP <item>")))
		return
	}
	targetItem := strings.ToLower(strings.Join(args, " "))

	var itemID string
	var foundIdx = -1
	
	c.hub.do(func() {
		for i, id := range c.inventory {
			itemData, exists := c.hub.worldMap.Items[id]
			if !exists {
				continue
			}
			if strings.ToLower(id) == targetItem || strings.ToLower(itemData.Name) == targetItem {
				itemID = id
				foundIdx = i
				break
			}
		}

		if foundIdx != -1 {
			// Remove from inventory
			c.inventory = append(c.inventory[:foundIdx], c.inventory[foundIdx+1:]...)
			// Add back to room
			if c.hub.roomItems[c.currentRoomID] == nil {
				c.hub.roomItems[c.currentRoomID] = make(map[string]bool)
			}
			c.hub.roomItems[c.currentRoomID][itemID] = true

			c.hub.logger.Info("item_dropped",
				"username", c.username,
				"room_id", c.currentRoomID,
				"item_id", itemID,
			)
		}
	})

	if foundIdx == -1 {
		c.reply([]byte(protocol.FormatErr(protocol.ErrItemNotInInventory, "ITEM_NOT_IN_INVENTORY")))
		return
	}

	c.reply([]byte(protocol.FormatOK("dropped=" + itemID)))
}

func (c *Client) handleInventory() {
	var inv []string
	c.hub.do(func() {
		inv = append([]string{}, c.inventory...)
	})

	data, err := json.Marshal(inv)
	if err != nil {
		c.reply([]byte(protocol.FormatErr(protocol.ErrSendFailed, "SEND_FAILED")))
		return
	}

	c.reply([]byte(protocol.FormatOK(string(data))))
}
