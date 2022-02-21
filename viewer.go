package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

//Viewer is a graphical representation of a freego game
type Viewer struct {
	input      *Input
	board      *ViewBoard
	boardImage *ebiten.Image
	gameState  *Game
}

//NewViewer creates a new viewer
func NewViewer(g *Game) (*Viewer, error) {
	v := &Viewer{
		input:     NewInput(),
		gameState: g,
	}
	var err error
	v.board, err = NewViewBoard(g)
	if err != nil {
		panic(err)
	}
	return v, nil
}

// Layout implements ebiten.Game's Layout.
func (v *Viewer) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 640
}

// Update updates the current game state.
func (v *Viewer) Update() error {
	v.input.Update()
	if err := v.board.Update(v.input); err != nil {
		return err
	}
	return nil
}

// Draw draws the current game to the given screen.
func (v *Viewer) Draw(screen *ebiten.Image) {
	if v.boardImage == nil {
		w, h := v.board.Size()
		v.boardImage = ebiten.NewImage(w, h)
	}
	screen.Fill(color.RGBA{0xfa, 0xf8, 0xef, 0xff})
	v.board.Draw(v.boardImage)
	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()
	bw, bh := v.boardImage.Size()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(v.boardImage, op)
}
