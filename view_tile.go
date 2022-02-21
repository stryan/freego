package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//ViewTile is a tile in the viewer
type ViewTile struct {
	gameTile *Tile

	// next represents a next tile information after moving.
	// next is empty when the tile is not about to move.
	//next TileData

	//movingCount       int
	//startPoppingCount int
	//poppingCount      int
}

var (
	tileImage = ebiten.NewImage(tileSize, tileSize)
)

var (
	mplusSmallFont  font.Face
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

const (
	tileSize   = 80
	tileMargin = 4
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusSmallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

//NewViewTile creates a new view tile
func NewViewTile(x, y int, t *Tile) *ViewTile {
	return &ViewTile{
		gameTile: t,
	}
}

// Update updates the tile's animation states.
func (t *ViewTile) Update() error {
	return nil
}

// Draw draws the current tile to the given boardImage.
func (t *ViewTile) Draw(boardImage *ebiten.Image) {
	j, i := t.gameTile.x, t.gameTile.y
	v := t.gameTile.entity
	if v == nil && t.gameTile.Passable() {
		return
	}
	op := &ebiten.DrawImageOptions{}
	x := i*tileSize + (i+1)*tileMargin
	y := j*tileSize + (j+1)*tileMargin
	op.GeoM.Translate(float64(x), float64(y))
	//op.ColorM.ScaleWithColor(tileBackgroundColor(v))
	boardImage.DrawImage(tileImage, op)
	str := "NA"
	if v == nil {
		str = "river"
	} else {
		str = v.Rank.String()
	}

	f := mplusBigFont
	switch {
	case 3 < len(str):
		f = mplusSmallFont
	case 2 < len(str):
		f = mplusNormalFont
	}

	bound, _ := font.BoundString(f, str)
	w := (bound.Max.X - bound.Min.X).Ceil()
	h := (bound.Max.Y - bound.Min.Y).Ceil()
	x = x + (tileSize-w)/2
	y = y + (tileSize-h)/2 + h
	pieceColor := color.RGBA{0xf9, 0xf6, 0xf2, 0xff}
	if v != nil {
		if t.gameTile.entity.Owner == Red {
			pieceColor = color.RGBA{0xff, 0x0, 0x0, 0xff}
		} else if t.gameTile.entity.Owner == Blue {
			pieceColor = color.RGBA{0x0, 0x0, 0xff, 0xff}
		}
	}
	text.Draw(boardImage, str, f, x, y, pieceColor)
}
