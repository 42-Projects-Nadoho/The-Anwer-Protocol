package server

import (
	"encoding/json"
	"strings"

	"the_answer_protocol/tap/src/internal/protocol"
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
	var npcType string

	c.hub.do(func() {
		for id, dynamicNpc := range c.hub.roomNPCs[c.currentRoomID] {
			npcData, exists := c.hub.worldMap.NPCs[dynamicNpc.NPCType]
			if !exists {
				continue
			}
			if strings.ToLower(id) == targetName || strings.ToLower(npcData.Name) == targetName {
				found = true
				npcID = id
				npcType = dynamicNpc.NPCType
				dialogue = npcData.Dialogue
				break
			}
		}
	})

	if !found {
		c.reply([]byte(protocol.FormatErr(protocol.ErrNPCNotFound, "NPC_NOT_FOUND")))
		return
	}

	var questUnlocked string
	var healToFull bool
	var completedQuest string
	c.hub.do(func() {
		if npcType == "chest" {
			hasQuest := false
			for _, q := range c.quests {
				if q == "find_wayfinder" {
					hasQuest = true
					break
				}
			}
			if hasQuest {
				dialogue = []string{"You open the chest and find the Wayfinder! Return to Yen Sid."}
				for i, q := range c.quests {
					if q == "find_wayfinder" {
						c.quests[i] = "find_wayfinder_found"
						break
					}
				}
			} else {
				dialogue = []string{"It's locked tight."}
			}
		} else if npcType == "yen_sid" {
			hasFindWayfinderActive := false
			hasFoundWayfinder := false
			hasDefeatShadowActive := false
			hasFoundShadow := false
			isAllCompleted := false

			for _, q := range c.quests {
				if q == "find_wayfinder" {
					hasFindWayfinderActive = true
				} else if q == "find_wayfinder_found" {
					hasFoundWayfinder = true
				} else if q == "defeat_shadow" {
					hasDefeatShadowActive = true
				} else if q == "defeat_shadow_found" {
					hasFoundShadow = true
				} else if q == "defeat_shadow_completed" {
					isAllCompleted = true
				}
			}
			
			if hasFoundWayfinder {
				dialogue = []string{"Ah, you found it! Now, you must push back the darkness."}
				questUnlocked = "defeat_shadow"
				for i, q := range c.quests {
					if q == "find_wayfinder_found" {
						c.quests[i] = "find_wayfinder_completed"
						break
					}
				}
				// Give them the next quest automatically
				hasNext := false
				for _, q := range c.quests {
					if q == "defeat_shadow" || q == "defeat_shadow_completed" || q == "defeat_shadow_found" {
						hasNext = true
						break
					}
				}
				if !hasNext {
					c.quests = append(c.quests, "defeat_shadow")
				}
			} else if hasFoundShadow {
				dialogue = []string{"You have proven your strength. Your heart is fully restored!"}
				healToFull = true
				completedQuest = "defeat_shadow"
				for i, q := range c.quests {
					if q == "defeat_shadow_found" {
						c.quests[i] = "defeat_shadow_completed"
						break
					}
				}
			} else if isAllCompleted {
				dialogue = []string{"You have done well. The islands are safe for now."}
			} else if hasDefeatShadowActive {
				dialogue = []string{"The Heartless still linger. You must defeat the Shadow Heartless."}
			} else if hasFindWayfinderActive {
				dialogue = []string{"What are you waiting for? Find the chest in the Secret Cave!"}
			}
		}
	})

	if questUnlocked != "" {
		msg := "You completed a quest: Bonds of Friendship! Unlocked new quest: Push Back the Darkness!"
		c.hub.BroadcastGlobal([]byte(protocol.FormatEvt("GLOBAL", "CHAT", "QuestSys "+msg)))
	}

	if completedQuest != "" {
		if healToFull {
			c.hub.do(func() { c.hp = 100 })
		}
		var qName, qReward string
		c.hub.do(func() { 
			qData := c.hub.worldMap.Quests[completedQuest]
			qName = qData.Name
			qReward = qData.Reward
		})
		msg := "You completed a quest: " + qName + "! " + qReward
		c.hub.BroadcastGlobal([]byte(protocol.FormatEvt("GLOBAL", "CHAT", "QuestSys "+msg)))
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
		c.reply([]byte(protocol.FormatErr(protocol.ErrSendFailed, "SEND_FAILED")))
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
				if active == q || active == q+"_completed" || active == q+"_found" {
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

	var qDesc, qReward string
	c.hub.do(func() {
		qData := c.hub.worldMap.Quests[questID]
		qDesc = qData.Goal
		qReward = qData.Reward
	})

	resp := map[string]interface{}{
		"quest_id":    questID,
		"description": qDesc,
		"reward":      qReward,
		"status":      "available",
	}
	data, _ := json.Marshal(resp)
	c.reply([]byte(protocol.FormatOK(string(data))))
}

func (c *Client) handleQuests() {
	activeQuests := make([]map[string]interface{}, 0)
	
	c.checkQuestCompletion()

	c.hub.do(func() {
		for _, q := range c.quests {
			if strings.HasPrefix(q, "defeated:") {
				continue
			}
			if strings.HasSuffix(q, "_completed") {
				base := strings.TrimSuffix(q, "_completed")
				if _, ok := c.hub.worldMap.Quests[base]; ok {
					activeQuests = append(activeQuests, map[string]interface{}{
						"quest_id": base,
						"status":   "completed",
					})
				}
			} else if strings.HasSuffix(q, "_found") {
				base := strings.TrimSuffix(q, "_found")
				if _, ok := c.hub.worldMap.Quests[base]; ok {
					activeQuests = append(activeQuests, map[string]interface{}{
						"quest_id": base,
						"status":   "return to npc",
						"progress": "1/1",
					})
				}
			} else {
				if _, ok := c.hub.worldMap.Quests[q]; ok {
					activeQuests = append(activeQuests, map[string]interface{}{
						"quest_id": q,
						"status":   "active",
						"progress": "0/1",
					})
				}
			}
		}
	})

	data, err := json.Marshal(activeQuests)
	if err != nil {
		c.reply([]byte(protocol.FormatErr(protocol.ErrSendFailed, "SEND_FAILED")))
		return
	}
	c.reply([]byte(protocol.FormatOK(string(data))))
}

func (c *Client) checkQuestCompletion() {
	c.hub.do(func() {
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
				for _, item := range c.inventory {
					if item == questData.Target {
						completed = true
						break
					}
				}
			} else if questData.Type == "defeat" {
				for _, hist := range c.quests {
					if hist == "defeated:"+questData.Target {
						completed = true
						break
					}
				}
			}

			if completed {
				if q == "defeat_shadow" {
					c.quests[i] = q + "_found"
				} else {
					c.quests[i] = q + "_completed"
					c.hub.logger.Info("quest_completed",
						"username", c.username,
						"quest_id", q,
					)
					msg := "You completed a quest: " + questData.Name + "! " + questData.Reward
					select {
					case c.send <- []byte(protocol.FormatEvt("GLOBAL", "CHAT", "QuestSys "+msg)):
					default:
					}
				}
			}
		}
	})
}
