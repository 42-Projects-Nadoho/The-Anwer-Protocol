package server

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"the_answer_protocol/internal/protocol"
)

func (c *Client) handleAttack(args []string) {
	if len(args) == 0 {
		c.reply([]byte(protocol.FormatErr(protocol.ErrUnknownCommand, "USAGE: ATTACK <npc>")))
		return
	}
	targetName := strings.ToLower(strings.Join(args, " "))

	var npcID string
	var npcType string
	var found bool
	var hostile bool
	var dmgDealt int = 15 // Base player damage
	var enemyHp int
	var counterDmg int = 0
	var enemyDied bool

	c.hub.do(func() {
		// Find NPC in room
		for id, dynamicNpc := range c.hub.roomNPCs[c.currentRoomID] {
			npcData, exists := c.hub.worldMap.NPCs[dynamicNpc.NPCType]
			if !exists {
				continue
			}
			if strings.ToLower(id) == targetName || strings.ToLower(npcData.Name) == targetName {
				npcID = id
				npcType = dynamicNpc.NPCType
				found = true
				if npcData.Role == "enemy" {
					hostile = true
					// Deal damage
					dynamicNpc.HP -= dmgDealt
					enemyHp = dynamicNpc.HP
					if dynamicNpc.HP <= 0 {
						enemyDied = true
						// Delete from room
						delete(c.hub.roomNPCs[c.currentRoomID], id)
						
						c.hub.logger.Info("npc_defeated",
							"username", c.username,
							"npc_id", id,
							"room_id", c.currentRoomID,
						)

						// Schedule respawn (e.g. 30 seconds)
						respawnID := id
						respawnRoom := c.currentRoomID
						respawnMaxHP := npcData.Stats.HP
						respawnType := npcType
						
						time.AfterFunc(30*time.Second, func() {
							c.hub.do(func() {
								if c.hub.roomNPCs[respawnRoom] == nil {
									c.hub.roomNPCs[respawnRoom] = make(map[string]*DynamicNPC)
								}
								c.hub.roomNPCs[respawnRoom][respawnID] = &DynamicNPC{
									ID:      respawnID,
									NPCType: respawnType,
									HP:      respawnMaxHP,
								}
							})
							// Notify the room that the enemy respawned
							c.hub.BroadcastRoom(respawnRoom, []byte(protocol.FormatEvt("ROOM", "RESPAWN", "The air shifts... "+respawnID+" has respawned!")), nil)
						})
					} else {
						// Counter attack
						counterDmg = npcData.Stats.Damage
						c.hp -= counterDmg
					}
				}
				break
			}
		}
	})

	if !found {
		c.reply([]byte(protocol.FormatErr(protocol.ErrNPCNotFound, "NPC_NOT_FOUND")))
		return
	}
	if !hostile {
		c.reply([]byte(protocol.FormatErr(protocol.ErrNPCNotHostile, "NPC_NOT_HOSTILE")))
		return
	}

	// Broadcast combat logs
	if enemyDied {
		msg := fmt.Sprintf("%s dealt %d damage to %s. The enemy is defeated!", c.username, dmgDealt, npcID)
		c.hub.BroadcastRoom(c.currentRoomID, []byte(protocol.FormatEvt("ROOM", "COMBAT", msg)), nil)
		
		// Broadcast custom combat event for client parsing
		customEvt := fmt.Sprintf("DEFEAT %s %s", c.username, npcID)
		c.hub.BroadcastRoom(c.currentRoomID, []byte(protocol.FormatEvt("ROOM", "COMBAT", customEvt)), nil)

		
		var php int
		c.hub.do(func() { php = c.hp })
		resp := map[string]interface{}{
			"attacker_hp": php,
			"target_hp":   0,
			"damage":      dmgDealt,
			"status":      "victory",
		}
		data, _ := json.Marshal(resp)
		c.reply([]byte(protocol.FormatOK(string(data))))
		
		// If player had a defeat quest for this, mark it (simple logic)
		c.hub.do(func() {
			c.quests = append(c.quests, "defeated:"+npcType)
		})
		c.checkQuestCompletion()
	} else {
		msg := fmt.Sprintf("%s dealt %d damage to %s. %s has %d HP left. %s counter-attacked for %d damage!", 
			c.username, dmgDealt, npcID, npcID, enemyHp, npcID, counterDmg)
		c.hub.BroadcastRoom(c.currentRoomID, []byte(protocol.FormatEvt("ROOM", "COMBAT", msg)), nil)
		
		var php int
		c.hub.do(func() { php = c.hp })
		resp := map[string]interface{}{
			"attacker_hp": php,
			"target_hp":   enemyHp,
			"damage":      dmgDealt,
			"status":      "combat",
		}
		data, _ := json.Marshal(resp)
		c.reply([]byte(protocol.FormatOK(string(data))))

		// Check if player died from counter-attack
		var died bool
		c.hub.do(func() {
			if c.hp <= 0 {
				died = true
				c.hp = 100 // Respawn with full health (or reduced)
				
				// Move to start room
				oldRoom := c.currentRoomID
				c.currentRoomID = c.hub.worldMap.StartRoomID
				
				c.hub.logger.Info("player_died",
					"username", c.username,
					"killed_by", npcID,
					"room_id", oldRoom,
				)
			}
		})

		if died {
			deathMsg := fmt.Sprintf("%s has been defeated and sent back to safety.", c.username)
			c.hub.BroadcastGlobal([]byte(protocol.FormatEvt("GLOBAL", "CHAT", "CombatSys "+deathMsg)))
			c.reply([]byte(protocol.FormatEvt("ROOM", "CHAT", "CombatSys You died! Respawning...")))
			// Need to notify player of room change
			c.reply([]byte(protocol.FormatOK("room=" + c.hub.worldMap.StartRoomID)))
		}
	}
}

func (c *Client) handleStatus() {
	var hp int
	c.hub.do(func() {
		hp = c.hp
	})
	status := "healthy"
	if hp <= 0 {
		status = "dead"
	}
	resp := map[string]interface{}{
		"hp":     hp,
		"max_hp": 100,
		"status": status,
	}
	data, _ := json.Marshal(resp)
	c.reply([]byte(protocol.FormatOK(string(data))))
}
