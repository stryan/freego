package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var cmdRegxp = regexp.MustCompile("([a-zA-Z])([1-9])(x|-)([a-zA-Z])([1-9])")
var ranks = "ABCDEFHI"

//RawCommand is a game command, converted from algebraic notation
type RawCommand struct {
	srcX int
	srcY int
	dstX int
	dstY int
	act  string
}

//NewRawCommand creates a RawCommand struct from a algebraic notation string
func NewRawCommand(cmd string) (*RawCommand, error) {
	res := cmdRegxp.FindAllString(cmd, -1)
	if res == nil {
		return nil, errors.New("error creating command from string")
	}
	sx := strings.Index(ranks, res[0])
	dx := strings.Index(ranks, res[3])
	sy, err := strconv.Atoi(res[1])
	if err != nil {
		return nil, err
	}
	dy, err := strconv.Atoi(res[4])
	if err != nil {
		return nil, err
	}
	return &RawCommand{sx, sy, dx, dy, res[2]}, nil
}

func (c *RawCommand) String() string {
	return fmt.Sprintf("%s%d%s%s%d", string(ranks[c.srcX]), c.srcY, c.act, string(ranks[c.dstX]), c.dstY)
}

//ParsedCommand is a game command after being run through the engine
type ParsedCommand struct {
	srcX    int
	srcY    int
	srcRank string
	dstX    int
	dstY    int
	dstRank string
	act     string
}

func (c *ParsedCommand) String() string {
	if c.dstRank != "" {
		return fmt.Sprintf("%s %s%d %s %s %s%d", c.srcRank, string(ranks[c.srcX]), c.srcY, c.act, c.dstRank, string(ranks[c.dstX]), c.dstY)
	}
	return fmt.Sprintf("%s %s%d %s %s %s%d", c.srcRank, string(ranks[c.srcX]), c.srcY, c.act, "empty", string(ranks[c.dstX]), c.dstY)

}
