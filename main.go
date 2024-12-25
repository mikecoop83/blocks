package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mikecoop83/blocks/lib"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Game struct represents the game state.
type Game struct {
	board *lib.Board

	currentPiece *lib.Piece
}

// Update is called every tick (1/60 seconds by default) to update the game state.
func (g *Game) Update() error {
	if g.currentPiece == nil {
		randPieceIdx := rand.Intn(len(lib.AllPieces))
		piece := lib.AllPieces[randPieceIdx]
		g.currentPiece = &piece
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyN) {
		piece := lib.AllPieces[rand.Intn(len(lib.AllPieces))]
		g.currentPiece = &piece
	}
	return nil
}

var white = color.RGBA{R: 0xf0, G: 0xf0, B: 0xf0, A: 0xff}
var green = color.RGBA{R: 0x00, G: 0xcc, B: 0x66, A: 0xff}
var red = color.RGBA{R: 0xff, G: 0x66, B: 0x66, A: 0xff}
var blue = color.RGBA{R: 0x66, G: 0x99, B: 0xff, A: 0xff}
var orange = color.RGBA{R: 0xff, G: 0xa5, B: 0x00, A: 0xff}

const cellSize = 100

var cellStateToColor = map[lib.CellState]color.Color{
	lib.Empty:    white,
	lib.Pending:  green,
	lib.Invalid:  red,
	lib.FullLine: orange,
	lib.Occupied: blue,
}

// Draw is called every frame to render the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the background.
	vector.DrawFilledRect(
		screen,
		0, 0,
		float32(screen.Bounds().Max.X),
		float32(screen.Bounds().Max.Y),
		white,
		false,
	)

	grid := g.board.GetGrid()
	if g.currentPiece != nil {
		piece := *g.currentPiece
		mouseX, mouseY := ebiten.CursorPosition()
		pieceWidth := len(piece.Shape)
		pieceHeight := len(piece.Shape[0])
		adjustedX := mouseX - pieceWidth*cellSize/2
		adjustedY := mouseY - pieceHeight*cellSize/2
		cellX, cellY := adjustedX/cellSize, adjustedY/cellSize
		// Clamp the piece to the board.
		if cellX < 0 {
			cellX = 0
		}
		if cellY < 0 {
			cellY = 0
		}
		if cellX > lib.BoardSize-pieceWidth {
			cellX = lib.BoardSize - pieceWidth
		}
		if cellY > lib.BoardSize-pieceHeight {
			cellY = lib.BoardSize - pieceHeight
		}

		pending := true
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			pending = false
		}

		pendingGrid, valid := g.board.AddPiece(
			lib.PieceLocation{
				Piece: piece,
				Loc:   lib.Location{X: cellX, Y: cellY},
			},
			pending,
		)
		if valid || pending {
			grid = pendingGrid
		}
		if valid && !pending {
			g.currentPiece = nil
		}
	}
	for c := range grid {
		for r := range grid[c] {
			state := grid[c][r]
			cellColor := cellStateToColor[state]
			vector.DrawFilledRect(
				screen,
				float32(c*cellSize), float32(r*cellSize),
				float32(cellSize), float32(cellSize),
				cellColor,
				false,
			)
		}
	}

	// Draw gridlines
	for i := 0; i <= lib.BoardSize; i++ {
		vector.StrokeLine(
			screen,
			float32(i*cellSize), 0,
			float32(i*cellSize), float32(lib.BoardSize*cellSize),
			1,
			color.Black,
			false,
		)
		vector.StrokeLine(
			screen,
			0, float32(i*cellSize),
			float32(lib.BoardSize*cellSize), float32(i*cellSize),
			1,
			color.Black,
			false,
		)
	}
}

// Layout returns the logical screen dimensions. The game window will scale to fit this size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	minDim := min(outsideWidth, outsideHeight)
	return minDim, minDim
}

func main() {
	// Initialize the game object.
	board := lib.NewBoard()
	game := &Game{
		board: &board,
	}

	// Set the window title.
	ebiten.SetWindowTitle("Blocks")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(lib.BoardSize*cellSize, lib.BoardSize*cellSize)

	// Run the game.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
