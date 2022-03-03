package freego

import (
	"errors"
	"fmt"
)

//DummyGame Creates a dummygame
func DummyGame() (*Game, error) {
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

//func main() {
//	//red := NewDummyPlayer(Red)
//blue := NewDummyPlayer(Blue)
//	_, err := DummyGame()
//	if err != nil {
//		panic(err)
//	}
//	return
//}

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
	res, err := game.Board.AddTerrain(x, y, 1)
	if err != nil {
		panic(err)
	}
	if !res {
		panic("can't river")
	}
}

func printboardcolours(g *Game) {
	for i := range g.Board.board {
		for j := range g.Board.board[i] {
			c := "X"
			if g.Board.board[i][j].colour == Red {
				c = "R"
			} else if g.Board.board[i][j].colour == Blue {
				c = "B"
			}
			fmt.Printf("%v", c)
		}
		fmt.Println("")
	}
}
