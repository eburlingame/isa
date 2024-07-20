package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
	"gopkg.in/olahol/melody.v1"
)

type Command struct {
	ReqId string            `json:"reqId"`
	Verb  string            `json:"v"`
	Data  map[string]string `json:"d"`
}

type Response struct {
	ReqId string            `json:"reqId"`
	Verb  string            `json:"v"`
	Data  map[string]string `json:"d"`
	Error bool              `json:"err"`
}

type GameStatus struct {
	GameId        string      `json:"gameId"`
	GamePneumonic string      `json:"gamePneumonic"`
	IsHost        bool        `json:"isHost"`
	Game          PlayersGame `json:"game"`
	Abandoned     bool        `json:"abandoned"`
}

type GameResponse struct {
	ReqId string     `json:"reqId"`
	Verb  string     `json:"v"`
	Game  GameStatus `json:"d"`
}

type GameUpdate struct {
	Verb string     `json:"v"`
	Game GameStatus `json:"d"`
}

func parseCommand(msg []byte) (Command, error) {
	var command Command

	err := json.Unmarshal(msg, &command)
	if err != nil {
		return command, err
	}

	return command, nil
}

func errorResponse(cmd *Command, message string) Response {
	return Response{
		ReqId: cmd.ReqId,
		Verb:  cmd.Verb,
		Error: true,
		Data: map[string]string{
			"message": message,
		},
	}
}

func unrecognizedCommandResponse(cmd *Command) Response {
	return errorResponse(cmd, "Unrecognized command")
}

func sendResponse(session *melody.Session, response Response) {
	text, _ := json.Marshal(response)
	session.Write(text)
}

func SendGameUpdate(session *melody.Session, gameId string, game *Game, abandoned bool) {
	persistentSession, _ := GetPersistentSession(session)

	playerIndex := GetPlayerIndex(game, persistentSession.PlayerName)

	// Abandoned status
	response := GameUpdate{
		Verb: "gameState",
		Game: GameStatus{
			GameId:        gameId,
			GamePneumonic: "",
			IsHost:        persistentSession.GameHost,
			Abandoned:     true,
		},
	}

	if !(playerIndex == -1 || abandoned) {
		playersGame := GetPlayersGame(game, playerIndex)

		response = GameUpdate{
			Verb: "gameState",
			Game: GameStatus{
				GameId:        gameId,
				GamePneumonic: game.GamePneumonic,
				IsHost:        persistentSession.GameHost,
				Game:          *playersGame,
				Abandoned:     abandoned,
			},
		}
	}

	text, _ := json.Marshal(response)
	session.Write(text)
}

func SendGameResponse(session *melody.Session, cmd *Command, gameId string, game *Game, abandoned bool) {
	persistentSession, _ := GetPersistentSession(session)

	response := GameResponse{
		ReqId: cmd.ReqId,
		Verb:  cmd.Verb,
		Game: GameStatus{
			GameId:        gameId,
			GamePneumonic: "",
			IsHost:        persistentSession.GameHost,
			Abandoned:     true,
		},
	}

	playerIndex := GetPlayerIndex(game, persistentSession.PlayerName)

	if !(playerIndex == -1 || abandoned) {
		playersGame := GetPlayersGame(game, playerIndex)
		response = GameResponse{
			ReqId: cmd.ReqId,
			Verb:  cmd.Verb,
			Game: GameStatus{
				GameId:        gameId,
				GamePneumonic: game.GamePneumonic,
				IsHost:        persistentSession.GameHost,
				Game:          *playersGame,
				Abandoned:     abandoned,
			},
		}
	}

	text, _ := json.Marshal(response)
	session.Write(text)
}

func subscribeToGame(ctx *context.Context, rdb *redis.Client, session *melody.Session, gameId string) chan bool {
	doneChan := make(chan bool)

	log.Printf("Subscribing to gameId: %s", gameId)

	// Kick off a go routine to watch the game state topic
	// When the session is closed, a true value will be send to the doneChan
	// and the routine will return
	go WatchGame(ctx, rdb, session, gameId, doneChan)

	return doneChan
}

func openSession(ctx *context.Context, rdb *redis.Client, session *melody.Session, cmd *Command) error {
	var persistentSession *PersistentSession

	// If an existing sessionId was provided, try to reuse the session from Redis
	if key, ok := cmd.Data["sessionId"]; ok {
		log.Printf("Trying to reopen session: %s", key)
		persistentSession = FetchPersistentSession(ctx, session, rdb, key)
		log.Printf("Reopened persistent session: %s", persistentSession.SessionId)
	} else {
		persistentSession = NewPersistentSession()
	}

	// Store the session data in the melody session store
	SetPersistentSession(ctx, session, rdb, persistentSession)

	// If we are attached to a game, broadcast the game state to the client
	gameId := persistentSession.ActiveGame
	if gameId != "" {
		// Load the game from redis
		game, err := LoadGame(ctx, rdb, gameId)

		// If the game doesn't exist, detach it from the session
		if err != nil || game.State == GameAbandoned {
			persistentSession.ActiveGame = ""
			SetPersistentSession(ctx, session, rdb, persistentSession)
		} else {
			// Re-subscribe to the active game
			persistentSession.UnsubChan = subscribeToGame(ctx, rdb, session, gameId)
			SetPersistentSession(ctx, session, rdb, persistentSession)

			SendGameUpdate(session, gameId, game, false)
		}
	}

	sendResponse(session, Response{
		ReqId: cmd.ReqId,
		Verb:  "openSession",
		Data: map[string]string{
			"sessionId": persistentSession.SessionId,
		},
	})

	return nil
}

func createGame(ctx *context.Context, rdb *redis.Client, session *melody.Session, cmd *Command) error {
	persistentSession, err := GetPersistentSession(session)
	if err != nil {
		return errors.New("Something went wrong loading your session")
	}

	gameId := NewGameCode()
	gamePneumonic := MakeGamePneumonic(gameId)
	fmt.Println(gamePneumonic)

	if playerName, ok := cmd.Data["playerName"]; ok {
		game := EmptyGame(gameId, gamePneumonic)
		game = AddPlayer(game, playerName)

		SaveGame(ctx, rdb, gameId, game)

		persistentSession.GameHost = true
		persistentSession.PlayerName = playerName
		persistentSession.ActiveGame = gameId
		persistentSession.UnsubChan = subscribeToGame(ctx, rdb, session, gameId)

		SetPersistentSession(ctx, session, rdb, persistentSession)
		SendGameResponse(session, cmd, gameId, game, false)

		return nil
	}

	return errors.New("Expected gameId and playerName to be supplied")
}

func joinGame(ctx *context.Context, rdb *redis.Client, session *melody.Session, cmd *Command) error {
	persistentSession, err := GetPersistentSession(session)
	if err != nil {
		return errors.New("Something went wrong loading your session")
	}

	gameId, ok1 := cmd.Data["gameId"]
	playerName, ok2 := cmd.Data["playerName"]

	if ok1 && ok2 {
		if GameExists(ctx, rdb, gameId) {
			game, err := LoadGame(ctx, rdb, gameId)
			if err != nil {
				return errors.New("Error fetching game")
			}

			for _, player := range game.Players {
				if player.Name == playerName {
					return errors.New("Player already exists")
				}
			}

			game = AddPlayer(game, playerName)
			SaveGame(ctx, rdb, gameId, game)

			persistentSession.GameHost = false
			persistentSession.PlayerName = playerName
			persistentSession.ActiveGame = gameId
			persistentSession.UnsubChan = subscribeToGame(ctx, rdb, session, gameId)

			SetPersistentSession(ctx, session, rdb, persistentSession)
			SendGameResponse(session, cmd, gameId, game, false)

			return nil
		} else {
			return errors.New("Game not found")
		}
	}

	return errors.New("Expected gameId and playerName to be supplied")
}

func leaveGame(ctx *context.Context, rdb *redis.Client, session *melody.Session, cmd *Command) error {
	persistentSession, err := GetPersistentSession(session)
	if err != nil {
		return errors.New("Something went wrong loading your session")
	}

	gameId := persistentSession.ActiveGame
	playerName := persistentSession.PlayerName

	if GameExists(ctx, rdb, gameId) {
		game, err := LoadGame(ctx, rdb, gameId)
		if err != nil {
			return errors.New("Can't find that game")
		}

		game = RemovePlayer(game, playerName)

		if len(game.Players) == 1 && game.State == GamePlaying {
			game = EndGame(game)
		}

		SaveGame(ctx, rdb, gameId, game)

		persistentSession.GameHost = false
		persistentSession.PlayerName = playerName
		persistentSession.ActiveGame = ""
		persistentSession.UnsubChan = nil

		SetPersistentSession(ctx, session, rdb, persistentSession)
		SendGameResponse(session, cmd, gameId, game, true)

		return nil
	} else {
		return errors.New("Game not found")
	}

	return errors.New("Expected gameId and playerName to be supplied")
}

func startGame(ctx *context.Context, rdb *redis.Client, session *melody.Session, cmd *Command) error {
	persistentSession, err := GetPersistentSession(session)
	if err != nil {
		return errors.New("Something went wrong loading your session")
	}

	gameId := persistentSession.ActiveGame

	game, err := LoadGame(ctx, rdb, gameId)
	if err != nil {
		return errors.New("Error fetching game")
	}

	if !persistentSession.GameHost {
		return errors.New("Only the game host can start the game")
	}

	game = DrawHands(game)
	game = StartGame(game)

	SaveGame(ctx, rdb, gameId, game)
	SendGameResponse(session, cmd, gameId, game, false)

	return nil
}

func restartGame(ctx *context.Context, rdb *redis.Client, session *melody.Session, cmd *Command) error {
	persistentSession, err := GetPersistentSession(session)
	if err != nil {
		return errors.New("Something went wrong loading your session")
	}

	gameId := persistentSession.ActiveGame

	game, err := LoadGame(ctx, rdb, gameId)
	if err != nil {
		return errors.New("Error fetching game")
	}

	if !persistentSession.GameHost {
		return errors.New("Only the game host can restart the game")
	}

	game.State = GameCreated

	SaveGame(ctx, rdb, gameId, game)
	SendGameResponse(session, cmd, gameId, game, false)

	return nil
}

func endGame(ctx *context.Context, rdb *redis.Client, session *melody.Session, cmd *Command) error {
	persistentSession, err := GetPersistentSession(session)
	if err != nil {
		return errors.New("Something went wrong loading your session")
	}

	gameId := persistentSession.ActiveGame

	game, err := LoadGame(ctx, rdb, gameId)
	if err != nil {
		return errors.New("Error fetching game")
	}

	if !persistentSession.GameHost {
		return errors.New("Only the game host can end the game")
	}

	game = EndGame(game)

	SaveGame(ctx, rdb, gameId, game)
	SendGameResponse(session, cmd, gameId, game, true)

	return nil
}

func playCard(ctx *context.Context, rdb *redis.Client, session *melody.Session, cmd *Command) error {
	persistentSession, err := GetPersistentSession(session)
	if err != nil {
		return errors.New("Something went wrong loading your session")
	}

	cardIndex, ok1 := cmd.Data["cardIndex"]
	wildColor, ok2 := cmd.Data["wildColor"]

	if !ok1 && !ok2 {
		return errors.New("Expected cardIndex and wildColor to be supplied")
	}

	gameId := persistentSession.ActiveGame
	game, err := LoadGame(ctx, rdb, gameId)
	if err != nil {
		return errors.New("Error fetching game")
	}

	if game.ActivePlayer != GetPlayerIndex(game, persistentSession.PlayerName) {
		return errors.New("It's not your turn")
	}

	index, err := strconv.ParseInt(cardIndex, 10, 32)
	card := int(index)
	if err != nil || card < 0 || card > len(game.Players[game.ActivePlayer].Cards) {
		return errors.New("Invalid card index")
	}

	game, err = PlayCard(game, card, wildColor)
	if err != nil {
		return err
	}

	SaveGame(ctx, rdb, gameId, game)
	SendGameResponse(session, cmd, gameId, game, false)

	return nil
}

func drawCard(ctx *context.Context, rdb *redis.Client, session *melody.Session, cmd *Command) error {
	persistentSession, err := GetPersistentSession(session)
	if err != nil {
		return errors.New("Something went wrong loading your session")
	}

	gameId := persistentSession.ActiveGame
	game, err := LoadGame(ctx, rdb, gameId)
	if err != nil {
		return errors.New("Error fetching game")
	}

	if game.ActivePlayer != GetPlayerIndex(game, persistentSession.PlayerName) {
		return errors.New("It's not your turn")
	}

	game, err = DrawCard(game)
	if err != nil {
		return err
	}

	SaveGame(ctx, rdb, gameId, game)
	SendGameResponse(session, cmd, gameId, game, false)

	return nil
}

func doneDrawing(ctx *context.Context, rdb *redis.Client, session *melody.Session, cmd *Command) error {
	persistentSession, err := GetPersistentSession(session)
	if err != nil {
		return errors.New("Something went wrong loading your session")
	}

	gameId := persistentSession.ActiveGame
	game, err := LoadGame(ctx, rdb, gameId)
	if err != nil {
		return errors.New("Error fetching game")
	}

	if game.ActivePlayer != GetPlayerIndex(game, persistentSession.PlayerName) {
		return errors.New("It's not your turn")
	}

	game = DoneDrawing(game)

	SaveGame(ctx, rdb, gameId, game)
	SendGameResponse(session, cmd, gameId, game, false)

	return nil
}

// DispatchMessage handles an incoming game message
func DispatchMessage(ctx *context.Context, rdb *redis.Client, session *melody.Session, msg []byte) {
	cmd, err := parseCommand(msg)
	if err != nil {
		log.Println(err)
		return
	}

	switch cmd.Verb {
	case "openSession":
		log.Println("Opening session...")
		err = openSession(ctx, rdb, session, &cmd)
		break

	case "createGame":
		log.Println("Creating new game")
		err = createGame(ctx, rdb, session, &cmd)
		break

	case "joinGame":
		log.Println("Player joining a game")
		err = joinGame(ctx, rdb, session, &cmd)
		break

	case "leaveGame":
		log.Println("Player leaving a game")
		err = leaveGame(ctx, rdb, session, &cmd)
		break

	case "startGame":
		log.Println("Starting the game")
		err = startGame(ctx, rdb, session, &cmd)
		break

	case "restartGame":
		log.Println("Restarting the game")
		err = restartGame(ctx, rdb, session, &cmd)
		break

	case "endGame":
		log.Println("Ending the game")
		err = endGame(ctx, rdb, session, &cmd)
		break

	case "playCard":
		log.Println("Playing a card")
		err = playCard(ctx, rdb, session, &cmd)
		break

	case "drawCard":
		log.Println("Drawing a card")
		err = drawCard(ctx, rdb, session, &cmd)
		break

	case "doneDrawing":
		log.Println("Done drawing cards")
		err = doneDrawing(ctx, rdb, session, &cmd)
		break

	default:
		err = errors.New("Unrecognized command")
		break
	}

	if err != nil {
		sendResponse(session, errorResponse(&cmd, fmt.Sprint(err)))
	}
}
