package world

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)


type Room struct {
	ID			string				`yaml:"id"`
	Name		string				`yaml:"name"`
	Description	string				`yaml:"description"`
	Exit		map[string]string	`yaml:"exits"`
}

type World struct {
	Rooms []Room `yaml:"rooms"`
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

	return &w, nil
}
