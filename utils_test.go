package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/FuzzyStatic/blizzard/v1/d3gd"
)

func TestGetPlayersFromData(t *testing.T) {
	testData, err := ioutil.ReadFile("__tests__/data/season_23_rift_wizard.json")
	if err != nil {
		log.Fatalf("Unable to read test data.\n%s", err)
	}

	// unmarshaling is failing. Leaving test JSON out of commit.
	var leaderboardData d3gd.Leaderboard
	err = json.Unmarshal(testData, &leaderboardData)
	if err != nil {
		t.Errorf("Unable to setup leaderboard data.\n%s", err)
	}

	t.Run("Can load all players", func(t *testing.T) {
		players := getPlayersFromData(&leaderboardData)
		expectedPlayerCount := 1000
		if len(players) != expectedPlayerCount {
			t.Errorf("Players != %d", expectedPlayerCount)
			t.FailNow()
		}

	})

}
