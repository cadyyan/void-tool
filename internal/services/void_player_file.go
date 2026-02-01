package services

import (
	"context"
	"fmt"
	"io/fs"
	"time"

	"github.com/BurntSushi/toml"
)

type VoidPlayerFileService struct {
	fs fs.FS
}

func NewVoidPlayerFileService(fileSystem fs.FS) *VoidPlayerFileService {
	return &VoidPlayerFileService{
		fs: fileSystem,
	}
}

var _ VoidPlayerService = (*VoidPlayerFileService)(nil)

func (service *VoidPlayerFileService) GetAllPlayers(
	ctx context.Context,
) ([]Player, error) {
	playerFiles, err := fs.Glob(service.fs, "*.toml")
	if err != nil {
		return nil, fmt.Errorf("unable to find player files: %w", err)
	}

	players := make([]Player, len(playerFiles))
	for index, playerFile := range playerFiles {
		save, err := service.parsePlayerFile(service.fs, playerFile)
		if err != nil {
			return nil, fmt.Errorf("unable to get some player data: %w", err)
		}

		players[index] = save
	}

	return players, nil
}

func (service *VoidPlayerFileService) parsePlayerFile(
	fileSystem fs.FS,
	filePath string,
) (Player, error) {
	var save PlayerSaveFileFormat

	_, err := toml.DecodeFS(fileSystem, filePath, &save)
	if err != nil {
		return Player{}, fmt.Errorf("unable to read player save file: %w", err)
	}

	experience := make(map[string]float64)
	levels := make(map[string]int)

	for index, skill := range skillOrder {
		experience[skill] = float64(save.Experience[index]) / 10.0
		levels[skill] = save.Levels[index]
	}

	creationTime := time.UnixMilli(save.Variables.Creation)

	return Player{
		AccountName: save.AccountName,
		Experience:  experience,
		Levels:      levels,
		CreatedOn:   creationTime,
	}, nil
}

type PlayerSaveFileFormat struct {
	AccountName string                        `toml:"accountName"`
	Experience  []int                         `toml:"experience"`
	Levels      []int                         `toml:"levels"`
	Variables   PlayerSaveFileVariablesFormat `toml:"variables"`
}

type PlayerSaveFileVariablesFormat struct {
	Creation int64 `toml:"creation"`
}

var skillOrder = []string{
	"Attack",
	"Defence",
	"Strength",
	"Constitution",
	"Ranged",
	"Prayer",
	"Magic",
	"Cooking",
	"Woodcutting",
	"Fletching",
	"Fishing",
	"Firemaking",
	"Crafting",
	"Smithing",
	"Mining",
	"Herblore",
	"Agility",
	"Thieving",
	"Slayer",
	"Farming",
	"Runecrafting",
	"Hunter",
	"Construction",
	"Summoning",
	"Dungeoneering",
}
