package resources

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed restart.png
var restartImageData []byte
var RestartImage *ebiten.Image

func init() {
	bytesReader := bytes.NewReader(restartImageData)
	var err error
	RestartImage, _, err = ebitenutil.NewImageFromReader(bytesReader)
	if err != nil {
		panic(err)
	}
}
