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

func CreateDefaultGame() (*Chess, error) {
	game, err := CreateGame("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0h")

	return game, err
}

func CreateGame(fen string) (*Chess, error) {
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
	y := 0

	// populate board
	for _, char := range info[0] {
		switch char {
		case '/':
			y++
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
		c.EnPassantLoc = -999
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

	if c.Board[dest] == nil && p.Rank != "Pawn" {
		c.HalfMoveClock += 1
	}
	if p.Rank == "Pawn" {
		if dest == c.EnPassantLoc {
			fmt.Printf("En Passant at %v\n", dest)
			if c.Turn == "White" {
				c.Board[dest+8] = nil
			} else {
				c.Board[dest-8] = nil
			}
		}
		if dest-src == 16 {
			c.EnPassantLoc = dest - 8
		} else if dest-src == -16 {
			c.EnPassantLoc = dest + 8
		} else {
			c.EnPassantLoc = -999
		}
	} else {
		c.EnPassantLoc = -999
	}

	c.Board[dest] = p
	c.Board[src] = nil

	fmt.Printf("Moving %v from (%v, %v) to (%v, %v).\n", p.Rank, src%8, src/8, dest%8, dest/8)

	if c.Turn == "White" {
		c.Turn = "Black"
	} else {
		c.Turn = "White"
		c.FullMoveClock += 1
	}
}
