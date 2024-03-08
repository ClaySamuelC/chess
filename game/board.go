package game

import (
	"fmt"
	"strconv"
	"strings"
)

type Chess struct {
	Board         [64]*Piece
	Turn          string
	EnPassantLoc  int
	HalfMoveClock int
	FullMoveClock int
	PlayerInfo    map[string]*Info
}

type Info struct {
	IsKingCastleValid  bool
	IsQueenCastleValid bool
}

func CreateGame(fen string) (*Chess, error) {
	// default game: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0h"
	pieceMap := map[byte]Piece{
		'P': {"White", "Pawn"},
		'N': {"White", "Knight"},
		'B': {"White", "Bishop"},
		'R': {"White", "Rook"},
		'Q': {"White", "Queen"},
		'K': {"White", "King"},
		'p': {"Black", "Pawn"},
		'n': {"Black", "Knight"},
		'b': {"Black", "Bishop"},
		'r': {"Black", "Rook"},
		'q': {"Black", "Queen"},
		'k': {"Black", "King"},
	}

	info := strings.Fields(fen)

	c := &Chess{}
	x := 0
	y := 7

	// populate board
	for _, char := range info[0] {
		switch char {
		case '/':
			y--
			x = 0
		case '1', '2', '3', '4', '5', '6', '7', '8':
			x += int(char - '0')
		default:
			piece, ok := pieceMap[byte(char)]
			if !ok {
				return nil, fmt.Errorf("Invalid FEN char: %c", char)
			}
			c.Board[y*8+x] = &piece
			x++
		}
	}

	// set current turn
	if info[1] == "w" {
		c.Turn = "White"
	} else {
		c.Turn = "Black"
	}

	// set castle info
	c.PlayerInfo = make(map[string]*Info)

	c.PlayerInfo["White"] = &Info{false, false}
	c.PlayerInfo["Black"] = &Info{false, false}
	for _, char := range info[2] {
		switch char {
		case 'K':
			c.PlayerInfo["White"].IsKingCastleValid = true
		case 'Q':
			c.PlayerInfo["White"].IsQueenCastleValid = true
		case 'k':
			c.PlayerInfo["Black"].IsKingCastleValid = true
		case 'q':
			c.PlayerInfo["Black"].IsQueenCastleValid = true
		}
	}

	// set en passant targets
	if info[3] == "-" {
		c.EnPassantLoc = -1
	} else {
		c.EnPassantLoc = (int(info[3][1])-49)*8 + int(info[3][0]) - 97
	}

	// set half move clock
	c.HalfMoveClock, _ = strconv.Atoi(info[4])

	// set full move clock
	c.FullMoveClock, _ = strconv.Atoi(info[5])

	return c, nil
}

func (c *Chess) Move(src int, dest int) {
	p := c.Board[src]

	c.Board[dest] = p
	c.Board[src] = nil

	fmt.Printf("Moving %v from (%v, %v) to (%v, %v).\n", p.Rank, src/8, src%8, dest/8, dest%8)

	// special king cases
	if p.Rank == "King" {
		c.PlayerInfo[c.Turn].IsKingCastleValid = true

		// if castling, move the rook as well
		if dest-src == -2 {
			c.Board[dest-1] = c.Board[dest]
			c.Board[dest-4] = nil
		} else if dest-src == 2 {
			c.Board[dest+1] = c.Board[dest-4]
			c.Board[dest+3] = nil
		}
	}

	// special pawn cases
	if p.Rank == "Pawn" {
		if dest/8 == 0 || dest/8 == 7 { // Promotion
			fmt.Println("Promote!!!!")
			p.Rank = "Queen"
		} else if dest-src == 2 { // If pawn double moved, update en passant info
			c.EnPassantLoc = dest - 8
			fmt.Printf("Making double move, En Passant loc = (%v, %v), En Passant/Current Turn (%v, %v).\n", c.EnPassantLoc.X, c.EnPassantLoc.Y, c.LastEnPassant, c.CurrentTurn)
		} else if dest-src == -2 {
			c.EnPassantLoc = dest + 8
			fmt.Printf("Making double move, En Passant loc = (%v, %v), En Passant/Current Turn (%v, %v).\n", c.EnPassantLoc.X, c.EnPassantLoc.Y, c.LastEnPassant, c.CurrentTurn)
		}
	}

	if c.Turn == "White" {
		c.Turn = "Black"
	} else {
		c.Turn = "White"
	}
}
