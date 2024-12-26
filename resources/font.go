package resources

import (
	_ "embed"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

//go:embed Roboto-Regular.ttf
var robotoFontData []byte

var FontFace font.Face

func init() {
	var err error
	ttf, err := sfnt.Parse(robotoFontData)
	if err != nil {
		panic(err)
	}
	FontFace, err = opentype.NewFace(
		ttf, &opentype.FaceOptions{
			Size: 50,
			DPI:  72,
		},
	)
	if err != nil {
		panic(err)
	}
}
