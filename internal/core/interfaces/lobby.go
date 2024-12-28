package interfaces

import (
	"fmt"

	"github.com/AraanBranco/meepo/internal/config"
)

type PostLobbyRequest struct {
	ReferenceID     string         `json:"reference_id"`
	LeagueID        string         `json:"league_id"`
	GameMode        string         `json:"game_mode"`
	Region          string         `json:"region"`
	LobbyName       string         `json:"lobby_name"`
	LobbyPassword   string         `json:"lobby_password"`
	AllowSpectating bool           `json:"allow_spectating"`
	Visibility      string         `json:"visibility"`
	Players         []PlayersMatch `json:"players"`
}

type PlayersMatch struct {
	SteamID string `json:"steam_id"`
	Team    string `json:"team"`
}

func (pr *PostLobbyRequest) Validate(configs config.Config) (bool, []string) {
	var validationErrors []string

	if pr.ReferenceID == "" {
		validationErrors = append(validationErrors, "reference_id is required")
	}

	if pr.GameMode == "" {
		validationErrors = append(validationErrors, "game_mode is required")
	}

	if pr.GameMode != "" {
		switch pr.GameMode {
		// 1 - All Pick
		// 2 - Captains Mode
		case "1", "2":
			// is valid

		default:
			validationErrors = append(validationErrors, fmt.Sprintf("game_mode %s is invalid", pr.GameMode))
		}

	}

	if pr.Region == "" {
		validationErrors = append(validationErrors, "region is required")
	}

	if pr.Region != "" {
		switch pr.Region {
		// 1 - US West
		// 2 - US East
		// 3 - Europe
		// 8 - Stockholm
		// 10 - Brazil
		// 15 - Peru
		case "1", "2", "3", "8", "10", "15":
			// is valid

		default:
			validationErrors = append(validationErrors, fmt.Sprintf("region %s is invalid", pr.Region))
		}
	}

	if pr.LobbyName == "" {
		validationErrors = append(validationErrors, "lobby_name is required")
	}

	if pr.LobbyPassword == "" {
		validationErrors = append(validationErrors, "lobby_password is required")
	}

	maxPlayers := configs.GetInt("lobby.maxPlayers")
	if len(pr.Players) < maxPlayers {
		validationErrors = append(validationErrors, fmt.Sprintf("must have %d players", maxPlayers))
	}

	if len(pr.Players) == maxPlayers {
		for _, player := range pr.Players {
			if player.SteamID == "" {
				validationErrors = append(validationErrors, "have player without steam_id")
			}

			if player.Team == "" {
				validationErrors = append(validationErrors, "have player without team")
			}

			if player.Team != "" {
				switch player.Team {
				case "dire", "radiant":
					// is valid

				default:
					validationErrors = append(validationErrors, fmt.Sprintf("team %s is invalid", player.Team))
				}
			}

			if player.SteamID == "" && player.Team == "" {
				return true, validationErrors
			}
		}
	}

	if pr.Visibility == "" {
		validationErrors = append(validationErrors, "visibility is required")
	}

	if pr.Visibility != "" {
		switch pr.Visibility {
		// 0 - Public
		// 1 - Friends - Do not use it, as it will be visible to the bot's friends
		// 2 - Private
		case "0", "2":
			// is valid

		default:
			validationErrors = append(validationErrors, fmt.Sprintf("visibility %s is invalid", pr.Visibility))
		}
	}

	if len(validationErrors) > 0 {
		return true, validationErrors
	}

	return false, validationErrors
}

type PostLobbyResponse struct {
	ReferenceID string `json:"reference_id"`
	LobbyStatus string `json:"lobby_status"`
}

type GetLobbyResponse struct {
	ReferenceID     string      `json:"reference_id"`
	LeagueID        string      `json:"league_id"`
	MatchID         string      `json:"match_id"`
	LobbyID         string      `json:"lobby_id"`
	GameMode        string      `json:"game_mode"`
	Region          string      `json:"region"`
	LobbyName       string      `json:"lobby_name"`
	LobbyPassword   string      `json:"lobby_password"`
	AllowSpectating bool        `json:"allow_spectating"`
	Visibility      string      `json:"visibility"`
	Victory         string      `json:"victory"`
	Players         []BotPlayer `json:"players"`
	Duration        int         `json:"duraration"`
	LobbyStatus     string      `json:"lobby_status"`
}
