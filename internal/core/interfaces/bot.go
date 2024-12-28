package interfaces

type BotConfig struct {
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
}

type BotPlayer struct {
	SteamID string `json:"steam_id"`
	Team    string `json:"team"`
}
