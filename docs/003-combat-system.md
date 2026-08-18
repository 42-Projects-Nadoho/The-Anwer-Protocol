# Combat System

The Answer Protocol features a streamlined, real-time combat system. Rather than relying on a complex, slow, turn-based initiative order, the combat is completely synchronous.

## Mechanics

### Attacking
Players initiate combat by sending the `ATTACK <npc>` command.
1. The server verifies the NPC exists in the player's current room and is marked with the `enemy` role.
2. A flat 15 base damage is instantly dealt to the NPC.
3. If the NPC's HP drops to 0, it dies, and the server broadcasts a custom `DEFEAT <username> <npc_id>` event to the room, followed by scheduling the NPC to respawn after 30 seconds.
4. If the NPC survives, it immediately counter-attacks, dealing its stat-based damage back to the player.

### Death and Respawning
Permanent death is not implemented to ensure continuous gameplay flow.
- When a player's HP drops to 0 from an NPC counter-attack, they are "defeated".
- The server instantly restores their HP to 50 and teleports them to the `StartRoomID` (Destiny Islands).
- A `GLOBAL CHAT` message is broadcasted announcing their defeat.
- Appropriate `PRESENCE LEAVE` and `PRESENCE ENTER` events are broadcasted to update the world state for all clients seamlessly.

### Client Synchronization
All combat actions return a structured JSON response containing the updated `attacker_hp`, `target_hp`, `damage`, and combat `status` (`combat` or `victory`). This allows clients (especially the GUI) to instantly render health bars and damage numbers without needing to poll the server.
