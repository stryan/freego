package main

import "testing"

func dummyMiniGame() *Game {
	g := &Game{
		board: NewBoard(4),
		state: gameSetup,
	}
	//Setup terrain
	g.board.AddTerrain(1, 1, 1)
	g.board.AddTerrain(2, 2, 1)
	//Setup blue (5 pieces)
	g.SetupPiece(0, 0, NewPiece(0, Blue))
	g.SetupPiece(0, 1, NewPiece(1, Blue))
	g.SetupPiece(0, 2, NewPiece(2, Blue))
	g.SetupPiece(0, 3, NewPiece(3, Blue))
	g.SetupPiece(1, 3, NewPiece(4, Blue))
	//Setup red (5 pieces)
	g.SetupPiece(3, 0, NewPiece(0, Red))
	g.SetupPiece(3, 1, NewPiece(1, Red))
	g.SetupPiece(3, 2, NewPiece(2, Red))
	g.SetupPiece(3, 3, NewPiece(3, Red))
	g.SetupPiece(2, 1, NewPiece(4, Red))

	return g
}

func TestNewGame(t *testing.T) {
	g := NewGame()
	if g == nil {
		t.Fatal("couldn't create game")
	}
	if g.state != gameLobby {
		t.Error("Game created with weird state")
	}
}

func TestSetupPiece(t *testing.T) {
	g := NewGame()
	p1 := NewPiece(0, Blue)
	p2 := NewPiece(0, Red)
	res, err := g.SetupPiece(0, 0, p1)
	if err == nil || res == true {
		t.Errorf("Expected to fail setup piece but didn't")
	}
	g.state = gameSetup
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
	res, err = g.SetupPiece(4, 0, p2)
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
	g.state = gameSetup
	res = g.Start()
	if !res {
		t.Fatal("expected game to start")
	}
}
