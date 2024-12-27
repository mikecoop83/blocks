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

//go:embed logo.png
var logoImageData []byte
var LogoImage *ebiten.Image

func init() {
	restartImageReader := bytes.NewReader(restartImageData)
	var err error
	RestartImage, _, err = ebitenutil.NewImageFromReader(restartImageReader)
	if err != nil {
		panic(err)
	}
	logoImageReader := bytes.NewReader(logoImageData)
	LogoImage, _, err = ebitenutil.NewImageFromReader(logoImageReader)
	if err != nil {
		panic(err)
	}
}
