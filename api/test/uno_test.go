package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/go-test/deep"
)

func startedGame() Game {
	rand.Seed(0)

	game := EmptyGame()
	game = AddPlayer(game, "Nia")
	game = AddPlayer(game, "Eric")
	game = DrawHands(game)

	game = StartGame(game)

	return game
}

func printGame(game Game) {
	b, _ := json.Marshal(game)
	fmt.Printf("%s\n", b)
}

func TestEmptyGame(t *testing.T) {
	rand.Seed(0)

	game := EmptyGame()

	expected := Game{
		State:         GameCreated,
		Players:       []Player{},
		ActivePlayer:  0,
		GameDirection: 1,
		DrawPile:      []string{"B4", "R4", "B0", "G+2", "R5", "Yrev", "G6", "R1", "R+2", "R3", "B8", "Y+2", "Gskip", "Y8", "Grev", "G8", "R+2", "Y2", "Yskip", "wild", "wild", "R0", "B6", "B2", "Y9", "R3", "B9", "wild+4", "Bskip", "B+2", "R6", "wild+4", "R9", "wild", "B5", "G2", "R4", "G7", "G2", "R6", "Grev", "R9", "B8", "Y3", "B3", "Y0", "B5", "B1", "Yskip", "B1", "Y1", "G5", "G8", "B9", "B3", "Y1", "Y4", "R2", "B6", "Y+2", "B2", "Y7", "wild+4", "Bskip", "G5", "Gskip", "R7", "G1", "Y6", "Y2", "Yrev", "G+2", "wild+4", "Y5", "Y4", "R2", "Rskip", "R8", "G9", "B4", "Y8", "G7", "Y9", "R8", "Rrev", "G9", "G3", "G6", "Y7", "Brev", "R1", "B+2", "Rskip", "B7", "G0", "Y6", "Y3", "G3", "R7", "Y5", "B7", "Rrev", "G4", "G4", "R5", "Brev", "G1", "wild"},
		DiscardPile:   []string{},
	}

	if diff := deep.Equal(game, expected); diff != nil {
		t.Error(diff)
	}
}

func TestAddPlayer(t *testing.T) {
	rand.Seed(0)

	game := EmptyGame()

	game = AddPlayer(game, "Nia")
	game = AddPlayer(game, "Eric")

	expected := []Player{
		Player{Name: "Nia", Cards: []string{}},
		Player{Name: "Eric", Cards: []string{}},
	}

	if diff := deep.Equal(game.Players, expected); diff != nil {
		t.Error(diff)
	}
}

func TestDrawHands(t *testing.T) {
	rand.Seed(0)

	game := EmptyGame()

	game = AddPlayer(game, "Nia")
	game = AddPlayer(game, "Eric")

	game = DrawHands(game)

	expectedPlayers := []Player{
		Player{Name: "Nia", Cards: []string{"B4", "R4", "B0", "G+2", "R5", "Yrev", "G6"}},
		Player{Name: "Eric", Cards: []string{"R1", "R+2", "R3", "B8", "Y+2", "Gskip", "Y8"}},
	}

	expectedDiscard := []string{"Grev"}

	if diff := deep.Equal(game.Players, expectedPlayers); diff != nil {
		t.Error(diff)
	}

	if diff := deep.Equal(game.DiscardPile, expectedDiscard); diff != nil {
		t.Error(diff)
	}
}

func TestStartGame(t *testing.T) {
	game := startedGame()

	if game.State != GamePlaying {
		t.Error("Expected game state to be Playing")
	}

	if game.ActivePlayer != 0 {
		t.Error("Expected player 0 to be active")
	}
}

func TestPlayCard(t *testing.T) {
	game := startedGame()

	game, err := PlayCard(game, 6, "")

	if err != nil {
		t.Error("Should not have returned error")
	}

	expectedPlayerCards := []string{
		"B4", "R4", "B0", "G+2", "R5", "Yrev",
	}
	if diff := deep.Equal(game.Players[0].Cards, expectedPlayerCards); diff != nil {
		t.Error(diff)
	}

	if game.DiscardPile[0] != "G6" {
		t.Error("Top of discard pile should be G6")
	}
}

func TestValidateCardPlay(t *testing.T) {
	type CardPlay struct {
		topCard   string
		wildColor string
		newCard   string
		valid     bool
	}

	plays := []CardPlay{
		// Colors match
		{"Grev", "", "Grev", true},
		{"Grev", "", "G2", true},
		{"Grev", "", "R2", false},

		{"Brev", "", "Brev", true},
		{"Brev", "", "B2", true},
		{"Brev", "", "R2", false},

		{"Rrev", "", "Rrev", true},
		{"Rrev", "", "R2", true},
		{"Rrev", "", "G2", false},

		{"Yrev", "", "Yrev", true},
		{"Yrev", "", "Y2", true},
		{"Yrev", "", "R2", false},

		// Wilds
		{"wild", "G", "G2", true},
		{"wild+4", "G", "G2", true},
		{"wild", "Y", "G2", false},
		{"wild+4", "Y", "G2", false},

		{"wild", "", "wild+4", true},
		{"wild+4", "", "wild", true},

		// Reverses
		{"Grev", "", "Yrev", true},
		{"Grev", "", "G+2", true},
		{"Grev", "", "Y+2", false},
		{"G+2", "", "Grev", true},

		// +2s
		{"G+2", "", "Y+2", true},
		{"G+2", "", "Grev", true},
		{"G+2", "", "Yrev", false},
		{"Grev", "", "G+2", true},

		// Skips
		{"Gskip", "", "Yskip", true},
		{"Gskip", "", "Grev", true},
		{"Gskip", "", "Yrev", false},
		{"Grev", "", "Gskip", true},
	}

	for _, play := range plays {
		err := ValidateCardPlay(play.topCard, play.wildColor, play.newCard)

		if play.valid && err != nil {
			t.Error(fmt.Sprintf("Expected %s played on %s (%s) to be valid", play.newCard, play.topCard, play.wildColor))
		}

		if !play.valid && err == nil {
			t.Error(fmt.Sprintf("Expected %s played on %s (%s) to be invalid", play.newCard, play.topCard, play.wildColor))
		}
	}
}

func TestAdvancePlayerPlusTwo(t *testing.T) {

	game := Game{
		State: GameCreated,
		Players: []Player{
			Player{Name: "0", Cards: []string{"R+2"}},
			Player{Name: "1", Cards: []string{"R+2"}},
			Player{Name: "2", Cards: []string{"R+2"}},
		},
		ActivePlayer:  0,
		MustDraw:      0,
		GameDirection: Clockwise,
		DrawPile:      []string{"B4", "R4", "B0", "G+2"},
		DiscardPile:   []string{"R0"},
	}

	// Play a +2
	game, err := PlayCard(game, 0, "")
	if err != nil {
		t.Error(fmt.Sprintf("Didn't expect an error: %s", err))
	}

	if !(game.ActivePlayer == 1 && game.MustDraw == 2) {
		t.Error("Expected player 1 to be active")
	}

	game, err = PlayCard(game, 0, "")
	if err == nil {
		t.Error("Expected an error. The player must draw")
	}

	game, err = DrawCard(game)
	game, err = DrawCard(game)

	if !(game.ActivePlayer == 2 && game.MustDraw == 0) {
		t.Error("Expected player 2 to be active")
	}
}

func TestAdvancePlayerPlusFour(t *testing.T) {

	game := Game{
		State: GameCreated,
		Players: []Player{
			Player{Name: "0", Cards: []string{"R+4"}},
			Player{Name: "1", Cards: []string{"R+4"}},
			Player{Name: "2", Cards: []string{"R+4"}},
		},
		ActivePlayer:  0,
		MustDraw:      0,
		GameDirection: Clockwise,
		DrawPile:      []string{"B4", "R4", "B0", "G+2"},
		DiscardPile:   []string{"R0"},
	}

	// Play a +4
	game, err := PlayCard(game, 0, "")
	if err != nil {
		t.Error(fmt.Sprintf("Didn't expect an error: %s", err))
	}

	if !(game.ActivePlayer == 1 && game.MustDraw == 4) {
		t.Error("Expected player 1 to be active")
	}

	game, err = PlayCard(game, 0, "")
	if err == nil {
		t.Error("Expected an error. The player must draw")
	}

	game, err = DrawCard(game)
	game, err = DrawCard(game)
	game, err = DrawCard(game)
	game, err = DrawCard(game)

	if !(game.ActivePlayer == 2 && game.MustDraw == 0) {
		t.Error("Expected player 2 to be active")
	}
}

func TestAdvancePlayerSkip(t *testing.T) {

	game := Game{
		State: GameCreated,
		Players: []Player{
			Player{Name: "0", Cards: []string{"Rskip"}},
			Player{Name: "1", Cards: []string{"Rskip"}},
			Player{Name: "2", Cards: []string{"Rskip"}},
		},
		ActivePlayer:  0,
		MustDraw:      0,
		GameDirection: Clockwise,
		DrawPile:      []string{"B4", "R4", "B0", "G+2"},
		DiscardPile:   []string{"R0"},
	}

	// Play a Skip
	game, err := PlayCard(game, 0, "")
	if err != nil {
		t.Error(fmt.Sprintf("Didn't expect an error: %s", err))
	}

	if !(game.ActivePlayer == 2 && game.MustDraw == 0) {
		t.Error("Expected player 2 to be active")
	}

	game, err = DrawCard(game)
	game = DoneDrawing(game)

	if !(game.ActivePlayer == 0 && game.MustDraw == 0) {
		t.Error("Expected player 0 to be active")
	}
}

func TestAdvancePlayerReverse(t *testing.T) {

	game := Game{
		State: GameCreated,
		Players: []Player{
			Player{Name: "0", Cards: []string{"Rrev"}},
			Player{Name: "1", Cards: []string{"Rrev"}},
			Player{Name: "2", Cards: []string{"Rrev"}},
		},
		ActivePlayer:  0,
		MustDraw:      0,
		GameDirection: Clockwise,
		DrawPile:      []string{"B4", "R4", "B0", "G+2"},
		DiscardPile:   []string{"R0"},
	}

	// Play a reverse
	game, err := PlayCard(game, 0, "")
	if err != nil {
		t.Error(fmt.Sprintf("Didn't expect an error: %s", err))
	}

	printGame(game)

	if !(game.ActivePlayer == 2 && game.MustDraw == 0) {
		t.Error("Expected player 2 to be active")
	}

}

// Play a Reverse
// Play a Wild
// Play a Wild+4
