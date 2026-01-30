package services

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/BurntSushi/toml"
)

type UserFileService struct {
	fs fs.FS
}

func NewUserFileService(fileSystem fs.FS) *UserFileService {
	return &UserFileService{
		fs: fileSystem,
	}
}

var _ UserService = (*UserFileService)(nil)

func (service *UserFileService) GetAllUsers(
	ctx context.Context,
) ([]Player, error) {
	userFiles, err := fs.Glob(service.fs, "*.toml")
	if err != nil {
		return nil, fmt.Errorf("unable to find user files: %w", err)
	}

	users := make([]Player, len(userFiles))
	for index, userFile := range userFiles {
		save, err := service.parseUserFile(service.fs, userFile)
		if err != nil {
			return nil, fmt.Errorf("unable to get some user data: %w", err)
		}

		users[index] = save
	}

	return users, nil
}

func (service *UserFileService) parseUserFile(
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

	return Player{
		AccountName: save.AccountName,
		Experience:  experience,
		Levels:      levels,
	}, nil
}

type PlayerSaveFileFormat struct {
	AccountName string `toml:"accountName"`
	Experience  []int  `toml:"experience"`
	Levels      []int  `toml:"levels"`
}

var (
	skillOrder = []string{
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
)
