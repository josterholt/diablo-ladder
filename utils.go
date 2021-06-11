package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/FuzzyStatic/blizzard/v1"
	"github.com/FuzzyStatic/blizzard/v1/d3gd"
	"github.com/joho/godotenv"
)

type Player struct {
	BattleTag     string
	HeroClass     string
	HeroLevel     int
	ParagonLevel  int
	Rank          int
	RiftLevel     int
	RiftTime      time.Duration
	CompletedTime time.Time
}

type LadderContext struct {
	Players []Player
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file.\n%s", err)
	}
}

func getPlayers(playerClass string, hardcore bool) []Player {
	blizz := blizzard.NewClient(os.Getenv("CLIENT_ID"), os.Getenv("SECRET_KEY"), blizzard.US, blizzard.EnUS)

	err := blizz.AccessTokenRequest()
	if err != nil {
		log.Fatalf("Could not connect to Blizzard API.\n%s", err)
	}

	currentSeason, err := strconv.Atoi(os.Getenv("CURRENT_SEASON"))
	if err != nil {
		log.Fatalf("Error parsing season.\n%s", err)
	}

	//classes := []string {"Barbarian", "Crusader", "DemonHunter", "Monk", "Necromancer", "WitchDoctor", "Wizard"}
	// playerClass := "BARBARIAN"
	// hardcore := false
	var data *d3gd.Leaderboard
	var _ []byte

	switch playerClass {
	case "BARBARIAN":
		if !hardcore {
			data, _, err = blizz.D3SeasonLeaderboardBarbarian(currentSeason)
		} else {
			data, _, err = blizz.D3SeasonLeaderboardHardcoreBarbarian(currentSeason)
		}
	case "CRUSADER":
		if !hardcore {
			data, _, err = blizz.D3SeasonLeaderboardCrusader(currentSeason)
		} else {
			data, _, err = blizz.D3SeasonLeaderboardHardcoreCrusader(currentSeason)
		}
	case "DEMONHUNTER":
		if !hardcore {
			data, _, err = blizz.D3SeasonLeaderboardDemonHunter(currentSeason)
		} else {
			data, _, err = blizz.D3SeasonLeaderboardHardcoreDemonHunter(currentSeason)
		}
	case "MONK":
		if !hardcore {
			data, _, err = blizz.D3SeasonLeaderboardMonk(currentSeason)
		} else {
			data, _, err = blizz.D3SeasonLeaderboardHardcoreMonk(currentSeason)
		}
	case "NECROMANCER":
		if !hardcore {
			data, _, err = blizz.D3SeasonLeaderboardNecromancer(currentSeason)
		} else {
			data, _, err = blizz.D3SeasonLeaderboardHardcoreNecromancer(currentSeason)
		}
	case "WITCHDOCTOR":
		if !hardcore {
			blizz.D3SeasonLeaderboardWitchDoctor(currentSeason)
		} else {
			blizz.D3SeasonLeaderboardHardcoreWitchDoctor(currentSeason)
		}
	case "WIZARD":
		if !hardcore {
			data, _, err = blizz.D3SeasonLeaderboardWizard(currentSeason)
		} else {
			data, _, err = blizz.D3SeasonLeaderboardHardcoreWizard(currentSeason)
		}
	}

	//data, _, err := blizz.D3SeasonLeaderboardWizard(currentSeason)
	if err != nil {
		log.Fatalf("Could not retrieve season data.\n%s", err)
	}

	return getPlayersFromData(data)
}

func FormatDuration(riftTime time.Duration) string {
	secondsInMinute := 60
	riftMinutes := int(riftTime.Minutes())
	riftSeconds := int(riftTime.Seconds()) - (riftMinutes * secondsInMinute)
	return fmt.Sprintf("%dm %ds", riftMinutes, riftSeconds)
}

func getPlayersFromData(data *d3gd.Leaderboard) []Player {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("stacktrace from panic: \n%s", string(debug.Stack()))
		}
	}()

	var players []Player
	var ladderPlayer Player

	// data.Row[0].Player[0].Data[0].ID
	for h := 0; h < len(data.Row); h++ {
		row := data.Row[h]

		// Gather Player information
		for i := 0; i < len(row.Player); i++ {
			player := row.Player[i]

			for j := 0; j < len(player.Data); j++ {
				playerData := player.Data[j]

				switch playerData.ID {
				case "HeroBattleTag":
					ladderPlayer.BattleTag = strings.Split(playerData.String, "#")[0]
				case "HeroClass":
					ladderPlayer.HeroClass = strings.Title(strings.ToLower(playerData.String))
				case "HeroLevel":
					ladderPlayer.HeroLevel = playerData.Number
				case "ParagonLevel":
					ladderPlayer.ParagonLevel = playerData.Number
				}
			}
		}

		// Gather Rift information
		for i := 0; i < len(row.Data); i++ {
			riftDataPoint := row.Data[i]

			switch riftDataPoint.ID {
			case "Rank":
				ladderPlayer.Rank = riftDataPoint.Number
			case "RiftLevel":
				ladderPlayer.RiftLevel = riftDataPoint.Number
			case "RiftTime":
				secondsToNanoseconds := int64(1000000)
				ladderPlayer.RiftTime = time.Duration(riftDataPoint.Timestamp * secondsToNanoseconds)
			case "CompletedTime":
				// Comleted time is in microseconds
				ladderPlayer.CompletedTime = time.Unix(riftDataPoint.Timestamp/1000, 0)
			}
		}
		players = append(players, ladderPlayer)
	}
	return players
}
