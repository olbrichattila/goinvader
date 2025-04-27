package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

var ErrExit = errors.New("normal exit")

type Game struct {
	resized bool
}

func (g *Game) Update() error {
	// Update the game logic here
	if g.resized {
		return nil
	}

	g.resized = true
	g.resize()

	return ErrExit
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Draw things on the screen here
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Set the screen size here
	return 640, 480
}

type imageData struct {
	path   string
	width  int
	height int
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ebiten Blank Boilerplate")

	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) resize() {
	resizeInfo := []imageData{
		// 640*480
		{"internal/images/BGS/robotitle.png", 640, 480},
		{"internal/images/BGS/sky.png", 640, 480},
		{"internal/images/youwon.png", 640, 480},
		{"internal/images/lost.png", 640, 480},
		// 40*40
		{"internal/images/robot-fighter.png", 40, 40},
		{"internal/images/Ships/Spaceship.png", 40, 40},
		{"internal/images/Ships/Spaceship2.png", 40, 40},
		{"internal/images/Ships/Spaceship3.png", 40, 40},
		{"internal/images/Ships/Spaceship4.png", 40, 40},
		{"internal/images/Ships/Spaceship5.png", 40, 40},
		{"internal/images/Ships/Spaceship6.png", 40, 40},
		{"internal/images/Ships/Spaceship7.png", 40, 40},
		{"internal/images/Ships/Spaceship8.png", 40, 40},
		// 100*100
		{"internal/images/Ships/Spaceship9.png", 40, 40},
		{"internal/images/Rocks/up00000.png", 100, 100},
		{"internal/images/Rocks/up00001.png", 100, 100},
		{"internal/images/Rocks/up00002.png", 100, 100},
		{"internal/images/Rocks/up00003.png", 100, 100},
		{"internal/images/Rocks/up00004.png", 100, 100},
		{"internal/images/Rocks/up00005.png", 100, 100},
		{"internal/images/Rocks/up00006.png", 100, 100},
		{"internal/images/Rocks/up00007.png", 100, 100},
		{"internal/images/Rocks/up00008.png", 100, 100},
		{"internal/images/Rocks/up00009.png", 100, 100},
		{"internal/images/Rocks/up00010.png", 100, 100},
		// 20*20
		{"internal/images/Objects/xff2.png", 20, 20},
		{"internal/images/Objects/star3.png", 20, 20},
	}

	for _, imageData := range resizeInfo {
		fmt.Println("Resizing" + imageData.path)
		img, _, _, err := g.rescaleImageToFit(imageData.path, imageData.width, imageData.height)
		if err != nil {
			fmt.Println(err)
			continue
		}

		newPath := strings.TrimPrefix(imageData.path, "internal/")

		err = g.ensureFolderForFile(newPath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = g.saveEbitenImageAsPNG(img, newPath)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (g *Game) loadImageFromFile(path string) (*ebiten.Image, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode image (PNG, JPEG, etc.)
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Convert to ebiten.Image
	ebitenImg := ebiten.NewImageFromImage(img)
	return ebitenImg, nil
}

func (g *Game) rescaleImageToFit(fileName string, targetWidth, targetHeight int) (*ebiten.Image, int, int, error) {
	img, err := g.loadImageFromFile(fileName)
	// img, _, err := ebitenutil.NewImageFromFile(imageName)
	if err != nil {
		return nil, 0, 0, err
	}

	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	widthRatio := float64(targetWidth) / float64(origWidth)
	heightRatio := float64(targetHeight) / float64(origHeight)

	scale := math.Min(widthRatio, heightRatio)
	newWidth := int(float64(origWidth) * scale)
	newHeight := int(float64(origHeight) * scale)

	rescaled := ebiten.NewImage(newWidth, newHeight)

	geom := ebiten.GeoM{}
	geom.Scale(scale, scale)

	rescaled.DrawImage(img, &ebiten.DrawImageOptions{
		GeoM: geom,
	})

	return rescaled, newWidth, newHeight, nil
}

func (g *Game) saveEbitenImageAsPNG(img *ebiten.Image, filename string) error {
	w, h := img.Size()

	rgbaImg := ebiten.NewImage(w, h)
	rgbaImg.DrawImage(img, nil)

	pixels := make([]byte, 4*w*h)
	rgbaImg.ReadPixels(pixels)

	goImg := &image.RGBA{
		Pix:    pixels,
		Stride: 4 * w,
		Rect:   image.Rect(0, 0, w, h),
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, goImg)
}

func (g *Game) ensureFolderForFile(filePath string) error {
	dir := filepath.Dir(filePath) // Get the directory part

	// Check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create all necessary parent directories
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create folders: %w", err)
		}
	}
	return nil
}
