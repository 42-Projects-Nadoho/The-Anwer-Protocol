package server

import (
	"encoding/json"
	"strings"

	"the_answer_protocol/internal/protocol"
)

func (c *Client) handleTalk(args []string) {
	if len(args) == 0 {
		c.reply([]byte(protocol.FormatErr(protocol.ErrUnknownCommand, "USAGE: TALK <npc>")))
		return
	}
	targetName := strings.ToLower(strings.Join(args, " "))

	var found bool
	var dialogue []string
	var npcID string

	c.hub.do(func() {
		for id, dynamicNpc := range c.hub.roomNPCs[c.currentRoomID] {
			npcData, exists := c.hub.worldMap.NPCs[dynamicNpc.NPCType]
			if !exists {
				continue
			}
			if strings.ToLower(id) == targetName || strings.ToLower(npcData.Name) == targetName {
				found = true
				npcID = id
				dialogue = npcData.Dialogue
				break
			}
		}
	})

	if !found {
		c.reply([]byte(protocol.FormatErr(protocol.ErrNPCNotFound, "NPC_NOT_FOUND")))
		return
	}

	dialogueStr := ""
	if len(dialogue) > 0 {
		dialogueStr = strings.Join(dialogue, " ")
	}

	resp := map[string]string{
		"npc":      npcID,
		"dialogue": dialogueStr,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		c.reply([]byte(protocol.FormatErr(protocol.ErrSendFailed, "SERIALIZATION_FAILED")))
		return
	}
	c.reply([]byte(protocol.FormatOK(string(data))))
}

func (c *Client) handleQuest(args []string) {
	if len(args) == 0 {
		c.reply([]byte(protocol.FormatErr(protocol.ErrUnknownCommand, "USAGE: QUEST <npc>")))
		return
	}
	targetName := strings.ToLower(strings.Join(args, " "))

	var found bool
	var questGiver bool
	var questsToGive []string

	c.hub.do(func() {
		for id, dynamicNpc := range c.hub.roomNPCs[c.currentRoomID] {
			npcData, exists := c.hub.worldMap.NPCs[dynamicNpc.NPCType]
			if !exists {
				continue
			}
			if strings.ToLower(id) == targetName || strings.ToLower(npcData.Name) == targetName {
				found = true
				if npcData.Role == "quest_giver" {
					questGiver = true
					questsToGive = npcData.Quests
				}
				break
			}
		}
	})

	if !found {
		c.reply([]byte(protocol.FormatErr(protocol.ErrNPCNotFound, "NPC_NOT_FOUND")))
		return
	}
	if !questGiver || len(questsToGive) == 0 {
		c.reply([]byte(protocol.FormatErr(protocol.ErrNoQuestAvailable, "NO_QUEST_AVAILABLE")))
		return
	}

	// Just give the first quest they don't have yet.
	var questID string
	c.hub.do(func() {
		for _, q := range questsToGive {
			has := false
			for _, active := range c.quests {
				if active == q || active == q+"_completed" {
					has = true
					break
				}
			}
			if !has {
				questID = q
				c.quests = append(c.quests, q)
				
				c.hub.logger.Info("quest_accepted",
					"username", c.username,
					"quest_id", questID,
				)
				break
			}
		}
	})

	if questID == "" {
		c.reply([]byte(protocol.FormatErr(protocol.ErrNoQuestAvailable, "NO_QUEST_AVAILABLE")))
		return
	}

	c.reply([]byte(protocol.FormatOK("quest_accepted=" + questID)))
}

func (c *Client) handleQuests() {
	var activeQuests []string
	
	c.hub.do(func() {
		// First pass: check for completion
		for i, q := range c.quests {
			if strings.HasSuffix(q, "_completed") {
				continue
			}
			if strings.HasPrefix(q, "defeated:") {
				continue
			}
			questData, exists := c.hub.worldMap.Quests[q]
			if !exists {
				continue
			}
			
			completed := false
			if questData.Type == "fetch" {
				// check inventory
				for _, item := range c.inventory {
					if item == questData.Target {
						completed = true
						break
					}
				}
			} else if questData.Type == "defeat" {
				// check history
				for _, hist := range c.quests {
					if hist == "defeated:"+questData.Target {
						completed = true
						break
					}
				}
			}

			if completed {
				c.quests[i] = q + "_completed"
				
				c.hub.logger.Info("quest_completed",
					"username", c.username,
					"quest_id", q,
				)
				
				// Send a reward message
				msg := "You completed a quest: " + questData.Name + "! " + questData.Reward
				select {
				case c.send <- []byte(protocol.FormatEvt("GLOBAL", "CHAT", "QuestSys "+msg)):
				default:
				}
			}
		}

		// Second pass: gather names
		for _, q := range c.quests {
			if strings.HasPrefix(q, "defeated:") {
				continue
			}
			if strings.HasSuffix(q, "_completed") {
				base := strings.TrimSuffix(q, "_completed")
				if qData, ok := c.hub.worldMap.Quests[base]; ok {
					activeQuests = append(activeQuests, "[DONE] "+qData.Name)
				}
			} else {
				if qData, ok := c.hub.worldMap.Quests[q]; ok {
					activeQuests = append(activeQuests, "[ACTIVE] "+qData.Name)
				}
			}
		}
	})

	data, err := json.Marshal(activeQuests)
	if err != nil {
		c.reply([]byte(protocol.FormatErr(protocol.ErrSendFailed, "SERIALIZATION_FAILED")))
		return
	}
	c.reply([]byte(protocol.FormatOK(string(data))))
}
