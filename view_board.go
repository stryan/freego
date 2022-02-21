package main

import (
	"errors"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var errorTaskTerminated = errors.New("freego: task terminated")

type task func() error

// ViewBoard represents the game board.
type ViewBoard struct {
	size  int
	tiles map[*ViewTile]struct{}
	tasks []task
}

// NewViewBoard generates a new ViewBoard with giving a size.
func NewViewBoard(g *Game) (*ViewBoard, error) {
	size := g.board.size
	b := &ViewBoard{
		size:  size,
		tiles: map[*ViewTile]struct{}{},
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			b.tiles[NewViewTile(i, j, g.board.board[i][j])] = struct{}{}
		}
	}

	return b, nil
}

//func (b *ViewBoard) tileAt(x, y int) *Tile {
//	return tileAt(b.tiles, x, y)
//}

// Update updates the board state.
func (b *ViewBoard) Update(input *Input) error {
	for t := range b.tiles {
		if err := t.Update(); err != nil {
			return err
		}
	}
	if 0 < len(b.tasks) {
		t := b.tasks[0]
		if err := t(); err == errorTaskTerminated {
			b.tasks = b.tasks[1:]
		} else if err != nil {
			return err
		}
		return nil
	}
	if dir, ok := input.Dir(); ok {
		//if err := b.Move(dir); err != nil {
		//	return err
		//}
		log.Println(dir)
	}
	return nil
}

// Draw draws the board to the given boardImage.
func (b *ViewBoard) Draw(boardImage *ebiten.Image) {
	boardImage.Fill(color.RGBA{0xbb, 0xad, 0xa0, 0xff})
	for j := 0; j < b.size; j++ {
		for i := 0; i < b.size; i++ {
			//v := 0
			op := &ebiten.DrawImageOptions{}
			x := i*tileSize + (i+1)*tileMargin
			y := j*tileSize + (j+1)*tileMargin
			op.GeoM.Translate(float64(x), float64(y))
			//op.ColorM.ScaleWithColor(tileBackgroundColor(v))
			boardImage.DrawImage(tileImage, op)
		}
	}
	for t := range b.tiles {
		t.Draw(boardImage)
	}
}

// Size returns the board size.
func (b *ViewBoard) Size() (int, int) {
	x := b.size*tileSize + (b.size+1)*tileMargin
	y := x
	return x, y
}
