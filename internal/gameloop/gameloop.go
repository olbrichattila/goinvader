package gameloop

import (
	"bytes"
	"fmt"
	"io"
	"spaceinvader/internal/api"
	"spaceinvader/internal/button"
	"spaceinvader/internal/gametext"
	"spaceinvader/internal/inputbox"
	"spaceinvader/internal/sprite"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenW    = 640
	ScreenH    = 480
	sampleRate = 44100
)

const (
	statusDrawIntro drawStatus = iota
	statusDrawGame
	statusDrawInputScore
	statusDrawTop10
)

type drawStatus int

type Game interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
}

type images struct {
	bg         *ebiten.Image
	wonImage   *ebiten.Image
	lostImage  *ebiten.Image
	titleImage *ebiten.Image
}

type sounds struct {
	launchSndData []byte
	explosionSnd  *audio.Player
	bgMusic       *audio.Player
}

type gameStatus struct {
	gameOver bool
	loose    bool
	bulletId int
	lives    int
}

type keyboardStatuses struct {
	spaceDown bool
}

type sprites struct {
	playerSprite sprite.Sprite
	bullets      map[int]sprite.Sprite
}

type game struct {
	api               api.APIClient
	drawStatus        drawStatus
	audioContext      *audio.Context
	openScreenButtons button.Button
	winButtons        button.Button
	backButtons       button.Button
	inputBox          inputbox.InputBox
	sprites           sprites
	ufos              *ufos
	images            images
	sounds            sounds
	score             int
	keyboardStatuses  keyboardStatuses
	gameStatus        gameStatus
	level             int
	userScores        api.UserScores
}

func New() Game {
	g := &game{
		audioContext: audio.NewContext(sampleRate),
		inputBox:     inputbox.New(),
		api:          api.New(),
		level:        0,
	}

	g.preInit()
	g.init()
	return g
}

func (g *game) preInit() {
	g.loadSprites()
	g.loadImages()
	g.loadSounds()
	g.initiateButtons()

}

func (g *game) init() {
	g.gameStatus.lives = 3
	g.score = 0
	g.ufos = newUfoList(g.sprites.playerSprite, g.explosionCallback, g.hitCallback)
	g.sprites.playerSprite.SetY(ScreenH - 50)
	g.gameStatus.gameOver = false
	g.gameStatus.loose = false
}

func (g *game) Update() error {
	switch g.drawStatus {
	case statusDrawIntro:
		g.openScreenButtons.Update()
	case statusDrawInputScore:
		g.inputBox.Update()
		g.winButtons.Update()
	case statusDrawTop10:
		if g.userScores == nil {
			scores, err := g.api.Top10()
			if err != nil {
				return fmt.Errorf("Cannot load from API" + err.Error())
			}
			g.userScores = scores
		}
		g.backButtons.Update()
	default:
		g.playBgMusic()
		if g.handleGameOver() {
			return nil
		}
		g.handleFullScreen()
		g.handleShoot()
		g.handleNavigation()

	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	switch g.drawStatus {
	case statusDrawIntro:
		screen.DrawImage(g.images.titleImage, op)
		g.openScreenButtons.Render(screen)
	case statusDrawTop10:
		gametext.Draw(screen, "TOP 10", 250, 50)
		for i, score := range g.userScores {
			dateStr := score.CreatedAt
			dt1, err := time.Parse(time.RFC3339, score.CreatedAt)
			if err == nil {
				dateStr = dt1.Format("06-01-02 15:04")
			}
			gametext.Draw(screen, score.Name, 50, float64(100+i*25))
			gametext.Draw(screen, strconv.Itoa(score.Score), 250, float64(100+i*25))
			gametext.Draw(screen, dateStr, 370, float64(100+i*25))

		}
		g.backButtons.Render(screen)
	case statusDrawInputScore:
		winnerText := []string{"You can now enter your name"}

		winnerText = append(winnerText, "Your score is "+strconv.Itoa(g.score))
		for i, line := range winnerText {
			gametext.Draw(screen, line, 60, float64(60+i*25))
		}

		g.inputBox.Draw(screen, 200, float64(len(winnerText)*25+85))
		g.winButtons.Render(screen)

	default:
		if g.gameStatus.gameOver {
			screen.DrawImage(g.images.wonImage, op)
			return
		}

		if g.gameStatus.loose {
			screen.DrawImage(g.images.lostImage, op)
			return
		}

		screen.DrawImage(g.images.bg, op)
		g.sprites.playerSprite.Render(screen)
		g.drawBullets(screen)

		won, loose := g.ufos.render(screen)
		if won {
			if g.level == 4 {
				g.gameStatus.gameOver = true
				g.level = 1
				return
			}
			g.level++
			g.ufos.init()
		}

		if loose {
			g.gameStatus.loose = true
		}

		gametext.Draw(screen, "Lives: "+strconv.Itoa(g.gameStatus.lives)+" Level: "+strconv.Itoa(g.level), 20, 30)
	}
}

func (g *game) drawBullets(screen *ebiten.Image) {
	for _, bullet := range g.sprites.bullets {
		bullet.Render(screen)
		if bullet.IsMoving() == false {
			bullet.Close()
			delete(g.sprites.bullets, bullet.Id())
		}
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenW, ScreenH
}

func (g *game) handleGameOver() bool {
	if g.gameStatus.gameOver || g.gameStatus.loose {
		g.drawStatus = statusDrawInputScore
		// if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		// 	g.drawStatus = statusDrawInputScore
		// }

		return true
	}

	return false
}

func (*game) handleFullScreen() {
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
		if !ebiten.IsFullscreen() {
			ebiten.SetWindowSize(ScreenW, ScreenH)
		}
	}
}

func (g *game) handleShoot() {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if !g.keyboardStatuses.spaceDown {
			g.playLaunchSound()
			g.gameStatus.bulletId++
			collisionSprites := []sprite.Sprite{}
			for _, cSprite := range g.ufos.ufoList {
				if cSprite != nil {
					collisionSprites = append(collisionSprites, cSprite)
				}
			}
			bullet := sprite.New(
				[]string{
					"internal/images/Objects/star3.png",
				},
				20, 20,
				sprite.SpiteOptions{
					Id:                g.gameStatus.bulletId,
					AnimateOnMove:     true,
					SoftY:             50,
					CollisionSprites:  collisionSprites,
					CollisionCallback: g.handleCollision,
				},
			)

			bullet.SetX(g.sprites.playerSprite.GetX() + 15)
			bullet.SetY(g.sprites.playerSprite.GetY() - 5)
			bullet.Soft(true)
			bullet.SetY(0)

			g.sprites.bullets[g.gameStatus.bulletId] = bullet
		}
		g.keyboardStatuses.spaceDown = true
	} else {
		g.keyboardStatuses.spaceDown = false
	}
}

func (g *game) handleNavigation() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.sprites.playerSprite.MoveX(5)
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.sprites.playerSprite.MoveX(-5)
	}
}

func (g *game) handleCollision(thisSprite sprite.Sprite, collidedSprites []sprite.Sprite) {
	if g.sprites.bullets[thisSprite.Id()] != nil {
		g.sprites.bullets[thisSprite.Id()].Close()
		delete(g.sprites.bullets, thisSprite.Id())
	}

	for _, collidedSprite := range collidedSprites {
		if g.ufos.ufoList[collidedSprite.Id()] != nil {
			g.ufos.ufoList[collidedSprite.Id()].RunAfterAnimation()
		}
	}
	g.score++
	g.playExplosionSound()
}

func (g *game) explosionCallback(sp sprite.Sprite) {
	g.ufos.ufoList[sp.Id()].Close()
	delete(g.ufos.ufoList, sp.Id())
}

func (g *game) hitCallback(bombSprite sprite.Sprite) {
	bombSprite.Close()
	delete(g.ufos.bombs, bombSprite.Id())
	g.gameStatus.lives--
	g.sprites.playerSprite.RunAfterAnimation()

	g.playExplosionSound()
	if g.gameStatus.lives == 0 {
		g.level = 1
		g.gameStatus.loose = true
	}
}

func (g *game) loadImages() {
	titleImg, _, _ := ebitenutil.NewImageFromFile("internal/images/BGS/robotitle.png")
	bgImg, _, _ := ebitenutil.NewImageFromFile("internal/images/BGS/sky.png")
	wonImg, _, _ := ebitenutil.NewImageFromFile("internal/images/youwon.png")
	lostImg, _, _ := ebitenutil.NewImageFromFile("internal/images/lost.png")

	g.images = images{
		bg:         bgImg,
		wonImage:   wonImg,
		lostImage:  lostImg,
		titleImage: titleImg,
	}
}

func (g *game) loadSounds() {
	explosionSnd, _ := g.loadMp3Sound("internal/sound/fire.mp3")
	launchSndData, _ := g.loadMp3SoundData("internal/sound/rlauncher.mp3")
	bgMusic, _ := g.loadMp3Sound("internal/sound/music.mp3")

	g.sounds = sounds{
		explosionSnd:  explosionSnd,
		launchSndData: launchSndData,
		bgMusic:       bgMusic,
	}
}

func (g *game) initiateButtons() {
	g.openScreenButtons = button.New()
	g.openScreenButtons.New("Play the game", 50, 400, 130, 28, func() {
		g.drawStatus = statusDrawGame
		g.init()
	})
	g.openScreenButtons.New("Display scores", 450, 400, 135, 28, func() {
		g.drawStatus = statusDrawTop10
	})

	g.winButtons = button.New()
	g.winButtons.New("Cancel", 50, 400, 70, 28, func() {
		g.drawStatus = statusDrawIntro
		g.userScores = nil
	})
	g.winButtons.New("Save", 500, 400, 55, 28, func() {
		if len(g.inputBox.Text()) >= 3 {
			err := g.api.AddScore(g.inputBox.Text(), g.score)
			if err != nil {
				fmt.Println(err)
			}
			g.drawStatus = statusDrawIntro
			g.userScores = nil
		}
	})

	g.backButtons = button.New()
	g.backButtons.New("OK", 250, 400, 40, 28, func() {
		g.drawStatus = statusDrawIntro
		g.userScores = nil
	})
}

func (g *game) loadSprites() {
	g.sprites = sprites{
		playerSprite: sprite.New(
			[]string{
				"internal/images/robot-fighter.png",
			},
			50, 50,
			sprite.SpiteOptions{
				Soft:                 true,
				SoftY:                60,
				AnimateOnMove:        true,
				AfterAnimationImages: g.ufos.getExplosionImages(),
			},
		),
		bullets: map[int]sprite.Sprite{},
	}
}

func (g *game) loadMp3Sound(path string) (*audio.Player, error) {
	d, err := g.loadMp3SoundData(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	audioStream, err := mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(d))

	return g.audioContext.NewPlayer(audioStream)
}

func (g *game) loadMp3SoundData(path string) ([]byte, error) {
	file, err := ebitenutil.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	return io.ReadAll(file)
}

func (g *game) playExplosionSound() {
	if !g.sounds.explosionSnd.IsPlaying() {
		g.sounds.explosionSnd.Rewind()
		g.sounds.explosionSnd.Play()
	}
}

func (g *game) playLaunchSound() {
	go func() {
		audioStream, _ := mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(g.sounds.launchSndData))
		sound, _ := g.audioContext.NewPlayer(audioStream)

		if !sound.IsPlaying() {
			sound.Rewind()
			sound.Play()
		}
	}()
}

func (g *game) playBgMusic() {
	if g.sounds.bgMusic == nil {
		return
	}

	if !g.sounds.bgMusic.IsPlaying() {
		g.sounds.bgMusic.Rewind()
		g.sounds.bgMusic.Play()
	}
}
