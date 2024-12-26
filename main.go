package main

import (
	"image/color"
	"log"
	"math/rand"

	"golang.org/x/image/font"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mikecoop83/blocks/lib"
	"github.com/mikecoop83/blocks/resources"
)

// Game struct represents the game state.
type Game struct {
	board *lib.Board

	pieceOptions   [3]*lib.Piece
	chosenPieceIdx int

	score int64
}

func (g *Game) chosenPiece() *lib.Piece {
	return g.pieceOptions[g.chosenPieceIdx]
}

func getRandomRotatedPiece() *lib.Piece {
	randPieceIdx := rand.Intn(len(lib.AllPieces))
	piece := lib.AllPieces[randPieceIdx]
	rotateTimes := rand.Intn(4)
	for i := 0; i < rotateTimes; i++ {
		piece = piece.Rotate()
	}
	return &piece
}

// Update is called every tick (1/60 seconds by default) to update the game state.
func (g *Game) Update() error {
	// If the chosen piece is nil, choose the first available piece.
	if g.chosenPiece() == nil {
		for i := 0; i < len(g.pieceOptions); i++ {
			if g.pieceOptions[i] != nil {
				g.chosenPieceIdx = i
				break
			}
		}
	}
	// If none left, get 3 new pieces and set first piece to be chosen.
	if g.chosenPiece() == nil {
		for i := 0; i < 3; i++ {
			g.pieceOptions[i] = getRandomRotatedPiece()
		}
		g.chosenPieceIdx = 0
	}

	// Cheats...
	if inpututil.IsKeyJustReleased(ebiten.KeyS) {
		for i := 0; i < len(g.pieceOptions); i++ {
			g.pieceOptions[i] = getRandomRotatedPiece()
		}
	}
	return nil
}

var white = color.RGBA{R: 0xf0, G: 0xf0, B: 0xf0, A: 0xff}
var green = color.RGBA{R: 0x00, G: 0xcc, B: 0x66, A: 0xff}
var red = color.RGBA{R: 0xff, G: 0x66, B: 0x66, A: 0xff}
var blue = color.RGBA{R: 0x66, G: 0x99, B: 0xff, A: 0xff}
var orange = color.RGBA{R: 0xff, G: 0xa5, B: 0x00, A: 0xff}
var gray = color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff}
var paleYellow = color.RGBA{R: 0xff, G: 0xff, B: 0xcc, A: 0xff}

const (
	cellSize         = 100
	topAreaHeight    = 100
	boardWidth       = lib.BoardSize * cellSize
	boardHeight      = lib.BoardSize * cellSize
	bottomAreaHeight = lib.BoardSize * cellSize * 0.5
)

var cellStateToColor = map[lib.CellState]color.Color{
	lib.Empty:    white,
	lib.Pending:  green,
	lib.Invalid:  red,
	lib.FullLine: orange,
	lib.Occupied: blue,
	lib.Unchosen: gray,
	lib.Hovering: paleYellow,
}

var commaFormatter = message.NewPrinter(language.English)

// Draw is called every frame to render the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	mouseX, mouseY := ebiten.CursorPosition()
	// Draw the background.
	vector.DrawFilledRect(
		screen,
		0, 0,
		float32(screen.Bounds().Max.X),
		float32(screen.Bounds().Max.Y),
		white,
		false,
	)
	const topAreaOffset = 0
	// Draw the top area with the score.
	msg := commaFormatter.Sprintf("Score: %d", g.score)
	bounds, _ := font.BoundString(resources.FontFace, msg)

	// Calculate text width and height
	textWidth := (bounds.Max.X - bounds.Min.X) >> 6
	textHeight := (bounds.Max.Y - bounds.Min.Y) >> 6

	// Calculate the center position
	textX := (boardWidth - textWidth) / 2
	textY := (topAreaHeight - textHeight) / 2

	// Adjust for baseline alignment
	textY += textHeight
	text.Draw(
		screen,
		msg,
		resources.FontFace,
		int(textX), int(topAreaOffset+textY),
		color.Black,
	)

	const boardYOffset = topAreaHeight
	mouseOnBoard := mouseX >= 0 && mouseX < boardWidth && mouseY >= boardYOffset && mouseY < boardYOffset+boardHeight
	// Draw gridlines
	var gridColor = color.Gray16{Y: 0xBBBB}
	for i := 0; i <= lib.BoardSize; i++ {
		// Horizontal line
		vector.StrokeLine(
			screen,
			float32(i*cellSize), float32(boardYOffset),
			float32(i*cellSize), float32(boardYOffset+boardHeight),
			1,
			gridColor,
			false,
		)
		// Vertical line
		vector.StrokeLine(
			screen,
			0, float32(boardYOffset+i*cellSize),
			boardWidth, float32(boardYOffset+i*cellSize),
			1,
			gridColor,
			false,
		)
	}

	grid := g.board.GetGrid()
	if mouseOnBoard && g.chosenPiece() != nil {
		piece := *g.chosenPiece()
		cellC := mouseX / cellSize
		cellR := (mouseY - boardYOffset) / cellSize

		// Clamp the piece to the board if the mouse is on the board
		if cellC < 0 {
			cellC = 0
		}
		if cellR < 0 {
			cellR = 0
		}
		if cellC > lib.BoardSize-piece.Width() {
			cellC = lib.BoardSize - piece.Width()
		}
		if cellR > lib.BoardSize-piece.Height() {
			cellR = lib.BoardSize - piece.Height()
		}
		pending := !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
		pendingGrid, clearedLines, valid := g.board.AddPiece(
			lib.PieceLocation{
				Piece: piece,
				Loc:   lib.Location{C: cellC, R: cellR},
			},
			pending,
		)
		if valid || pending {
			grid = pendingGrid
		}
		if valid && !pending {
			numPoints := g.chosenPiece().NumBlocks()
			numPoints += clearedLines * 10
			g.score += int64(numPoints)
			g.pieceOptions[g.chosenPieceIdx] = nil
		}
	}
	for r := range grid {
		for c := range grid[r] {
			state := grid[r][c]
			if state == lib.Empty {
				continue
			}
			cellColor := cellStateToColor[state]
			vector.DrawFilledRect(
				screen,
				float32(c*cellSize), float32(boardYOffset+r*cellSize),
				cellSize, cellSize,
				cellColor,
				false,
			)
			// Draw rectangle around each filled cell.
			vector.StrokeRect(
				screen,
				float32(c*cellSize), float32(boardYOffset+r*cellSize),
				cellSize, cellSize,
				1,
				color.Black,
				false,
			)
		}
	}

	// Draw the bottom area with the piece options
	const bottomAreaOffset = topAreaHeight + boardHeight
	const pieceOptionCellSize = cellSize * 0.5
	pieceOptionWidth := boardWidth / 3
	for p, piece := range g.pieceOptions {
		if piece == nil {
			continue
		}
		pieceOptionColor := cellStateToColor[lib.Unchosen]
		if p == g.chosenPieceIdx {
			pieceOptionColor = cellStateToColor[lib.Pending]
		}
		yOffset := (bottomAreaHeight - piece.Height()*pieceOptionCellSize) / 2
		xOffset := (pieceOptionWidth - piece.Width()*pieceOptionCellSize) / 2
		// If the mouse is hovering over an unselected piece, change the color.  Select it if it was just clicked.
		pieceX := xOffset + p*pieceOptionWidth
		pieceY := bottomAreaOffset + yOffset
		if g.chosenPieceIdx != p &&
			mouseX >= pieceX && mouseX < pieceX+piece.Width()*pieceOptionCellSize &&
			mouseY >= pieceY && mouseY < pieceY+piece.Height()*pieceOptionCellSize {
			pieceOptionColor = cellStateToColor[lib.Hovering]
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				g.chosenPieceIdx = p
			}
		}
		for r := range piece.Shape {
			for c := range piece.Shape[r] {
				if !piece.Shape[r][c] {
					continue
				}
				vector.DrawFilledRect(
					screen,
					float32(pieceX+c*pieceOptionCellSize),
					float32(pieceY+r*pieceOptionCellSize),
					float32(pieceOptionCellSize),
					float32(pieceOptionCellSize),
					pieceOptionColor,
					false,
				)
				// Draw rectangle around each filled cell.
				vector.StrokeRect(
					screen,
					float32(pieceX+c*pieceOptionCellSize),
					float32(pieceY+r*pieceOptionCellSize),
					float32(pieceOptionCellSize),
					float32(pieceOptionCellSize),
					1,
					color.Black,
					false,
				)
			}
		}
	}
}

// Layout returns the logical screen dimensions. The game window will scale to fit this size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return boardWidth, topAreaHeight + boardHeight + bottomAreaHeight
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
	ebiten.SetWindowSize(boardWidth, boardHeight+bottomAreaHeight)

	// Run the game.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
