package font

import (
	_ "embed"

	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed Minecraft.ttf
var fontData []byte

var (
    Minecraft font.Face
)

func init() {
	
	// parse mincraft font
    tt, err := opentype.Parse(fontData)
    if err != nil {
        log.Fatal(err)
    }

    Minecraft, err = opentype.NewFace(tt, &opentype.FaceOptions{
        Size:    18,
        DPI:     72,
        Hinting: font.HintingFull,
    })

    if err != nil {
        log.Fatal(err)
    }
}