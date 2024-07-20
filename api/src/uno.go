package main

import (
	"errors"
	"math/rand"
	"sort"
	"strings"
)

type GameState int

const (
	GameCreated   GameState = 0
	GamePlaying   GameState = 1
	GameComplete  GameState = 2
	GameAbandoned GameState = 3
)

type GameDirection int

const (
	Clockwise        GameDirection = 1
	CounterClockwise GameDirection = -1
)

type Player struct {
	Name  string   `json:"name"`
	Cards []string `json:"cards"`
}

type OtherPlayer struct {
	Name     string `json:"name"`
	NumCards int    `json:"numCards"`
}

type Game struct {
	GameCode      string
	GamePneumonic string

	State         GameState
	Players       []Player
	ActivePlayer  int
	GameDirection GameDirection
	MustDraw      int
	WildColor     string
	DrawPile      []string
	DiscardPile   []string
}

type PlayersGame struct {
	State         GameState     `json:"state"`
	ActivePlayer  int           `json:"activePlayer"`
	GameDirection GameDirection `json:"direction"`

	You          Player        `json:"you"`
	OtherPlayers []OtherPlayer `json:"otherPlayers"`

	MustDraw         int    `json:"mustDraw"`
	WildColor        string `json:"wildColor"`
	DrawPileCount    int    `json:"drawPileCount"`
	DiscardPileTop   string `json:"discardPileTop"`
	DiscardPileCount int    `json:"discardPileCount"`
}

type GameError struct {
	message string
}

func (e *GameError) Error() string {
	return e.message
}

var colors = [...]string{"R", "G", "B", "Y"}
var numbers = [...]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var modifiers = [...]string{"rev", "+2", "skip"}
var specials = [...]string{"wild", "wild", "wild", "wild", "wild+4", "wild+4", "wild+4", "wild+4"}

// GetPlayersGame returns a modifier Game object for a particular player, only showing information relevant
// to them
func GetPlayersGame(game *Game, playersIndex int) *PlayersGame {
	you := game.Players[playersIndex]

	// Collect the status of other players
	otherPlayers := make([]OtherPlayer, 0)
	for _, player := range game.Players {
		otherPlayers = append(otherPlayers, OtherPlayer{
			NumCards: len(player.Cards),
			Name:     player.Name,
		})
	}

	discardPileTop := ""
	if len(game.DiscardPile) > 0 {
		discardPileTop = game.DiscardPile[0]
	}

	// Build the player's game object
	return &PlayersGame{
		State:         game.State,
		ActivePlayer:  game.ActivePlayer,
		GameDirection: game.GameDirection,

		You:          you,
		OtherPlayers: otherPlayers,

		MustDraw:         game.MustDraw,
		WildColor:        game.WildColor,
		DrawPileCount:    len(game.DrawPile),
		DiscardPileTop:   discardPileTop,
		DiscardPileCount: len(game.DiscardPile),
	}
}

// GetPlayerIndex returns the index of the player with the given name, in the provided game
func GetPlayerIndex(game *Game, playerName string) int {
	for i, player := range game.Players {
		if player.Name == playerName {
			return i
		}
	}

	return -1
}

// Deck returns a new uno deck
func Deck() []string {
	var cards []string

	for _, color := range colors {
		for _, number := range numbers {
			cards = append(cards, color+number)
		}

		// Uno includes two sets of number cards per color, expect 0 (don't know why...)
		for _, number := range numbers[1:] {
			cards = append(cards, color+number)
		}

		for _, modifier := range modifiers {
			cards = append(cards, color+modifier)
			cards = append(cards, color+modifier)
		}
	}

	for _, special := range specials {
		cards = append(cards, special)
	}

	return cards
}

// Shuffle returns a new deck that has been shuffled
func Shuffle(deck []string) []string {
	newDeck := []string{}

	for _, card := range deck {
		newDeck = append(newDeck, card)
	}

	rand.Shuffle(
		len(newDeck),
		func(i, j int) { newDeck[i], newDeck[j] = newDeck[j], newDeck[i] },
	)

	return newDeck
}

// EmptyGame returns a new Game instance
func EmptyGame(gameCode string, gamePneumonic string) *Game {
	return &Game{
		GameCode:      gameCode,
		GamePneumonic: gamePneumonic,
		State:         GameCreated,
		ActivePlayer:  0,
		GameDirection: Clockwise,
		WildColor:     "R",
		DrawPile:      Shuffle(Deck()),
		DiscardPile:   []string{},
		Players:       []Player{},
	}
}

// AddPlayer returns a game with a new player added
func AddPlayer(game *Game, name string) *Game {
	game.Players = append(game.Players, Player{
		Name:  name,
		Cards: []string{},
	})

	return game
}

// RemovePlayer returns a game with the given player removed
func RemovePlayer(game *Game, name string) *Game {
	index := GetPlayerIndex(game, name)
	if index > 0 {
		// Release the players cards back into the deck
		game.DrawPile = append(game.DrawPile, Shuffle(game.Players[index].Cards)...)

		// Remove the player
		game.Players = append(game.Players[:index], game.Players[index+1:]...)

		// Ensure the current player isn't active
		if game.ActivePlayer == index {
			game = AdvancePlayer(game)
		}
	}

	return game
}

// Draw returns a game and an array with numCards cards
func Draw(game *Game, numCards int) (*Game, []string) {
	cards := []string{}

	for i := 0; i < numCards; i++ {
		if len(game.DrawPile) == 0 {
			game.DrawPile = Shuffle(game.DiscardPile[1:])
			game.DiscardPile = game.DiscardPile[0:1]
		}

		cards = append(cards, game.DrawPile[0])
		game.DrawPile = game.DrawPile[1:]
	}

	return game, cards
}

// DrawHands returns a game with all the players getting 7 cards from the draw pile
func DrawHands(game *Game) *Game {
	// Reset the game
	game.GameDirection = Clockwise
	game.DrawPile = Shuffle(Deck())
	game.DiscardPile = []string{}

	// Each player starts with 7 cards
	for i := range game.Players {
		// Draw 7 cards from the deck
		game, game.Players[i].Cards = Draw(game, 7)

		// Sort the player cards after adding
		game = SortPlayerCards(game, i)
	}

	// Draw a card for the discard pile
	game, game.DiscardPile = Draw(game, 1)
	game.WildColor = colors[rand.Intn(len(colors))]

	return game
}

// StartGame returns a game which has been started
func StartGame(game *Game) *Game {
	game.State = GamePlaying
	game.ActivePlayer = rand.Intn(len(game.Players))

	return game
}

// EndGame returns a game which has been started
func EndGame(game *Game) *Game {
	game.State = GameAbandoned

	return game
}

// RemovePlayerCard removes a card from a players hand and returns it
func RemovePlayerCard(game *Game, playerIndex int, cardIndex int) (*Game, string) {
	playerCards := game.Players[playerIndex].Cards
	card := playerCards[cardIndex]

	game.Players[playerIndex].Cards = append(
		playerCards[:cardIndex],
		playerCards[cardIndex+1:]...,
	)

	return game, card
}

// DiscardCard places the given card on top of the discard pile
func DiscardCard(game *Game, card string) *Game {
	game.DiscardPile = append([]string{card}, game.DiscardPile...)
	return game
}

// ValidateCardPlay player returns an error if the give play is invalid. Otherwise
// it returns nil.
func ValidateCardPlay(topCard string, wildColor string, card string) error {
	// If the new card is a wild card (which can be played on anything)
	if strings.HasPrefix(card, "wild") {
		return nil
	}

	// If the new card is being placed on a wild
	if strings.HasPrefix(topCard, "wild") && len(wildColor) == 1 {
		// Playing a wild on a wild
		if strings.HasPrefix(card, "wild") {
			return nil
		}

		// Otherwise, the colors must match
		if card[0] == wildColor[0] {
			return nil
		}
	}

	// If the colors match
	if card[0] == topCard[0] {
		return nil
	}

	// Both cards are a reverse
	if strings.HasSuffix(card, "rev") && strings.HasSuffix(topCard, "rev") {
		return nil
	}

	// Both cards are a +2
	if strings.HasSuffix(card, "+2") && strings.HasSuffix(topCard, "+2") {
		return nil
	}

	// Both cards are a skip
	if strings.HasSuffix(card, "skip") && strings.HasSuffix(topCard, "skip") {
		return nil
	}

	// Otherwise, both cards are number cards. Ensure the numbers match
	if card[1] == topCard[1] {
		return nil
	}

	return &GameError{message: "Can't play that card"}
}

func reverseDirection(direction GameDirection) GameDirection {
	if direction == Clockwise {
		return CounterClockwise
	}

	return Clockwise
}

// ApplyModifiers should be called after the player plays a card on the discard pile. It updates
// the ActivePlayer wit the new active player
func ApplyModifiers(game *Game) *Game {
	topCard := game.DiscardPile[0]

	// if MustDraw > 0, then the subsequent players turn must be used to draw 2 from the pile
	if strings.HasSuffix(topCard, "+2") {
		game.MustDraw = 2
	} else if strings.HasSuffix(topCard, "+4") {
		game.MustDraw = 4
	}

	if strings.HasSuffix(topCard, "skip") {
		game = AdvancePlayer(game)
	}

	if strings.HasSuffix(topCard, "rev") {
		game.GameDirection = reverseDirection(game.GameDirection)

		// People expect a reverse to skip the next player (those that's now really how it works...)
		if len(game.Players) == 2 {
			game = AdvancePlayer(game)
		}
	}

	return game
}

// AdvancePlayer increments the ActivePlayer by 1 in the current GameDirection
func AdvancePlayer(game *Game) *Game {
	newPlayer := (int(game.ActivePlayer) + int(game.GameDirection))

	if newPlayer >= len(game.Players) {
		newPlayer -= len(game.Players)
	}

	if newPlayer < 0 {
		newPlayer += len(game.Players)
	}

	game.ActivePlayer = newPlayer

	return game
}

// CheckForWinner checks if any player has 0 cards, and if so sets the state to complete
func CheckForWinner(game *Game) *Game {
	for _, player := range game.Players {
		if len(player.Cards) == 0 {
			game.State = GameComplete
			return game
		}
	}

	return game
}

// PlayCard takes a card from the active player's hand and places it on the top of the discard pile
func PlayCard(game *Game, cardIndex int, wildColor string) (*Game, error) {
	topCard := game.DiscardPile[0]
	playerIndex := game.ActivePlayer
	cardToPlay := game.Players[playerIndex].Cards[cardIndex]

	// Check if the player is allows to play
	if game.MustDraw > 0 {
		return game, errors.New("player cannot play, they must draw")
	}

	if len(wildColor) > 0 {
		game.WildColor = wildColor
	}

	// Validate the card move
	if err := ValidateCardPlay(topCard, game.WildColor, cardToPlay); err != nil {
		return game, err
	}

	// Remove the players card from their hand, and discard it on top of the DiscardPile
	game, playedCard := RemovePlayerCard(game, playerIndex, cardIndex)
	game = DiscardCard(game, playedCard)

	// Apply modifier cards (+2, +4, skip, reverse)
	game = ApplyModifiers(game)

	// Move to the next player
	game = AdvancePlayer(game)

	// Check if any player has won
	game = CheckForWinner(game)

	return game, nil
}

func SortPlayerCards(game *Game, playerIndex int) *Game {
	// Sort the new cards
	sort.Slice(game.Players[playerIndex].Cards, func(a, b int) bool {
		return game.Players[playerIndex].Cards[a] < game.Players[playerIndex].Cards[b]
	})

	return game
}

// DrawCard draws a card from the DrawPile and places it in the active players hand
func DrawCard(game *Game) (*Game, error) {
	playerIndex := game.ActivePlayer
	currentCards := game.Players[playerIndex].Cards

	// Draw one card from the deck
	newGame, newCards := Draw(game, 1)
	newGame.Players[playerIndex].Cards = append(currentCards, newCards[0])

	// Sort the player cards after adding
	newGame = SortPlayerCards(game, playerIndex)

	if newGame.MustDraw > 0 {
		newGame.MustDraw--

		if newGame.MustDraw == 0 {
			newGame = DoneDrawing(newGame)
		}
	}

	return newGame, nil
}

// DoneDrawing indicates that the current player is done drawing, and the next person should play
func DoneDrawing(game *Game) *Game {
	return AdvancePlayer(game)
}
