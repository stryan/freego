package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var cmdRegxp = regexp.MustCompile("([a-zA-Z])([0-9])(x|-)([a-zA-Z])([0-9])")
var ranks = "ABCDEFHIJKLMNOPQRSTUVWXYZ"

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
	rawRes := cmdRegxp.FindAllStringSubmatch(cmd, -1)
	if rawRes == nil {
		return nil, errors.New("error creating command from string")
	}
	res := rawRes[0]
	if len(res) != 6 {
		return nil, fmt.Errorf("expected more fields from command string 5!=%v, %v", len(res), res)
	}
	sx := strings.Index(ranks, strings.ToUpper(res[1]))
	if sx == -1 {
		return nil, fmt.Errorf("bad rank value: %v", res[1])
	}
	dx := strings.Index(ranks, strings.ToUpper(res[4]))
	if dx == -1 {
		return nil, fmt.Errorf("bad rank value: %v", strings.ToUpper(res[4]))
	}
	sy, err := strconv.Atoi(res[2])
	if err != nil {
		return nil, err
	}
	dy, err := strconv.Atoi(res[5])
	if err != nil {
		return nil, err
	}
	return &RawCommand{sx, sy, dx, dy, res[3]}, nil
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
