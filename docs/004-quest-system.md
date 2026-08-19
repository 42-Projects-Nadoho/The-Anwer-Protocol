# Quest System
The Answer Protocol features a robust backend quest engine implemented in `tap/src/internal/server/handlers_quests.go`. It manages quest state directly on the server to prevent client-side manipulation.

## Quest Mechanics

### Acquisition
Quests are defined within `tap/data/world.yaml`. Players can initiate quests by approaching NPCs assigned the `quest_giver` role (e.g., Master Yen Sid).
Using the `QUEST <npc>` command, the server will check if the player already has an active quest from that NPC. If not, the quest is added to the player's internal state array (`Client.quests`).

### Quest Types and Validation
To provide varied gameplay, the engine supports multiple quest types:
1. **`defeat` (Combat-based):** 
   When an enemy is defeated in combat, the server checks the player's active quests. If a `defeat` quest requires that specific enemy type, it is marked as completed.
2. **`multi_stage` (Interaction-based):**
   These quests require the player to navigate the world and interact with specific targets (e.g., a "Mysterious Chest"). The server validates progress dynamically when the player interacts with the target using the `TALK` command.

### Progression and Rewards
Players can monitor their active and completed objectives at any time using the `QUESTS` command, which returns a structured JSON list for the GUI to render.
Upon completing a quest's objective, the server triggers the reward logic. Rewards are flexible and can include:
- **Narrative Progression:** Unlocking the next stage of a questline or revealing new dialogue options.
- **Stat Restoration:** Automatically restoring the player's HP to maximum after a difficult objective.
