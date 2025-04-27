package gametext

import (
	"embed"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var fontFace = loadFontFace()

//go:embed fonts/*
var embeddedFonts embed.FS

func loadFontFace() font.Face {
	fontBytes, err := embeddedFonts.ReadFile("fonts/Righteous-Regular.ttf")
	if err != nil {
		panic(err)
	}

	tt, err := opentype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	const dpi = 72
	fontFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    18,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	return fontFace
}

func Draw(screen *ebiten.Image, drawText string, x, y float64) float64 {
	textColor := color.RGBA{255, 165, 0, 255}

	return DrawWithColor(screen, drawText, x, y, textColor)
}

func DrawWithColor(screen *ebiten.Image, drawText string, x, y float64, textColor color.RGBA) float64 {
	text.Draw(screen, drawText, fontFace, int(x), int(y), textColor)

	return calcTextWidth(drawText)
}

func calcTextWidth(text string) float64 {
	bounds, _ := font.BoundString(fontFace, text)

	return float64(bounds.Max.X-bounds.Min.X) / 64.0
}
