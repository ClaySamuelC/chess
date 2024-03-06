package main

import (
	"chess/game"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenW  = 502
	screenH  = 500
	tileSize = 58
	xOffset  = 19
	yOffset  = 19
)

func PrintSquares(squares []*game.Vector2) {
	fmt.Println("Possible moves: [")
	for _, square := range squares {
		fmt.Printf("  Can move to (%d, %d)\n", square.X, square.Y)
	}
	fmt.Println("]")
}

func getScreenSpace(tile_y int, tile_x int) (float64, float64) {
	y := tile_y*tileSize + yOffset - 3
	x := tile_x*tileSize + xOffset - 2

	return float64(y), float64(x)
}

func getTileSpace(screen_y int, screen_x int) (int, int) {
	y := math.Floor(float64((screen_y - yOffset) / tileSize))
	x := math.Floor(float64((screen_x - xOffset) / tileSize))

	return int(y), int(x)
}

var (
	backgroundImg *ebiten.Image
	pieceImageMap map[string]*ebiten.Image

	board         *game.Board
	selectedColor color.RGBA
	moveableColor color.RGBA

	moveableTiles []*game.Vector2
)

func getImage(imagePath string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func getPieceImage(piece *game.Piece) *ebiten.Image {
	var key string = piece.Team + piece.Rank

	img, ok := pieceImageMap[key]
	if ok == false {
		log.Fatalf("Image %s doesn't exist", key)
	}

	return img
}

func isInMoveableTiles(x int, y int) bool {
	for _, tile := range moveableTiles {
		if tile.X == x && tile.Y == y {
			return true
		}
	}

	return false
}

func init() {
	backgroundImg = getImage("assets/board.png")
	selectedColor = color.RGBA{0x00, 0x80, 0xFF, 0x80}
	moveableColor = color.RGBA{0xFF, 0x80, 0x00, 0x80}

	pieceImageMap = make(map[string]*ebiten.Image)
	pieceImageMap["WhitePawn"] = getImage("assets/WhitePawn.png")
	pieceImageMap["WhiteRook"] = getImage("assets/WhiteRook.png")
	pieceImageMap["WhiteKnight"] = getImage("assets/WhiteKnight.png")
	pieceImageMap["WhiteBishop"] = getImage("assets/WhiteBishop.png")
	pieceImageMap["WhiteQueen"] = getImage("assets/WhiteQueen.png")
	pieceImageMap["WhiteKing"] = getImage("assets/WhiteKing.png")
	pieceImageMap["BlackPawn"] = getImage("assets/BlackPawn.png")
	pieceImageMap["BlackRook"] = getImage("assets/BlackRook.png")
	pieceImageMap["BlackKnight"] = getImage("assets/BlackKnight.png")
	pieceImageMap["BlackBishop"] = getImage("assets/BlackBishop.png")
	pieceImageMap["BlackQueen"] = getImage("assets/BlackQueen.png")
	pieceImageMap["BlackKing"] = getImage("assets/BlackKing.png")

	board = game.NewBoard()
}

type Game struct{}

func (g *Game) Update() error {

	// check if mouse is clicked
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {

		// get mouse coords and convert to game space
		mouseX, mouseY := ebiten.CursorPosition()
		mouseTileY, mouseTileX := getTileSpace(mouseY, mouseX)

		// check if clicked square is within bounds
		if mouseTileX >= 0 && mouseTileX <= 7 && mouseTileY >= 0 && mouseTileY <= 7 {

			// Check if it is appropriate to make move
			if board.IsHighlighted && isInMoveableTiles(mouseTileX, mouseTileY) {

				// make the move
				board.Move(board.HighlightX, board.HighlightY, mouseTileX, mouseTileY)
				board.IsHighlighted = false

			} else {

				// Attempt to highlight clicked square
				success := board.HighlightSquare(mouseTileX, mouseTileY)
				if success {
					moveableTiles = board.GetPossibleMoves(board.Squares[mouseTileY][mouseTileX], &game.Vector2{mouseTileX, mouseTileY})

					PrintSquares(moveableTiles)
				}

			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Create Background
	screen.Fill(color.RGBA{0x31, 0x2e, 0x2b, 0xff})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, yOffset)
	screen.DrawImage(backgroundImg, op)

	// Highlight Selected Square
	if board.IsHighlighted {
		y, x := getScreenSpace(board.HighlightY, board.HighlightX)
		vector.DrawFilledRect(screen, float32(x), float32(y), tileSize, tileSize, selectedColor, false)
	}

	// Highlight Available Moves
	if moveableTiles != nil && board.IsHighlighted {
		for _, tile := range moveableTiles {
			y, x := getScreenSpace(tile.Y, tile.X)
			vector.DrawFilledRect(screen, float32(x), float32(y), tileSize, tileSize, moveableColor, false)
		}
	}

	// Draw Chess Pieces
	for y, row := range board.Squares {
		for x, piece := range row {
			if piece != nil {
				yDraw, xDraw := getScreenSpace(y, x)

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(xDraw, yDraw)
				screen.DrawImage(getPieceImage(piece), op)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenW, screenH
}

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Render an image")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
