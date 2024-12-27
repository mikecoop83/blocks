package resources

import (
	_ "embed"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

//go:embed text.ttf
var textFontData []byte
var TextFontFace font.Face

func init() {
	var err error
	ttf, err := sfnt.Parse(textFontData)
	if err != nil {
		panic(err)
	}
	TextFontFace, err = opentype.NewFace(
		ttf, &opentype.FaceOptions{
			Size: 50,
			DPI:  72,
		},
	)
	if err != nil {
		panic(err)
	}
}
