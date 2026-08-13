package world

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Spawn struct {
	NPCType string `yaml:"npc_type"`
	Count   int    `yaml:"count"`
}

type Room struct {
	ID          string            `yaml:"-"`
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Exits       map[string]string `yaml:"exits"`
	Spawns      []Spawn           `yaml:"spawns"`
	Items       []string          `yaml:"items"`
}

type Item struct {
	ID          string `yaml:"-"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Obtainable  bool   `yaml:"obtainable"`
}

type Stats struct {
	HP     int `yaml:"hp"`
	Damage int `yaml:"damage"`
}

type NPC struct {
	ID          string   `yaml:"-"`
	Name        string   `yaml:"name"`
	Role        string   `yaml:"role"`
	Description string   `yaml:"description"`
	Dialogue    []string `yaml:"dialogue"`
	Stats       Stats    `yaml:"stats"`
	Quests      []string `yaml:"quests"`
}

type Quest struct {
	ID          string `yaml:"-"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Target      string `yaml:"target"`
	Reward      string `yaml:"reward"`
}

type WorldData struct {
	Locations map[string]*Room  `yaml:"locations"`
	Items     map[string]*Item  `yaml:"items"`
	NPCs      map[string]*NPC   `yaml:"npcs"`
	Quests    map[string]*Quest `yaml:"quests"`
}

type yamlRoot struct {
	World WorldData `yaml:"world"`
}

type World struct {
	Rooms  map[string]*Room
	Items  map[string]*Item
	NPCs   map[string]*NPC
	Quests map[string]*Quest

	StartRoomID string
}

func (w *World) RoomByID(id string) (*Room, bool) {
	r, ok := w.Rooms[id]
	return r, ok
}

func LoadWorld(filename string) (*World, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read world file: %w", err)
	}

	var root yamlRoot
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	w := &World{
		Rooms:  root.World.Locations,
		Items:  root.World.Items,
		NPCs:   root.World.NPCs,
		Quests: root.World.Quests,
	}

	if w.Rooms == nil {
		w.Rooms = make(map[string]*Room)
	}
	if w.Items == nil {
		w.Items = make(map[string]*Item)
	}
	if w.NPCs == nil {
		w.NPCs = make(map[string]*NPC)
	}
	if w.Quests == nil {
		w.Quests = make(map[string]*Quest)
	}

	if err := w.index(); err != nil {
		return nil, err
	}

	if err := w.validate(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *World) index() error {
	// Set IDs
	for id, r := range w.Rooms {
		r.ID = id
		if w.StartRoomID == "" {
			w.StartRoomID = id // just pick the first one as start for now if needed, though server might need a specific one
		}
	}
	// Let's enforce destiny_islands as start room if it exists
	if _, ok := w.Rooms["destiny_islands"]; ok {
		w.StartRoomID = "destiny_islands"
	} else if _, ok := w.Rooms["start"]; ok {
		w.StartRoomID = "start"
	}

	for id, i := range w.Items {
		i.ID = id
	}
	for id, n := range w.NPCs {
		n.ID = id
	}
	for id, q := range w.Quests {
		q.ID = id
	}
	return nil
}

// Rejects any exit pointing to a room ID that doesn't exist, items that don't exist, etc.
func (w *World) validate() error {
	for roomID, r := range w.Rooms {
		for direction, target := range r.Exits {
			if _, ok := w.Rooms[target]; !ok {
				return fmt.Errorf("world data invalid: room %q exit %q points to unknown room %q", roomID, direction, target)
			}
		}
		for _, itemID := range r.Items {
			if _, ok := w.Items[itemID]; !ok {
				return fmt.Errorf("world data invalid: room %q contains unknown item %q", roomID, itemID)
			}
		}
		for _, spawn := range r.Spawns {
			if _, ok := w.NPCs[spawn.NPCType]; !ok {
				return fmt.Errorf("world data invalid: room %q spawns unknown npc %q", roomID, spawn.NPCType)
			}
		}
	}
	return nil
}
