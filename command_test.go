package main

import (
	"fmt"
	"testing"
)

func TestNewRawCommand(t *testing.T) {
	var tests = []struct {
		input string
		res   bool
		sx    int
		sy    int
		dx    int
		dy    int
		act   string
	}{
		{"A1xC3", true, 0, 1, 2, 3, "x"},
		{"d7-d1", true, 3, 7, 3, 1, "-"},
		{"AA-k3", false, 0, 0, 0, 0, ""},
		{"a1/b4", false, 0, 0, 0, 0, ""},
	}
	for _, tt := range tests {
		tname := fmt.Sprintf("input: %v", tt.input)
		t.Run(tname, func(t *testing.T) {
			res, err := NewRawCommand(tt.input)
			if tt.res && err != nil {
				t.Fatal(err)
			}
			if tt.res {
				if tt.sx != res.srcX {
					t.Errorf("bad sourceX: %v != %v", tt.sx, res.srcX)
				}
				if tt.sy != res.srcY {
					t.Errorf("bad sourceY: %v != %v", tt.sy, res.srcY)
				}
				if tt.dx != res.dstX {
					t.Errorf("bad destX: %v != %v", tt.dx, res.dstX)
				}
				if tt.dy != res.dstY {
					t.Errorf("bad destY: %v != %v", tt.dy, res.dstY)
				}
				if tt.act != res.act {
					t.Errorf("bad action: %v != %v", tt.act, res.act)
				}
			}
		})
	}
}
