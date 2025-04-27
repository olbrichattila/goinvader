package gameloop

import (
	"math/rand"
	"spaceinvader/internal/sprite"

	"github.com/hajimehoshi/ebiten/v2"
)

func newUfoList(playerSprite sprite.Sprite, explosionCallback func(sprite.Sprite), hitCallback func(sprite.Sprite)) *ufos {
	u := &ufos{
		moveOffset:        2,
		ufoList:           map[int]sprite.Sprite{},
		bombs:             map[int]sprite.Sprite{},
		playerSprite:      playerSprite,
		hitCallback:       hitCallback,
		explosionCallback: explosionCallback,

		randDelay: 200,
	}
	u.explosionImages = u.getExplosionImages()

	return u
}

type ufos struct {
	ufoList           map[int]sprite.Sprite
	playerSprite      sprite.Sprite
	hitCallback       func(sprite.Sprite)
	explosionCallback func(sprite.Sprite)
	explosionImages   []*ebiten.Image
	moveOffset        float64
	blockXLeft        float64
	bombs             map[int]sprite.Sprite
	bombId            int
	frameId           int
	randDelay         int
}

func (u *ufos) init() {

	ufoImages := [][]string{
		[]string{
			"internal/images/Ships/Spaceship.png",
			"internal/images/Ships/Spaceship2.png",
			"internal/images/Ships/Spaceship3.png",
		},
		[]string{
			"internal/images/Ships/Spaceship4.png",
			"internal/images/Ships/Spaceship5.png",
			"internal/images/Ships/Spaceship6.png",
		},
		[]string{
			"internal/images/Ships/Spaceship7.png",
			"internal/images/Ships/Spaceship8.png",
			"internal/images/Ships/Spaceship9.png",
		},
	}
	ufoId := 0
	for h := 0; h < 5; h++ {
		for v := 0; v < 3; v++ {
			u.ufoList[ufoId] =
				sprite.New(
					ufoImages[v],
					40, 40,
					sprite.SpiteOptions{
						X:                            float64(1 + h*80),
						Y:                            float64(45 + v*60),
						Soft:                         true,
						Id:                           ufoId,
						Animate:                      true,
						AfterAnimationImages:         u.explosionImages,
						AfterAnimationCallback:       u.explosionCallback,
						AfterAnimationAnimationDelay: 2,
					},
				)
			ufoId++
		}
	}
}

func (u *ufos) render(screen *ebiten.Image) (bool, bool) {
	hitEnd := false
	var bombingUfo sprite.Sprite

	for h := 0; h < 5; h++ {
		for v := 0; v < 3; v++ {
			newX := u.blockXLeft + float64(1+h*80)
			ufo := u.ufoList[h*3+v]
			if ufo == nil {
				continue
			}
			ufo.Render(screen)
			if ufo.SetX(newX) == false {
				hitEnd = true
			}

			if rand.Intn(8) == 3 {
				bombingUfo = ufo
			}
		}
	}

	if hitEnd {
		u.blockXLeft -= u.moveOffset
		u.moveOffset = -u.moveOffset
		// Speed up
		if u.moveOffset > 0 {
			u.moveOffset += 0.30
		}
		for _, ufo := range u.ufoList {
			if ufo.GetY() > 380 {
				return false, true
			}
			if ufo != nil {
				ufo.MoveY(10)
			}
		}
	} else {
		u.blockXLeft += u.moveOffset
	}

	// Add bombs render
	for _, rBomb := range u.bombs {
		rBomb.Render(screen)
		if !rBomb.IsMoving() {
			delete(u.bombs, rBomb.Id())
		}
	}

	if u.frameId == u.randDelay {

		u.frameId = 0
		u.randDelay = rand.Intn(120) + 20
		if bombingUfo != nil {
			u.bombId++
			bomb := sprite.New(
				[]string{
					"internal/images/Objects/xff2.png",
				},
				20, 20,
				sprite.SpiteOptions{
					Id:               u.bombId,
					AnimateOnMove:    true,
					SoftY:            40,
					SoftSpeedUp:      true,
					CollisionSprites: []sprite.Sprite{u.playerSprite},
					CollisionCallback: func(bombSprite sprite.Sprite, _ []sprite.Sprite) {
						u.hitCallback(bombSprite)
					},
				},
			)
			bomb.SetX(bombingUfo.GetX())
			bomb.SetY(bombingUfo.GetY())

			bomb.Soft(true)
			bomb.SetY(450)

			u.bombs[u.bombId] = bomb
		}
	} else {
		u.frameId++
	}

	return len(u.ufoList) == 0, false
}

func (u *ufos) getExplosionImages() []*ebiten.Image {
	i1, _, _, _ := sprite.RescaleImageToFit("internal/images/Rocks/up00000.png", 100, 100)
	i2, _, _, _ := sprite.RescaleImageToFit("internal/images/Rocks/up00001.png", 100, 100)
	i3, _, _, _ := sprite.RescaleImageToFit("internal/images/Rocks/up00002.png", 100, 100)
	i4, _, _, _ := sprite.RescaleImageToFit("internal/images/Rocks/up00003.png", 100, 100)
	i5, _, _, _ := sprite.RescaleImageToFit("internal/images/Rocks/up00004.png", 100, 100)
	i6, _, _, _ := sprite.RescaleImageToFit("internal/images/Rocks/up00005.png", 100, 100)
	i7, _, _, _ := sprite.RescaleImageToFit("internal/images/Rocks/up00006.png", 100, 100)
	i8, _, _, _ := sprite.RescaleImageToFit("internal/images/Rocks/up00007.png", 100, 100)
	i9, _, _, _ := sprite.RescaleImageToFit("internal/images/Rocks/up00008.png", 100, 100)
	i10, _, _, _ := sprite.RescaleImageToFit("internal/images/Rocks/up00009.png", 100, 100)
	i11, _, _, _ := sprite.RescaleImageToFit("internal/images/Rocks/up00010.png", 100, 100)

	return []*ebiten.Image{
		i1, i2, i3, i4, i5, i6, i7, i8, i9, i10, i11,
	}
}
