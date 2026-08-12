package world

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Room struct {
	ID          string            `yaml:"id"`
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Exits       map[string]string `yaml:"exits"`
}

type World struct {
	Rooms []Room `yaml:"rooms"`

	byID map[string]*Room
}

func (w *World) RoomByID(id string) (*Room, bool) {
	r, ok := w.byID[id]
	return r, ok
}

func LoadWorld(filename string) (*World, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read world file: %w", err)
	}

	var w World

	err = yaml.Unmarshal(data, &w)
	if err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	if err := w.index(); err != nil {
		return nil, err
	}

	if err := w.validate(); err != nil {
		return nil, err
	}

	return &w, nil
}

func (w *World) index() error {
	w.byID = make(map[string]*Room, len(w.Rooms))
	for i := range w.Rooms {
		r := &w.Rooms[i]
		if r.ID == "" {
			return fmt.Errorf("world data invalid: room %d has no id", i)
		}
		if _, exists := w.byID[r.ID]; exists {
			return fmt.Errorf("world data invalid: duplicate room id %q", r.ID)
		}
		w.byID[r.ID] = r
	}
	return nil
}

// Rejects any exit pointing to a room ID that doesn't exist.
func (w *World) validate() error {
	for _, r := range w.Rooms {
		for direction, target := range r.Exits {
			if _, ok := w.byID[target]; !ok {
				return fmt.Errorf(
					"world data invalid: room %q exit %q points to unknown room %q",
					r.ID, direction, target,
				)
			}
		}
	}
	return nil
}
