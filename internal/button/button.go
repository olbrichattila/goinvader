// Package button manages clickable buttons
package button

import (
	"image/color"
	"spaceinvader/internal/gametext"

	"github.com/hajimehoshi/ebiten/v2"
)

type Button interface {
	New(caption string, x, y, w, h int, onClick func())
	Remove(caption string)
	Update()
	Render(screen *ebiten.Image)
}

type btnEvent struct {
	hoverBackground *ebiten.Image
	background      *ebiten.Image
	x, y, w, h      int
	onHover         bool
	onClick         func()
}

type btn struct {
	mouseDown bool
	buttons   map[string]*btnEvent
}

func New() Button {
	return &btn{
		buttons: map[string]*btnEvent{},
	}
}

func (b *btn) New(caption string, x, y, w, h int, onClick func()) {
	btnBg := ebiten.NewImage(w, h)
	btnBg.Fill(color.RGBA{240, 33, 33, 255})

	btnBgHover := ebiten.NewImage(w, h)
	btnBgHover.Fill(color.RGBA{240, 120, 120, 255})

	textColor := color.RGBA{255, 255, 255, 255}
	gametext.DrawWithColor(btnBg, caption, 5, 20, textColor)

	gametext.DrawWithColor(btnBgHover, caption, 5, 20, textColor)
	b.buttons[caption] = &btnEvent{
		background:      btnBg,
		hoverBackground: btnBgHover,
		x:               x,
		y:               y,
		w:               w,
		h:               h,
		onClick:         onClick,
	}
}

func (b *btn) Remove(caption string) {
	delete(b.buttons, caption)
}

func (b *btn) Update() {
	justClicked := false
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !b.mouseDown {
			justClicked = true
		}
		b.mouseDown = true
	} else {
		b.mouseDown = false
	}

	x, y := ebiten.CursorPosition()
	for _, btn := range b.buttons {
		if x >= btn.x && x <= btn.x+btn.w && y >= btn.y && y <= btn.y+btn.h {
			btn.onHover = true
			if justClicked && btn.onClick != nil {
				btn.onClick()
			}
		} else {
			btn.onHover = false
		}
	}
}

func (b *btn) Render(screen *ebiten.Image) {
	for caption, btn := range b.buttons {
		b.renderButton(screen, caption, btn)
	}
}

func (b *btn) renderButton(screen *ebiten.Image, caption string, btn *btnEvent) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(btn.x), float64(btn.y))
	if btn.onHover {
		screen.DrawImage(btn.hoverBackground, op)
		return
	}
	screen.DrawImage(btn.background, op)
}
