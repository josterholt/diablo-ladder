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

type LaderContext struct {
	Players []Player
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file.\n%s", err)
	}
}

func getPlayers() []Player {
	blizz := blizzard.NewClient(os.Getenv("CLIENT_ID"), os.Getenv("SECRET_KEY"), blizzard.US, blizzard.EnUS)

	err := blizz.AccessTokenRequest()
	if err != nil {
		log.Fatalf("Could not connect to Blizzard API.\n%s", err)
	}

	current_season, err := strconv.Atoi(os.Getenv("CURRENT_SEASON"))
	if err != nil {
		log.Fatalf("Error parsing season.\n%s", err)
	}

	data, _, err := blizz.D3SeasonLeaderboardWizard(current_season)
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
	var ladder_player Player

	// data.Row[0].Player[0].Data[0].ID
	for h := 0; h < len(data.Row); h++ {
		row := data.Row[h]

		// Gather Player information
		for i := 0; i < len(row.Player); i++ {
			player := row.Player[i]

			for j := 0; j < len(player.Data); j++ {
				player_data := player.Data[j]

				switch player_data.ID {
				case "HeroBattleTag":
					ladder_player.BattleTag = strings.Split(player_data.String, "#")[0]
				case "HeroClass":
					ladder_player.HeroClass = strings.Title(strings.ToLower(player_data.String))
				case "HeroLevel":
					ladder_player.HeroLevel = player_data.Number
				case "ParagonLevel":
					ladder_player.ParagonLevel = player_data.Number
				}
			}
		}

		// Gather Rift information
		for i := 0; i < len(row.Data); i++ {
			rift_data_point := row.Data[i]

			switch rift_data_point.ID {
			case "Rank":
				ladder_player.Rank = rift_data_point.Number
			case "RiftLevel":
				ladder_player.RiftLevel = rift_data_point.Number
			case "RiftTime":
				seconds_to_nanoseconds := int64(1000000)
				ladder_player.RiftTime = time.Duration(rift_data_point.Timestamp * seconds_to_nanoseconds)
			case "CompletedTime":
				// Comleted time is in microseconds
				ladder_player.CompletedTime = time.Unix(rift_data_point.Timestamp/1000, 0)
			}
		}
		players = append(players, ladder_player)
	}
	return players
}
