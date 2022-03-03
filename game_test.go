package freego

import (
	"errors"
	"fmt"
	"testing"
)

func dummyMiniGame() (*Game, error) {
	g := &Game{
		Board: NewBoard(4),
		State: gameSetup,
	}
	//Setup terrain
	terrain := []struct {
		x, y, t int
	}{
		{1, 1, 1},
		{2, 2, 1},
	}
	for _, tt := range terrain {
		res, err := g.Board.AddTerrain(tt.x, tt.y, tt.t)
		if err != nil {
			return nil, err
		}
		if !res {
			return nil, errors.New("Error creating terrain")
		}
	}
	pieces := []struct {
		x, y int
		p    *Piece
	}{
		{0, 0, NewPiece(Flag, Blue)},
		{3, 0, NewPiece(Spy, Blue)},
		{2, 0, NewPiece(Captain, Blue)},
		{3, 1, NewPiece(Marshal, Blue)},
		{0, 1, NewPiece(Bomb, Blue)},

		{1, 2, NewPiece(Flag, Red)},
		{3, 2, NewPiece(Spy, Red)},
		{2, 3, NewPiece(Captain, Red)},
		{0, 2, NewPiece(Marshal, Red)},
		{0, 3, NewPiece(Bomb, Red)},
	}
	for _, tt := range pieces {
		res, err := g.SetupPiece(tt.x, tt.y, tt.p)
		if err != nil {
			return nil, fmt.Errorf("Piece %v,%v:%v", tt.x, tt.y, err)
		}
		if !res {
			return nil, errors.New("error placing dummy piece")
		}
	}
	_, err := g.SetupPiece(0, 0, NewPiece(Flag, Blue))
	if err != nil {
		return nil, err
	}

	return g, nil
}

func TestNewGame(t *testing.T) {
	g := NewGame()
	if g == nil {
		t.Fatal("couldn't create game")
	}
	if g.State != gameLobby {
		t.Error("Game created with weird state")
	}
}

func TestSetupPiece(t *testing.T) {
	g := NewGame()
	p1 := NewPieceFromInt(0, Blue)
	p2 := NewPieceFromInt(0, Red)
	res, err := g.SetupPiece(0, 0, p1)
	if err == nil || res == true {
		t.Errorf("Expected to fail setup piece but didn't")
	}
	g.State = gameSetup
	res, err = g.SetupPiece(0, 0, p2)
	if err == nil || res == true {
		t.Error("Expected to fail putting red piece on blue board")
	}
	res, err = g.SetupPiece(0, 0, p1)
	if err != nil {
		t.Fatal(err)
	}
	res, err = g.SetupPiece(9, 9, p2)
	if err == nil {
		t.Error("expected to fail setting up piece in invalid spot")
	}
	res, err = g.SetupPiece(4, 5, p2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStart(t *testing.T) {
	g := NewGame()
	res := g.Start()
	if res {
		t.Fatal("expected to fail starting game due to state")
	}
	g.State = gameSetup
	res = g.Start()
	if !res {
		t.Fatal("expected game to start")
	}
}

func TestMiniCreation(t *testing.T) {
	g, err := dummyMiniGame()
	if err != nil {
		t.Fatalf("mini game not created: %v", err)
	}
	if g.State != gameSetup {
		t.Fatal("mini game not in right state")
	}
}

func TestMiniGameDemo(t *testing.T) {
	g, err := dummyMiniGame()
	if err != nil {
		t.Fatal(err)
	}
	if !g.Start() {
		t.Fatal("could not start mini game")
	}
	if g.State != gameTurnRed {
		t.Error("game starting on wrong turn")
	}
	r := NewDummyPlayer(Red)
	b := NewDummyPlayer(Blue)
	var moves = []struct {
		input  string
		player Player
		res    bool
	}{
		{"c3-b3", r, true},
		{"c0-c1", b, true},
		{"d2xd1", r, true},
		{"c1-d1", b, true},
		{"b3-c3", r, true},
		{"d1xd2", b, true},
	}
	i := 0
	for _, tt := range moves {
		tname := fmt.Sprintf("%v input: %v", i, tt.input)
		i++
		t.Run(tname, func(t *testing.T) {
			raw, err := NewRawCommand(tt.input)
			if err != nil {
				t.Fatalf("error creating RawCommand from %v:%v", tt.input, err)
			}
			parsed, err := g.Parse(tt.player, raw)
			if err != nil {
				t.Fatal(err)
			}
			res, err := g.Mutate(parsed)
			if err != nil {
				t.Error(err)
			}
			if tt.res != res {
				t.Errorf("Expected command to be %v but was %v", tt.res, res)
			}
		})
	}
}

func TestMiniGameFull(t *testing.T) {
	g, err := dummyMiniGame()
	if err != nil {
		t.Fatal(err)
	}
	if !g.Start() {
		t.Fatal("could not start mini game")
	}
	if g.State != gameTurnRed {
		t.Error("game starting on wrong turn")
	}
	r := NewDummyPlayer(Red)
	b := NewDummyPlayer(Blue)
	var moves = []struct {
		input  string
		player Player
		res    bool
	}{
		{"c3-b3", r, true},
		{"c0-c1", b, true},
		{"d2xd1", r, true},
		{"c1-d1", b, true},
		{"b3-c3", r, true},
		{"d1xd2", b, true},
	}
	i := 0
	for _, tt := range moves {
		tname := fmt.Sprintf("%v input: %v", i, tt.input)
		i++
		t.Run(tname, func(t *testing.T) {
			raw, err := NewRawCommand(tt.input)
			if err != nil {
				t.Fatalf("error creating RawCommand from %v:%v", tt.input, err)
			}
			parsed, err := g.Parse(tt.player, raw)
			if err != nil {
				t.Fatal(err)
			}
			res, err := g.Mutate(parsed)
			if err != nil {
				t.Error(err)
			}
			if tt.res != res {
				t.Errorf("Expected command to be %v but was %v", tt.res, res)
			}
		})
	}
}
