# Diablo III Ladder Dashboard

This is a sample Go app that displays Diablo III ladder information.

Technology:

-   [Go](https://golang.org/)
-   [Goji](https://github.com/goji/goji) (Web Framework)
-   [Air](https://github.com/cosmtrek/air) (Live Reload)
-   [Go client library for Blizzard API](https://github.com/FuzzyStatic/blizzard)
-   [Blizzard API](https://develop.battle.net)

# SETUP

Uses the following environment variables to make Blizzard API queries. (Pulls in values with [godotenv](https://github.com/joho/godotenv))

-   CLIENT_ID - Blizzard API Client ID
-   SECRET_KEY - Blizzard API secret key
-   CURRENT_SEASON - Number for season to display

# Development

Run with `air` for live reloading.
