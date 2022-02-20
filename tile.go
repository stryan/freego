package main

import (
	"errors"
	"fmt"
)

//Tile represents a spot on the board
type Tile struct {
	x        int
	y        int
	passable bool
	entity   *Piece
	colour   Colour
}

func (t *Tile) String() string {
	icon := "?"
	if !t.passable {
		icon = "X"
	} else if t.entity != nil {
		if !t.entity.Hidden {
			icon = t.entity.Rank.String()
		} else {
			icon = "O"
		}
	} else {
		icon = " "
	}
	return fmt.Sprintf("|%v|", icon)
}

//Place a piece on the tile
func (t *Tile) Place(p *Piece) error {
	if p == nil {
		return errors.New("tried to place nil piece")
	}
	if !t.passable {
		return errors.New("tile not passable")
	}
	if t.entity != nil {
		return errors.New("entity already present")
	}
	t.entity = p
	return nil
}

//Remove the piece from a tile
func (t *Tile) Remove() {
	t.entity = nil
}

//Empty returns true if tile is empty
func (t *Tile) Empty() bool {
	return t.entity == nil
}

//Team returns the team of the piece currently occupying it, or 0
func (t *Tile) Team() Colour {
	if t.entity != nil {
		return t.entity.Owner
	}
	return NoColour
}

//Occupied returns true if the tile is currently occupied
func (t *Tile) Occupied() bool {
	return t.entity != nil
}

//Piece returns the current Piece if any
func (t *Tile) Piece() *Piece {
	return t.entity
}

//Passable returns if the tile is passable
func (t *Tile) Passable() bool {
	return t.passable
}

//Colour returns the tile colour
func (t *Tile) Colour() Colour {
	return t.colour
}

//X returns tile X position
func (t *Tile) X() int {
	return t.x
}

//Y returns tile y position
func (t *Tile) Y() int {
	return t.y
}

//AddTerrain adds specified terrain to position
func (t *Tile) AddTerrain(ter int) bool {
	if t.entity != nil {
		return false
	}
	t.passable = false
	return true
}
