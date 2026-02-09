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
) ([]VoidPlayer, error) {
	playerFiles, err := fs.Glob(service.fs, "*.toml")
	if err != nil {
		return nil, fmt.Errorf("unable to find player files: %w", err)
	}

	players := make([]VoidPlayer, len(playerFiles))
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
) (VoidPlayer, error) {
	var save PlayerSaveFileFormat

	_, err := toml.DecodeFS(fileSystem, filePath, &save)
	if err != nil {
		return VoidPlayer{}, fmt.Errorf("unable to read player save file: %w", err)
	}

	experience := make(map[string]float64)
	levels := make(map[string]int)

	for index, skill := range skillOrder {
		experience[skill] = float64(save.Experience[index]) / 10.0

		if skill == "Constitution" {
			levels[skill] = calculateLevelFromExperience(experience[skill])
		} else {
			levels[skill] = save.Levels[index]
		}
	}

	creationTime := time.UnixMilli(save.Variables.Creation)

	return VoidPlayer{
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

func calculateLevelFromExperience(exp float64) int {
	for index, requiredExp := range experienceTable {
		if exp < requiredExp {
			return index + 1
		}
	}

	return 99
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

var (
	experienceTable = []float64{
		0,
		83,
		174,
		276,
		388,
		512,
		650,
		801,
		969,
		1154,
		1358,
		1584,
		1833,
		2107,
		2411,
		2746,
		3115,
		3523,
		4470,
		5018,
		5624,
		6291,
		7028,
		7842,
		8740,
		9730,
		10824,
		12031,
		13363,
		14833,
		16456,
		18247,
		20224,
		22406,
		24815,
		27473,
		30408,
		33648,
		37224,
		41171,
		45529,
		50339,
		55649,
		61512,
		67983,
		75127,
		83014,
		91721,
		101333,
		111945,
		123660,
		136594,
		150872,
		166636,
		184040,
		203254,
		224466,
		247886,
		273742,
		302288,
		333804,
		368599,
		407015,
		449428,
		496254,
		547953,
		605032,
		668051,
		737627,
		814445,
		899257,
		992895,
		1096278,
		1210421,
		1336443,
		1475581,
		1629200,
		1798808,
		1986068,
		2192818,
		2424087,
		2673114,
		2951373,
		3258594,
		3597792,
		3972294,
		4385776,
		4842295,
		5346332,
		5902831,
		6517253,
		7195629,
		7944614,
		8771558,
		9684577,
		10692629,
		11805606,
		13034431,
	}
)
