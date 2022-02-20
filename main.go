package main

import "fmt"

func main() {
	red := NewDummyPlayer(Red)
	blue := NewDummyPlayer(Blue)
	game := NewGame()
	game.state = gameSetup
	addpiece(game, 0, Blue, 0, 0)
	addpiece(game, 1, Blue, 0, 1)
	addpiece(game, 2, Blue, 1, 2)
	addpiece(game, 3, Blue, 2, 4)
	addpiece(game, 3, Blue, 1, 3)
	addpiece(game, 4, Blue, 2, 6)
	addpiece(game, 5, Blue, 3, 3)
	addpiece(game, 9, Blue, 0, 6)
	addpiece(game, 9, Blue, 0, 7)
	addpiece(game, 0, Blue, 3, 5)
	addriver(game, 3, 2)
	addriver(game, 4, 2)
	addriver(game, 3, 6)
	addriver(game, 4, 6)
	addpiece(game, 0, Red, 4, 0)
	addpiece(game, 1, Red, 4, 1)
	addpiece(game, 2, Red, 5, 2)
	addpiece(game, 3, Red, 6, 4)
	addpiece(game, 3, Red, 5, 3)
	addpiece(game, 4, Red, 6, 6)
	addpiece(game, 5, Red, 7, 3)
	addpiece(game, 9, Red, 4, 5)
	addpiece(game, 9, Red, 4, 7)
	addpiece(game, 0, Red, 7, 5)

	fmt.Printf("%v %v\n%v\n", red, blue, game)
	game.Start()

	return
}

func addpiece(game *Game, rank int, c Colour, x int, y int) {
	res, err := game.SetupPiece(x, y, NewPiece(rank, c))
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
