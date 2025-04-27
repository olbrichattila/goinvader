package game

import (
	"log"
	"spaceinvader/internal/gameloop"

	"github.com/hajimehoshi/ebiten/v2"
)

func Run() {
	ebiten.SetWindowSize(gameloop.ScreenW, gameloop.ScreenH)
	ebiten.SetWindowTitle("Ali(en) Space invader)")

	if err := ebiten.RunGame(gameloop.New()); err != nil {
		log.Fatal(err)
	}
}
