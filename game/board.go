package game

import "fmt"

type Board struct {
	Squares       [8][8]*Piece
	PlayerTurn    string
	CurrentTurn   int
	IsHighlighted bool
	HighlightX    int
	HighlightY    int
	PossibleMoves []*Vector2
	EnPassantLoc  *Vector2
	LastEnPassant int
	PlayerInfo    map[string]*Info
}

type Info struct {
	CanCastle bool
	IsInCheck bool
}

type Vector2 struct {
	X int
	Y int
}

func (b *Board) CanEnPassant(pawnLoc *Vector2, dx int, dy int) bool {
	fmt.Printf("Current Turn: %v, Last En Passantable Move: %v\n", b.CurrentTurn, b.LastEnPassant)
	if b.LastEnPassant != b.CurrentTurn-1 {
		return false
	}
	if b.PlayerTurn == "White" {
		dy = -dy
	}
	return b.EnPassantLoc.X == pawnLoc.X+dx && b.EnPassantLoc.Y == pawnLoc.Y+dy
}

func (b *Board) HighlightSquare(x int, y int) bool {
	if b.IsHighlighted && b.HighlightX == x && b.HighlightY == y {
		return false
	}
	if b.Squares[y][x] != nil {
		if b.Squares[y][x].Team == b.PlayerTurn {
			fmt.Printf("Highlighting Square(%d, %d)\n", x, y)
			b.IsHighlighted = true
			b.HighlightX = x
			b.HighlightY = y

			return true
		}
	}

	b.IsHighlighted = false
	return false
}

func NewBoard() *Board {
	b := &Board{}
	b.SetDefault()
	return b
}

func (b *Board) Move(srcX int, srcY int, destX int, destY int) {
	p := b.Squares[srcY][srcX]

	b.Squares[destY][destX] = p
	b.Squares[srcY][srcX] = nil

	fmt.Printf("Moving from (%v, %v) to (%v, %v).\n", srcX, srcY, destX, destY)

	// special pawn cases
	if p.Rank == "Pawn" {
		if destY == 0 || destY == 7 { // Promotion
			fmt.Println("Promote!!!!")
			p.Rank = "Queen"
		} else if destY-srcY == 2 { // If pawn double moved, update en passant info
			b.EnPassantLoc = &Vector2{destX, destY - 1}
			b.LastEnPassant = b.CurrentTurn
			fmt.Printf("Making double move, En Passant loc = (%v, %v), En Passant/Current Turn (%v, %v).\n", b.EnPassantLoc.X, b.EnPassantLoc.Y, b.LastEnPassant, b.CurrentTurn)
		} else if destY-srcY == -2 {
			b.EnPassantLoc = &Vector2{destX, destY + 1}
			b.LastEnPassant = b.CurrentTurn
			fmt.Printf("Making double move, En Passant loc = (%v, %v), En Passant/Current Turn (%v, %v).\n", b.EnPassantLoc.X, b.EnPassantLoc.Y, b.LastEnPassant, b.CurrentTurn)
		}
		// Handle En Passant
		if b.LastEnPassant == b.CurrentTurn-1 && destX == b.EnPassantLoc.X && destY == b.EnPassantLoc.Y {
			// destroy pawn under move
			dy := -1
			if b.PlayerTurn == "White" {
				dy = 1
			}

			b.Squares[destY+dy][destX] = nil
		}
	}

	b.IsHighlighted = false
	b.CurrentTurn += 1
	if b.PlayerTurn == "White" {
		b.PlayerTurn = "Black"
	} else {
		b.PlayerTurn = "White"
	}
	fmt.Printf("Current turn: %v, Last Double Pawn Move: %v, En Passant Loc: (%v, %v)\n", b.CurrentTurn, b.LastEnPassant, b.EnPassantLoc.X, b.EnPassantLoc.Y)
}

func (b *Board) SetDefault() {
	for i := 0; i < 8; i++ {
		b.Squares[1][i] = NewPiece("Pawn", "Black")
		b.Squares[6][i] = NewPiece("Pawn", "White")
	}

	b.Squares[0][0] = NewPiece("Rook", "Black")
	b.Squares[0][1] = NewPiece("Knight", "Black")
	b.Squares[0][2] = NewPiece("Bishop", "Black")
	b.Squares[0][3] = NewPiece("Queen", "Black")
	b.Squares[0][4] = NewPiece("King", "Black")
	b.Squares[0][5] = NewPiece("Bishop", "Black")
	b.Squares[0][6] = NewPiece("Knight", "Black")
	b.Squares[0][7] = NewPiece("Rook", "Black")

	b.Squares[7][0] = NewPiece("Rook", "White")
	b.Squares[7][1] = NewPiece("Knight", "White")
	b.Squares[7][2] = NewPiece("Bishop", "White")
	b.Squares[7][3] = NewPiece("Queen", "White")
	b.Squares[7][4] = NewPiece("King", "White")
	b.Squares[7][5] = NewPiece("Bishop", "White")
	b.Squares[7][6] = NewPiece("Knight", "White")
	b.Squares[7][7] = NewPiece("Rook", "White")

	b.CurrentTurn = 1
	b.LastEnPassant = -1
	b.IsHighlighted = false
	b.PlayerTurn = "White"
}
