package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func dummyGame() (*Game, error) {
	g := &Game{
		board: NewBoard(4),
		state: gameSetup,
	}
	//Setup terrain
	terrain := []struct {
		x, y, t int
	}{
		{1, 1, 1},
		{2, 2, 1},
	}
	for _, tt := range terrain {
		res, err := g.board.AddTerrain(tt.x, tt.y, tt.t)
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
func main() {
	//red := NewDummyPlayer(Red)
	//blue := NewDummyPlayer(Blue)
	g, err := dummyGame()
	if err != nil {
		panic(err)
	}
	go func(g *Game) {
		r := NewDummyPlayer(Red)
		b := NewDummyPlayer(Blue)
		g.Start()
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
		for i, tt := range moves {
			time.Sleep(3 * time.Second)
			log.Printf("playing move %v\n", i)
			raw, err := NewRawCommand(tt.input)
			if err != nil {
				panic(err)
			}
			parsed, err := g.Parse(tt.player, raw)
			if err != nil {
				panic(err)
			}
			res, err := g.Mutate(parsed)
			if err != nil {
				panic(err)
			}
			if res {
				log.Printf("move %v successful\n", i)
			}

		}

	}(g)
	viewer(g)
	return
}

func viewer(g *Game) {
	v, err := NewViewer(g)
	if err != nil {
		panic(err)
	}
	ebiten.SetWindowSize(800, 640)
	ebiten.SetWindowTitle("Freego")
	if err := ebiten.RunGame(v); err != nil {
		panic(err)
	}
}

func addpiece(game *Game, rank int, c Colour, x int, y int) {
	res, err := game.SetupPiece(x, y, NewPieceFromInt(rank, c))
	if err != nil {
		panic(err)
	}
	if !res {
		panic("can't setup")
	}
}

func addriver(game *Game, x int, y int) {
	res, err := game.board.AddTerrain(x, y, 1)
	if err != nil {
		panic(err)
	}
	if !res {
		panic("can't river")
	}
}

func printboardcolours(g *Game) {
	for i := range g.board.board {
		for j := range g.board.board[i] {
			c := "X"
			if g.board.board[i][j].colour == Red {
				c = "R"
			} else if g.board.board[i][j].colour == Blue {
				c = "B"
			}
			fmt.Printf("%v", c)
		}
		fmt.Println("")
	}
}
