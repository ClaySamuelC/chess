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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenW  = 502
	screenH  = 500
	tileSize = 58
	xOffset  = 19
	yOffset  = 19
)

func getScreenSpace(tile int) (int, int) {
	x := (tile%8)*tileSize + xOffset - 2
	y := (tile/8)*tileSize + yOffset - 5

	return x, y
}

func getTileSpace(screen_x int, screen_y int) int {
	y := math.Floor(float64((screen_y - yOffset) / tileSize))
	x := math.Floor(float64((screen_x - xOffset) / tileSize))

	return int(y*8 + x)
}

var (
	backgroundImg *ebiten.Image
	pieceImageMap map[string]*ebiten.Image

	selectedColor color.RGBA
	moveableColor color.RGBA

	chess          *game.Chess
	selectedSquare int
	possibleMoves  []int
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

	return pieceImageMap[key]
}

func isIn(x int, arr []int) bool {
	for _, num := range arr {
		if x == num {
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

	var err error
	chess, err = game.CreateGame("r3k2r/1n6/8/8/8/8/6N1/R3K2R w KQkq - 0 1")
	if err != nil {
		log.Fatal(err)
	}

	selectedSquare = -1
}

type Game struct{}

func (g *Game) Update() error {
	// get mouse coords and convert to game space
	mouseX, mouseY := ebiten.CursorPosition()
	mouseTile := getTileSpace(mouseX, mouseY)

	// check if mouse is clicked
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		fmt.Printf("%v (%c%v)\n", mouseTile, 'a'+mouseTile%8, 8-mouseTile/8)

		// check if clicked square is within bounds
		if mouseTile >= 0 && mouseTile < 64 {

			// Check if it is appropriate to make move
			if selectedSquare != -1 && isIn(mouseTile, possibleMoves) {

				// make the move
				chess.Move(selectedSquare, mouseTile)
				selectedSquare = -1

				// Attempt to highlight clicked square
			} else if chess.Board[mouseTile] != nil && chess.Board[mouseTile].Team == chess.Turn {
				fmt.Printf("Selecting square %v (%v %v)\n", mouseTile, chess.Board[mouseTile].Team, chess.Board[mouseTile].Rank)
				selectedSquare = mouseTile
				possibleMoves = chess.GetPossibleMoves(chess.Board[mouseTile], mouseTile)
			} else {
				selectedSquare = -1
			}
		}
	} else if selectedSquare != -1 && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if isIn(mouseTile, possibleMoves) {

			// make the move
			chess.Move(selectedSquare, mouseTile)
			selectedSquare = -1

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

	// Highlight the selected square blue
	if selectedSquare != -1 {
		x, y := getScreenSpace(selectedSquare)
		vector.DrawFilledRect(screen, float32(x), float32(y), tileSize, tileSize, selectedColor, false)
	}

	// Highlight Available Moves
	if possibleMoves != nil && selectedSquare != -1 {
		for _, tile := range possibleMoves {
			x, y := getScreenSpace(tile)
			vector.DrawFilledRect(screen, float32(x), float32(y), tileSize, tileSize, moveableColor, false)
		}
	}

	// Draw Chess Pieces
	for i, piece := range chess.Board {
		if piece != nil {
			xDraw, yDraw := getScreenSpace(i)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(xDraw), float64(yDraw))
			screen.DrawImage(getPieceImage(piece), op)
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
