package main

import "testing"

func TestNewBoard(t *testing.T) {
	b := NewBoard(2)
	if b != nil {
		t.Error("Able to create too small board")
	}
	for size := 4; size < 9; size = size + 2 {
		b = NewBoard(size)
		if len(b.board[0]) != size {
			t.Errorf("board col is not size specified: %v != %v", 8, len(b.board[0]))
		}
		if len(b.board) != size {
			t.Errorf("board row is not size specified: %v != %v", 8, len(b.board[0]))
		}
		rTiles := 0
		bTiles := 0
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				ti := b.board[i][j]
				if ti == nil {
					t.Errorf("no tile at pos %v,%v", i, j)
				}
				if ti.entity != nil {
					t.Errorf("new board but %v,%v has a piece", i, j)
				}
				if !ti.passable {
					t.Errorf("New board but %v %v is inpassible", i, j)
				}
				if ti.colour == Red {
					rTiles++
				} else if ti.colour == Blue {
					bTiles++
				} else {
					t.Errorf("Tile %v,%v has invalid color %v", i, j, ti.colour)
				}
			}
		}
		if rTiles != bTiles {
			t.Errorf("Inequal number of coloured tiles: red tiles %v x blue tiles %v", rTiles, bTiles)
		}
	}

}

func TestGetPiece(t *testing.T) {
	b := NewBoard(4)
	p := NewPiece(2, Red)
	b.board[0][0].entity = p
	np, err := b.GetPiece(9, 9)
	if err == nil {
		t.Error("able to get piece from invalid location")
	}
	np, err = b.GetPiece(0, 0)
	if err != nil {
		t.Fatalf("GetPiece failed when it shouldn't: %v", err)
	}
	if np == nil {
		t.Errorf("Failed to get piece")
	}
	if np != p {
		t.Errorf("Got different piece from one placed")
	}

}

func TestPlace(t *testing.T) {
	b := NewBoard(4)
	p := NewPiece(2, Red)
	res, err := b.Place(8, 9, p)
	if err == nil {
		t.Errorf("able to place in invalid location")
	}
	res, err = b.Place(0, 0, p)
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Errorf("Failed to place piece")
	}
	p2, err := b.GetPiece(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if p2 == nil {
		t.Errorf("failed to get placed piece")
	}
	if p2 != p {
		t.Errorf("got different piece from one placed")
	}
}

func TestRemove(t *testing.T) {
	b := NewBoard(4)
	p := NewPiece(2, Red)

	res, err := b.Place(0, 0, p)
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Errorf("Failed to place piece")
	}
	p2, err := b.GetPiece(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if p2 == nil {
		t.Errorf("failed to get placed piece")
	}
	if p2 != p {
		t.Errorf("got different piece from one placed")
	}
	err = b.Remove(-1, 9)
	if err == nil {
		t.Errorf("able to remove from invalid location")
	}
	err = b.Remove(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if b.board[0][0].entity != nil {
		t.Fatalf("Ran remove but piece remained")
	}
}

func TestGetColor(t *testing.T) {
	b := NewBoard(4)
	if b.GetColor(0, 0) != Blue {
		t.Errorf("got wrong color for tile: %v", b.GetColor(0, 0))
	}
}

func TestAddTerrain(t *testing.T) {
	b := NewBoard(4)
	res, err := b.AddTerrain(5, 6, 1)
	if err == nil {
		t.Errorf("added terrain to invalid location")
	}
	res, err = b.AddTerrain(0, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Errorf("unable to add terrain")
	}
	if b.board[0][0].passable {
		t.Errorf("tile still passable even when there's terrain")
	}
}
