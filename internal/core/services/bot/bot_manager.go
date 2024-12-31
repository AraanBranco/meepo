package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/AraanBranco/meepow/internal/config"
	"github.com/AraanBranco/meepow/internal/core/interfaces"
	"github.com/paralin/go-dota2"
	"github.com/paralin/go-dota2/cso"
	"github.com/paralin/go-dota2/protocol"
	"github.com/paralin/go-dota2/state"
	"github.com/paralin/go-steam"
	"github.com/paralin/go-steam/protocol/steamlang"
	"github.com/paralin/go-steam/steamid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type BotManager struct {
	Redis       *redis.Client
	SteamClient *steam.Client
	DotaClient  *dota2.Dota2
	Config      config.Config
	Logger      *zap.Logger
}

func New(conf config.Config, rs *redis.Client, steam *steam.Client) *BotManager {
	return &BotManager{
		Redis:       rs,
		SteamClient: steam,
		Config:      conf,
		Logger:      zap.L().With(zap.String("service", "bot")),
	}
}

func (b *BotManager) GetLobbyData(referenceID string) (interfaces.BotConfig, error) {
	lobbyData, err := b.Redis.Get(context.Background(), fmt.Sprintf("lobby:%s", referenceID)).Result()
	if err != nil {
		if err == redis.Nil {
			return interfaces.BotConfig{
				ReferenceID: "",
			}, nil
		}
		b.Logger.Error("Error getting lobby status from Redis", zap.Error(err))
		return interfaces.BotConfig{}, err
	}

	var botConfig interfaces.BotConfig
	err = json.Unmarshal([]byte(lobbyData), &botConfig)
	if err != nil {
		b.Logger.Error("Error unmarshalling lobby data", zap.Error(err))
		return interfaces.BotConfig{}, err
	}

	return botConfig, nil
}

func (b *BotManager) StartupBot(lobbyData interfaces.BotConfig) {
	b.Logger.Info("Bot started")

	// Get Bots for logging
	logOnDetails := &steam.LogOnDetails{
		Username:               b.Config.GetString("bot.username"),
		Password:               b.Config.GetString("bot.password"),
		ShouldRememberPassword: true,
	}
	b.updateStatusRedis(lobbyData.ReferenceID, interfaces.LOBBY_CREATED, 0)

	b.Logger.Info("Lobby Data", zap.Any("lobbyData", lobbyData))

	// Log on to Steam
	b.SteamClient.Connect()
	for event := range b.SteamClient.Events() {
		switch event.(type) {
		case *steam.ConnectedEvent:
			b.SteamClient.Auth.LogOn(logOnDetails)

		case *steam.LoggedOnEvent:
			b.SteamClient.Social.SetPersonaState(steamlang.EPersonaState_LookingToPlay)
			logger := logrus.New()
			logger.Level = logrus.DebugLevel
			b.DotaClient = dota2.New(b.SteamClient, logger)
			b.DotaClient.SetPlaying(true)

			time.Sleep(5 * time.Second)
			b.DotaClient.SayHello()

			lobbyGameMode, err := strconv.Atoi(lobbyData.GameMode)
			if err != nil {
				b.Logger.Error("Error converting game mode to int", zap.Error(err))
				return
			}

			lobbyRegionSouthAmerica, err := strconv.Atoi(lobbyData.Region)
			if err != nil {
				b.Logger.Error("Error converting region to int", zap.Error(err))
				return
			}
			dota2s := state.Dota2State{
				Lobby: &protocol.CSODOTALobby{
					ServerRegion: proto.Uint32(uint32(lobbyRegionSouthAmerica)),
					GameMode:     proto.Uint32(uint32(lobbyGameMode)),
				},
			}

			var leagueDotaID int
			if lobbyData.LeagueID == "" {
				lobbyData.LeagueID = "0"
			}

			leagueDotaID, err = strconv.Atoi(lobbyData.LeagueID)
			if err != nil {
				b.Logger.Error("Error converting LeagueID to int", zap.Error(err))
				return
			}
			if leagueDotaID != 0 {
				dota2s.Lobby.LeaderId = proto.Uint64(uint64(leagueDotaID))
			}

			var visibility protocol.DOTALobbyVisibility
			if lobbyData.Visibility == "0" {
				visibility = protocol.DOTALobbyVisibility_DOTALobbyVisibility_Public
			} else {
				visibility = protocol.DOTALobbyVisibility_DOTALobbyVisibility_Unlisted
			}

			time.Sleep(15 * time.Second)
			if dota2s.IsReady() {
				lobbyDetails := &protocol.CMsgPracticeLobbySetDetails{
					AllowCheats:      proto.Bool(b.Config.GetBool("bot.allowCheats")),
					GameName:         proto.String(lobbyData.LobbyName),
					FillWithBots:     proto.Bool(false),
					PassKey:          proto.String(lobbyData.LobbyPassword),
					CustomMaxPlayers: proto.Uint32(10),
					AllowSpectating:  proto.Bool(true),
					Visibility:       visibility.Enum(),
					ServerRegion:     proto.Uint32(uint32(lobbyRegionSouthAmerica)),
					GameMode:         proto.Uint32(uint32(lobbyGameMode)),
					Allchat:          proto.Bool(true),
				}

				// Set League ID
				if leagueDotaID != 0 {
					lobbyDetails.Leagueid = proto.Uint32(uint32(leagueDotaID))
				}

				// Leave bot from lobby if it's already in one
				b.DotaClient.LeaveCreateLobby(context.TODO(), lobbyDetails, true)

				// Remove Bot from slot team
				b.DotaClient.KickLobbyMemberFromTeam(b.SteamClient.SteamId().GetAccountId())

				// Invite all members
				for _, player := range lobbyData.Players {
					steamID, _ := steamid.NewId(player.SteamID)
					b.DotaClient.InviteLobbyMember(steamID)
				}

				_ = b.updateStatusRedis(lobbyData.ReferenceID, interfaces.LOBBY_STARTED, 0)
			}

			eventCh, _, err := b.DotaClient.GetCache().SubscribeType(cso.Lobby)
			if err != nil {
				fmt.Print("Error in event: ", err)
				b.SteamClient.Disconnect()
				return
			}

			go func() {
				ticker := time.NewTicker(3 * time.Second)
				for range ticker.C {
					counterPlayers := 0
					b.DotaClient.SayHello()
					lobbyEvent := <-eventCh
					lob := lobbyEvent.Object.(*protocol.CSODOTALobby)

					if dota2s.IsReady() {
						if lob.GetState() == protocol.CSODOTALobby_UI {
							cacheCtr, _ := b.DotaClient.GetCache().GetContainerForTypeID(uint32(cso.Lobby))
							lobbyObjCache := cacheCtr.GetOne()
							lobCache := lobbyObjCache.(*protocol.CSODOTALobby)

							for _, member := range lobCache.GetAllMembers() {
								if member.GetName() != b.Config.GetString("bot.username") {
									counterPlayers++
								}

								if counterPlayers >= b.Config.GetInt("lobby.maxPlayers") {
									b.DotaClient.LaunchLobby()
								}
							}
						}

						if lob.GetState() == protocol.CSODOTALobby_RUN {
							lobbyID := fmt.Sprintf("%v", lob.GetLobbyId())
							matchID := fmt.Sprintf("%v", lob.GetMatchId())

							lobbyData.MatchID = matchID
							lobbyData.LobbyID = lobbyID
							_ = b.updateLobbyRedis(lobbyData.ReferenceID, lobbyData, 0)
							_ = b.updateStatusRedis(lobbyData.ReferenceID, interfaces.LOBBY_RUNNING, 0)
						}

						if lob.GetState() == protocol.CSODOTALobby_POSTGAME {
							fmt.Println("Lobby is finished")
							var teamVictory string
							if lob.GetMatchOutcome() == protocol.EMatchOutcome_k_EMatchOutcome_DireVictory {
								teamVictory = "DIRE"
							} else {
								teamVictory = "RADIANT"
							}

							lobbyData.Duration = int(lob.GetMatchDuration())
							lobbyData.Victory = teamVictory
							// Add expires for expire lobby data from Redis
							_ = b.updateLobbyRedis(lobbyData.ReferenceID, lobbyData, 3600)                  // 1 hour
							_ = b.updateStatusRedis(lobbyData.ReferenceID, interfaces.LOBBY_FINISHED, 3600) // 1 hour
							ticker.Stop()
							b.SteamClient.Disconnect()
						}
					}
				}
			}()

		case *steam.DisconnectedEvent:
			fmt.Printf("Bot %s disconnected from Steam", logOnDetails.Username)
			os.Exit(1)
			return
		}
	}
}

func (b *BotManager) updateLobbyRedis(referenceID string, lobbyData interfaces.BotConfig, expires int) error {
	data, err := json.Marshal(lobbyData)
	if err != nil {
		b.Logger.Error("Error marshalling lobby data", zap.Error(err))
		return err
	}

	err = b.Redis.Set(context.Background(), fmt.Sprintf("lobby:%s", referenceID), string(data), time.Duration(expires)).Err()
	if err != nil {
		b.Logger.Error("Error updating lobby status in Redis", zap.Error(err))
		return err
	}

	return nil
}

func (b *BotManager) updateStatusRedis(referenceID string, status string, expires int) error {
	err := b.Redis.Set(context.Background(), fmt.Sprintf("lobby:%s:status", referenceID), status, time.Duration(expires)).Err()
	if err != nil {
		b.Logger.Error("Error updating lobby status in Redis", zap.Error(err))
		return err
	}

	return nil
}
