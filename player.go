package main

//Player represents a player of the game
type Player interface {
	ID() int
	Colour() Colour
}

//DummyPlayer is the simplest implementation of a player
type DummyPlayer struct {
	id Colour
}

//ID returns player ID
func (d *DummyPlayer) ID() int {
	return int(d.id)
}

//Colour returns player colour (same as ID)
func (d *DummyPlayer) Colour() Colour {
	return d.id
}

//NewDummyPlayer creates a new player
func NewDummyPlayer(i Colour) *DummyPlayer {
	return &DummyPlayer{i}
}

func (d *DummyPlayer) String() string {
	if d.id == Red {
		return "Red Player"
	} else {
		return "Blue Player"
	}
}
