package game

import (
	"image/color"
	"log/slog"
	"math/rand"
	"time"

	"github.com/mikecoop83/blocks/persist"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mikecoop83/blocks/lib"
	"github.com/mikecoop83/blocks/resources"
)

var (
	offWhite    = color.RGBA{R: 0xf0, G: 0xf0, B: 0xf0, A: 0xff}
	green       = color.RGBA{R: 0x00, G: 0xcc, B: 0x66, A: 0xff}
	red         = color.RGBA{R: 0xff, G: 0x66, B: 0x66, A: 0xff}
	blue        = color.RGBA{R: 0x66, G: 0x99, B: 0xff, A: 0xff}
	orange      = color.RGBA{R: 0xff, G: 0xa5, B: 0x00, A: 0xff}
	gray        = color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff}
	darkGray    = color.RGBA{R: 0x40, G: 0x40, B: 0x40, A: 0xff}
	paleYellow  = color.RGBA{R: 0xff, G: 0xff, B: 0xcc, A: 0xff}
	reddishGray = color.RGBA{R: 0x99, G: 0x66, B: 0x66, A: 0xff}
)

const (
	splashDuration   = 1.25 * time.Second
	cellSize         = 100
	topAreaHeight    = 100
	numPieceOptions  = 3
	boardWidth       = lib.BoardSize * cellSize
	boardHeight      = lib.BoardSize * cellSize
	bottomAreaHeight = lib.BoardSize * cellSize * 0.5
	WindowWidth      = boardWidth
	WindowHeight     = topAreaHeight + boardHeight + bottomAreaHeight
)

var displayModeToCellColor = map[DisplayMode]map[lib.CellState]color.Color{
	DisplayModeNormal: cellStateToColor,
	DisplayModeDark:   darkCellStateToColor,
}

var displayModeToBackgroundColor = map[DisplayMode]color.Color{
	DisplayModeNormal: offWhite,
	DisplayModeDark:   darkGray,
}

var displayModeToForegroundColor = map[DisplayMode]color.Color{
	DisplayModeNormal: color.Black,
	DisplayModeDark:   offWhite,
}

var cellStateToColor = map[lib.CellState]color.Color{
	lib.Empty:    offWhite,
	lib.Pending:  green,
	lib.Invalid:  red,
	lib.FullLine: orange,
	lib.Occupied: blue,
	lib.Unchosen: gray,
	lib.Hovering: paleYellow,
}

var darkCellStateToColor = map[lib.CellState]color.Color{
	lib.Empty:    darkGray,
	lib.Pending:  green,
	lib.Invalid:  red,
	lib.FullLine: orange,
	lib.Occupied: blue,
	lib.Unchosen: gray,
	lib.Hovering: paleYellow,
}
var commaFormatter = message.NewPrinter(language.English)

type DisplayMode int

const (
	DisplayModeNormal DisplayMode = iota
	DisplayModeDark
)

var displayModeToName = map[DisplayMode]string{
	DisplayModeNormal: "normal",
	DisplayModeDark:   "dark",
}

var nameToDisplayMode = map[string]DisplayMode{
	"normal": DisplayModeNormal,
	"dark":   DisplayModeDark,
}

// Game struct represents the game state.
type Game struct {
	board *lib.Board

	gameID       uint64
	randSource   rand.Source
	updateGameID func(gameID uint64)

	pieceOptions [numPieceOptions]*lib.Piece

	clearedRows [lib.BoardSize]*animatedEntity
	clearedCols [lib.BoardSize]*animatedEntity

	touchEnabled       bool
	pressX, pressY     int
	dragX, dragY       int
	releaseX, releaseY int
	chosenPieceIdx     int

	score       int64
	highScore   int64
	gameOver    bool
	cheating    bool
	cheated     bool
	displayMode DisplayMode

	splashStart time.Time
}

func (g *Game) Reset(gameID uint64) {
	newBoard := lib.NewBoard()
	g.board = &newBoard
	g.pieceOptions = [numPieceOptions]*lib.Piece{}
	g.chosenPieceIdx = -1
	g.score = 0
	g.gameOver = false
	g.highScore = maybeGetHighScore()
	g.cheated = false
	displayModeText, err := persist.Load("displaymode")
	if err != nil {
		slog.Error("error loading display mode: %v", err)
	}
	g.displayMode = nameToDisplayMode[displayModeText]
	g.gameID = gameID
	g.randSource = rand.NewSource(int64(g.gameID))
	g.updateGameID(g.gameID)
}

func New(gameID uint64, updateGameID func(gameID uint64)) ebiten.Game {
	game := &Game{
		updateGameID: updateGameID,
	}
	game.Reset(gameID)
	return game
}

func (g *Game) chosenPiece() *lib.Piece {
	if g.cheating {
		return &lib.AllPieces[0]
	}
	if g.chosenPieceIdx < 0 || g.chosenPieceIdx >= len(g.pieceOptions) {
		return nil
	}
	return g.pieceOptions[g.chosenPieceIdx]
}

// Update is called every tick (1/60 seconds by default) to tick the game state.
func (g *Game) Update() error {
	switchMode := func() {
		g.displayMode = (g.displayMode + 1) % 2
		err := persist.Store("displaymode", displayModeToName[g.displayMode])
		if err != nil {
			slog.Error("error storing display mode: %v", err)
		}
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyM) {
		switchMode()
	}
	if g.splashStart.IsZero() {
		g.splashStart = time.Now()
	}

	if !g.cheated && g.score > g.highScore {
		g.highScore = g.score
		maybeUpdateHighScore(g.highScore)
	}
	g.cheating = ebiten.IsKeyPressed(ebiten.KeyMeta) && ebiten.IsKeyPressed(ebiten.KeyShift)
	var pressedTouchIDs, dragTouchIDs, releasedTouchIDs []ebiten.TouchID
	pressedTouchIDs = inpututil.AppendJustPressedTouchIDs(pressedTouchIDs)
	dragTouchIDs = ebiten.AppendTouchIDs(dragTouchIDs)
	releasedTouchIDs = inpututil.AppendJustReleasedTouchIDs(releasedTouchIDs)
	if len(pressedTouchIDs) > 0 || len(dragTouchIDs) > 0 || len(releasedTouchIDs) > 0 {
		g.touchEnabled = true
	}
	if g.touchEnabled {
		// Triple-touch screen switches mode
		if len(pressedTouchIDs) > 2 {
			switchMode()
		}
		for _, id := range pressedTouchIDs {
			g.pressX, g.pressY = ebiten.TouchPosition(id)
		}
		for _, id := range dragTouchIDs {
			dragX, dragY := ebiten.TouchPosition(id)
			// Offset touch dragY to be above your finger by a bit more than the height of the piece to see where you're
			// dragging it
			var dragYOffset int
			chosenPiece := g.chosenPiece()
			if chosenPiece != nil {
				dragYOffset = (chosenPiece.Height() + 1) * cellSize
			}
			g.dragX, g.dragY = dragX, dragY-dragYOffset
		}
		if len(releasedTouchIDs) > 0 {
			g.releaseX, g.releaseY = g.dragX, g.dragY
			g.dragX, g.dragY = -1, -1
			g.pressX, g.pressY = -1, -1
		}
	} else {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.pressX, g.pressY = ebiten.CursorPosition()
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			g.dragX, g.dragY = ebiten.CursorPosition()
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			g.releaseX, g.releaseY = ebiten.CursorPosition()
			g.dragX, g.dragY = -1, -1
			g.pressX, g.pressY = -1, -1
		}
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		g.Reset(rand.Uint64())
	}

	// If no pieces to choose, get numPieceOptions new pieces and set first piece to be chosen.
	if g.pieceOptions[0] == nil && g.pieceOptions[1] == nil && g.pieceOptions[2] == nil {
		for i := 0; i < numPieceOptions; i++ {
			randomPiece := lib.RandomRotatedPiece(g.randSource)
			g.pieceOptions[i] = &randomPiece
		}
	}
	// Check if the game is over.
	var hasValidMove bool
outer:
	for _, piece := range g.pieceOptions {
		if piece == nil {
			continue
		}
		for r := range lib.BoardSize {
			for c := range lib.BoardSize {
				pieceLoc := lib.PieceLocation{
					Piece: *piece,
					Loc:   lib.Location{C: c, R: r},
				}
				if g.board.ValidatePiece(pieceLoc, false) {
					hasValidMove = true
					break outer
				}
			}
		}
	}
	if !hasValidMove {
		g.gameOver = true
	}

	// Update the animations for cleared rows and columns.
	for _, rowsAndColumns := range [2]*[lib.BoardSize]*animatedEntity{&g.clearedRows, &g.clearedCols} {
		for i, entity := range rowsAndColumns {
			if entity == nil {
				continue
			}
			if entity.tick() {
				rowsAndColumns[i] = nil
			}
		}
	}
	return nil
}

// Draw is called every frame to render the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBackground(screen)

	if time.Since(g.splashStart) < splashDuration {
		g.drawSplash(screen)
		return
	}

	g.drawGame(screen)
}

func (g *Game) drawGame(screen *ebiten.Image) {
	defer func() {
		if g.releaseX >= 0 && g.releaseY >= 0 {
			g.chosenPieceIdx = -1
		}
		g.releaseX, g.releaseY = -1, -1
	}()

	g.drawHeader(screen)

	g.drawBoard(screen)

	g.drawOverlay(screen)

	g.drawPieceOptions(screen)
}

func (g *Game) drawPieceOptions(screen *ebiten.Image) {
	// Draw the bottom area with the piece options
	const bottomAreaOffset = topAreaHeight + boardHeight
	const pieceOptionCellSize = cellSize * 0.5
	pieceOptionWidth := boardWidth / numPieceOptions
	stateToColor := displayModeToCellColor[g.displayMode]
	for p, piece := range g.pieceOptions {
		if piece == nil {
			continue
		}
		pieceOptionColor := stateToColor[lib.Unchosen]
		if p == g.chosenPieceIdx && g.releaseX >= 0 && g.releaseY >= 0 {
			pieceOptionColor = stateToColor[lib.Pending]
		}
		yOffset := (bottomAreaHeight - piece.Height()*pieceOptionCellSize) / 2
		xOffset := (pieceOptionWidth - piece.Width()*pieceOptionCellSize) / 2
		// If the mouse is hovering over an unselected piece, change the color.  Select it if it was just clicked.
		pieceX := xOffset + p*pieceOptionWidth
		pieceY := bottomAreaOffset + yOffset
		// Break the bottom area into large sections so we don't require touching the piece itself to drag it on to the
		// board
		pieceAreaX := p * pieceOptionWidth
		pieceAreaY := bottomAreaOffset
		if g.pressX >= pieceAreaX && g.pressX < pieceAreaX+pieceOptionWidth &&
			g.pressY >= pieceAreaY && g.pressY < pieceAreaY+bottomAreaHeight {
			pieceOptionColor = stateToColor[lib.Hovering]
			g.chosenPieceIdx = p
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
					displayModeToForegroundColor[g.displayMode],
					false,
				)
			}
		}
	}
}

func (g *Game) drawOverlay(screen *ebiten.Image) {
	// Draw gridlines
	for i := 0; i <= lib.BoardSize; i++ {
		// Horizontal line
		vector.StrokeLine(
			screen,
			float32(i*cellSize), float32(topAreaHeight),
			float32(i*cellSize), float32(topAreaHeight+boardHeight),
			1,
			displayModeToForegroundColor[g.displayMode],
			false,
		)
		// Vertical line
		vector.StrokeLine(
			screen,
			0, float32(topAreaHeight+i*cellSize),
			boardWidth, float32(topAreaHeight+i*cellSize),
			1,
			displayModeToForegroundColor[g.displayMode],
			false,
		)
	}
	// If game over, gray out the board with transparency
	if g.gameOver {
		vector.DrawFilledRect(
			screen,
			0, float32(topAreaHeight),
			boardWidth, boardHeight,
			color.RGBA{R: 0, G: 0, B: 0, A: 0x80},
			false,
		)
	}
}

func (g *Game) drawBoard(screen *ebiten.Image) {
	grid := g.board.GetGrid()
	stateToColor := displayModeToCellColor[g.displayMode]

	// Either drag or click is the current mouse position.
	mouseX, mouseY := g.dragX, g.dragY
	if mouseX < 0 {
		mouseX, mouseY = g.releaseX, g.releaseY
	}
	onBoard := mouseX >= 0 &&
		mouseX < boardWidth &&
		mouseY >= topAreaHeight &&
		mouseY < topAreaHeight+boardHeight

	if onBoard && g.chosenPiece() != nil {
		piece := *g.chosenPiece()
		cellC := mouseX / cellSize
		cellR := (mouseY - topAreaHeight) / cellSize

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
		pending := true
		if g.releaseX >= 0 && g.releaseY >= 0 {
			pending = false
		}
		pendingGrid, clearedRows, clearedCols, valid := g.board.AddPiece(
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
			numClearedLines := len(clearedRows) + len(clearedCols)
			numPoints += numClearedLines * 10
			if grid.Empty() {
				numPoints += 300
			}
			for _, r := range clearedRows {
				g.clearedRows[r] = &animatedEntity{
					currentColor:  stateToColor[lib.FullLine],
					targetColor:   stateToColor[lib.Empty],
					animationTime: 1 * time.Second,
				}
			}
			for _, c := range clearedCols {
				g.clearedCols[c] = &animatedEntity{
					currentColor:  stateToColor[lib.FullLine],
					targetColor:   stateToColor[lib.Empty],
					animationTime: 1 * time.Second,
				}
			}
			g.score += int64(numPoints)
			if !g.cheating {
				g.pieceOptions[g.chosenPieceIdx] = nil
			} else {
				g.cheated = true
			}
		}
	}
	// Draw the cells
	for r := range grid {
		for c := range grid[r] {
			state := grid[r][c]
			displayColors := displayModeToCellColor[g.displayMode]
			displayColor := displayColors[state]
			if state == lib.Empty {
				if g.clearedCols[c] != nil {
					displayColor = g.clearedCols[c].currentColor
				}
				if g.clearedRows[r] != nil {
					displayColor = g.clearedRows[r].currentColor
				}
			}
			vector.DrawFilledRect(
				screen,
				float32(c*cellSize), float32(topAreaHeight+r*cellSize),
				cellSize, cellSize,
				displayColor,
				false,
			)
		}
	}
}

func (g *Game) drawHeader(screen *ebiten.Image) {
	// High score at top left
	op := &ebiten.DrawImageOptions{}
	const iconHeight = topAreaHeight * 0.75
	const iconWidth = topAreaHeight * 0.75
	scaleX := iconWidth / float64(resources.FirstPlaceImage.Bounds().Dx())
	scaleY := iconHeight / float64(resources.FirstPlaceImage.Bounds().Dy())
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(0, (topAreaHeight-iconHeight)/2)

	screen.DrawImage(resources.FirstPlaceImage, op)

	highScoreMsg := commaFormatter.Sprintf("%d", g.highScore)
	_, highScoreHeight := getTextSize(highScoreMsg, resources.TextFontFace)
	highScoreColor := displayModeToForegroundColor[g.displayMode]
	if g.cheated {
		highScoreColor = reddishGray
	}
	text.Draw(
		screen,
		highScoreMsg,
		resources.TextFontFace,
		// Text offset is at a weird spot towards the bottom of the letters, so we need to offset it by the height of the
		// text to center it.
		iconWidth, int(((topAreaHeight-highScoreHeight)/2)+highScoreHeight),
		highScoreColor,
	)

	// Score at top right
	scoreMsg := commaFormatter.Sprintf("%d", g.score)
	scoreWidth, scoreHeight := getTextSize(scoreMsg, resources.TextFontFace)
	text.Draw(
		screen,
		scoreMsg,
		resources.TextFontFace,
		// Offset it from the right edge a bit
		int(boardWidth-scoreWidth-20), int(((topAreaHeight-scoreHeight)/2)+scoreHeight),
		displayModeToForegroundColor[g.displayMode],
	)

	// Game over in the middle
	if g.gameOver {
		gameOverMsg := "Game Over"
		gameOverWidth, gameOverHeight := getTextSize(gameOverMsg, resources.TextFontFace)
		restartImageWidth := iconWidth
		restartImageHeight := iconHeight
		gameOverX := int((boardWidth - gameOverWidth - fixed.Int26_6(restartImageWidth)) / 2)
		gameOverY := int(((topAreaHeight - gameOverHeight) / 2) + gameOverHeight)
		text.Draw(
			screen,
			gameOverMsg,
			resources.TextFontFace,
			gameOverX,
			gameOverY,
			displayModeToForegroundColor[g.displayMode],
		)
		// put the restart image next to the game over text
		scaleX := restartImageWidth / float64(resources.RestartImage.Bounds().Dx())
		scaleY := restartImageHeight / float64(resources.RestartImage.Bounds().Dy())
		restartImageX := gameOverX + int(gameOverWidth) + 20
		restartImageY := (topAreaHeight - iconHeight) / 2
		op := &ebiten.DrawImageOptions{}
		op.Filter = ebiten.FilterLinear
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(float64(restartImageX), restartImageY)
		screen.DrawImage(resources.RestartImage, op)

		// Clicked on the restart image
		if g.releaseX >= restartImageX && g.releaseX <= restartImageX+int(restartImageWidth) &&
			g.releaseY >= int(restartImageY) && g.releaseY <= int(restartImageY+restartImageHeight) {
			g.Reset(rand.Uint64())
		}
	}
}

func (g *Game) drawBackground(screen *ebiten.Image) {
	vector.DrawFilledRect(
		screen,
		0, 0,
		float32(screen.Bounds().Max.X),
		float32(screen.Bounds().Max.Y),
		displayModeToBackgroundColor[g.displayMode],
		false,
	)
}

func getTextSize(scoreMsg string, face font.Face) (fixed.Int26_6, fixed.Int26_6) {
	bounds, _ := font.BoundString(face, scoreMsg)

	// Calculate text width and height
	textWidth := (bounds.Max.X - bounds.Min.X) >> 6
	textHeight := (bounds.Max.Y - bounds.Min.Y) >> 6
	return textWidth, textHeight
}

// Layout returns the logical screen dimensions. The game window will scale to fit this size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return boardWidth, topAreaHeight + boardHeight + bottomAreaHeight
}

func (g *Game) drawSplash(screen *ebiten.Image) {
	// Draw the splash screen
	splashBounds := resources.SplashImage.Bounds()
	scaleX := (boardWidth - 50) / float64(splashBounds.Dx())
	scaleY := boardHeight / float64(splashBounds.Dy())
	splashX := 25
	splashY := topAreaHeight
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(splashX), float64(splashY))
	screen.DrawImage(resources.SplashImage, op)
}
