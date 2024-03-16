package game

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	UpLeft    = &Vector2{-1, -1}
	UpRight   = &Vector2{-1, 1}
	DownLeft  = &Vector2{1, -1}
	DownRight = &Vector2{1, 1}
	Up        = &Vector2{0, -1}
	Right     = &Vector2{1, 0}
	Down      = &Vector2{0, 1}
	Left      = &Vector2{-1, 0}

	Diagonals      = &[]*Vector2{UpLeft, UpRight, DownRight, DownLeft}
	Straights      = &[]*Vector2{Left, Up, Right, Down}
	Adjacents      = &[]*Vector2{Left, UpLeft, Up, UpRight, Right, DownRight, Down, DownLeft}
	WhitePawnMoves = &[]*Vector2{UpLeft, UpRight}
	BlackPawnMoves = &[]*Vector2{DownLeft, DownRight}
	KnightMoves    = &[]*Vector2{{-2, -1}, {-1, -2}, {1, -2}, {2, -1}, {2, 1}, {1, 2}, {-1, 2}, {-2, 1}}
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
	KingPos            int
}

type Vector2 struct {
	X int
	Y int
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

	// set castle info
	c.PlayerInfo = make(map[string]*Info)

	c.PlayerInfo["White"] = &Info{false, false, 0}
	c.PlayerInfo["Black"] = &Info{false, false, 0}

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
			if char == 'K' {
				c.PlayerInfo["White"].KingPos = y*8 + x
			} else if char == 'k' {
				c.PlayerInfo["Black"].KingPos = y*8 + x
			}
			piece, ok := pieceMap[byte(char)]
			if !ok {
				return nil, fmt.Errorf("invalid FEN char: %c", char)
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

	for _, char := range info[2] {
		switch char {
		case 'K':
			fmt.Println("White can castle king side.")
			c.PlayerInfo["White"].IsKingCastleValid = true
		case 'Q':
			fmt.Println("White can castle queen side.")
			c.PlayerInfo["White"].IsQueenCastleValid = true
		case 'k':
			fmt.Println("Black can castle king side.")
			c.PlayerInfo["Black"].IsKingCastleValid = true
		case 'q':
			fmt.Println("Black can castle queen side.")
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

	fmt.Printf("White King Pos: %v\n", c.PlayerInfo["White"].KingPos)
	fmt.Printf("Black King Pos: %v\n", c.PlayerInfo["Black"].KingPos)

	return c, nil
}

func (c *Chess) Move(src int, dest int) {
	fmt.Printf("King Pos: %v\n", c.PlayerInfo[c.Turn].KingPos)
	p := c.Board[src]

	if p.Rank == "King" {
		c.PlayerInfo[c.Turn].IsKingCastleValid = false
		c.PlayerInfo[c.Turn].IsQueenCastleValid = false

		if dest-src == -2 {
			c.Board[dest+1] = c.Board[dest-2]
			c.Board[dest-2] = nil
		}
		if dest-src == 2 {
			c.Board[dest-1] = c.Board[dest+1]
			c.Board[dest+1] = nil
		}

		c.PlayerInfo[c.Turn].KingPos = dest
	} else if p.Rank == "Rook" {
		if src%8 == 0 {
			c.PlayerInfo[c.Turn].IsQueenCastleValid = false
		} else if src%8 == 7 {
			c.PlayerInfo[c.Turn].IsKingCastleValid = false
		}
	}

	if c.Board[dest] == nil && p.Rank != "Pawn" {
		c.HalfMoveClock += 1
	}
	if p.Rank == "Pawn" {
		if dest == c.EnPassantLoc {
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

	if c.Turn == "White" {
		c.Turn = "Black"
	} else {
		c.Turn = "White"
		c.FullMoveClock += 1
	}
}

func (c *Chess) GetAllMoves(team string) *map[int]*[]int {
	moveMap := map[int]*[]int{}

	for i, p := range c.Board {
		if p.Team == team {
			moveMap[i] = c.GetPossibleMoves(p, i)
		}
	}

	return &moveMap
}
