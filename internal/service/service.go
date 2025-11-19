package service

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sort"

	"github.com/mstorm/sprite/internal/image"
)

func NewSprite(filePath string, ratio float64) (*Sprite, error) {
	bytes, attr, err := image.Parse(filePath)
	if err != nil {
		return nil, err
	}
	return &Sprite{
		W:          int(float64(attr.Width()) * ratio),
		H:          int(float64(attr.Height()) * ratio),
		X:          0,
		Y:          0,
		PixelRatio: ratio,
		FilePath:   &filePath,
		Bytes:      bytes,
	}, nil
}

func SaveJsonSpriteMap(filePath string, sprites Sprites) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.Marshal(sprites.ToMap())
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func Load(files []string, ratio float64) (Sprites, error) {
	var sprites Sprites
	for _, svgFile := range files {
		newSprite, err := NewSprite(svgFile, ratio)
		if err != nil {
			return nil, err
		}

		sprites = append(sprites, newSprite)
	}

	return sprites, nil
}

// shelf represents a horizontal row in the bin packing algorithm
type shelf struct {
	x      int
	y      int
	width  int
	height int
}

// binPack implements the Shelf First Fit Decreasing Height (FFDH) algorithm
// Returns the optimal width and height for the sprite sheet
func binPack(sprites Sprites) (width, height int) {
	if len(sprites) == 0 {
		return 0, 0
	}

	// Calculate total area and estimate optimal width
	totalArea := 0
	maxWidth := 0
	maxHeight := 0
	for _, sprite := range sprites {
		totalArea += sprite.W * sprite.H
		if sprite.W > maxWidth {
			maxWidth = sprite.W
		}
		if sprite.H > maxHeight {
			maxHeight = sprite.H
		}
	}

	// Start with a width that aims for a square-ish layout
	// Use the next power of 2 that fits the estimated square root of total area
	estimatedSize := int(math.Sqrt(float64(totalArea)))
	targetWidth := nextPowerOfTwo(max(estimatedSize, maxWidth))

	// Try different widths to find the most efficient packing
	// Start with the target and allow it to grow if needed
	bestWidth := targetWidth
	bestHeight := packWithWidth(sprites, targetWidth)

	// Try a few larger widths to see if we get better aspect ratio
	for i := 1; i <= 2; i++ {
		tryWidth := targetWidth * (1 << i) // 2x, 4x
		tryHeight := packWithWidth(sprites, tryWidth)

		// Prefer more square-like layouts (aspect ratio closer to 1)
		currentRatio := float64(bestWidth) / float64(bestHeight)
		tryRatio := float64(tryWidth) / float64(tryHeight)

		if math.Abs(tryRatio-1.0) < math.Abs(currentRatio-1.0) {
			bestWidth = tryWidth
			bestHeight = tryHeight
		}
	}

	// Do the final packing with the best width
	height = packWithWidth(sprites, bestWidth)

	return bestWidth, height
}

// packWithWidth packs sprites with a given maximum width using shelf algorithm
// Returns the total height needed
func packWithWidth(sprites Sprites, maxWidth int) int {
	if len(sprites) == 0 {
		return 0
	}

	var shelves []shelf
	currentShelf := shelf{x: 0, y: 0, width: maxWidth, height: 0}

	for _, sprite := range sprites {
		// Check if sprite fits in current shelf
		if currentShelf.x+sprite.W <= maxWidth {
			// Place sprite on current shelf
			sprite.X = currentShelf.x
			sprite.Y = currentShelf.y

			currentShelf.x += sprite.W
			if sprite.H > currentShelf.height {
				currentShelf.height = sprite.H
			}
		} else {
			// Start a new shelf
			shelves = append(shelves, currentShelf)

			newY := currentShelf.y + currentShelf.height
			currentShelf = shelf{x: 0, y: newY, width: maxWidth, height: 0}

			// Place sprite on new shelf
			sprite.X = currentShelf.x
			sprite.Y = currentShelf.y

			currentShelf.x += sprite.W
			currentShelf.height = sprite.H
		}
	}

	// Add the last shelf
	shelves = append(shelves, currentShelf)

	// Calculate total height
	totalHeight := 0
	for _, shelf := range shelves {
		totalHeight = shelf.y + shelf.height
	}

	return totalHeight
}

// nextPowerOfTwo returns the next power of 2 greater than or equal to n
func nextPowerOfTwo(n int) int {
	if n <= 0 {
		return 1
	}

	// Check if already power of 2
	if n&(n-1) == 0 {
		return n
	}

	power := 1
	for power < n {
		power *= 2
	}
	return power
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Gen(name string, files []string) error {
	sprites, err := Load(files, 1.0)
	if err != nil {
		return err
	}

	// Layouts
	sort.Slice(sprites, func(i, j int) bool {
		if sprites[i].H != sprites[j].H {
			return sprites[i].H > sprites[j].H
		}

		if sprites[i].W != sprites[j].W {
			return sprites[i].W > sprites[j].W
		}

		return sprites[i].GetName() < sprites[j].GetName()
	})

	// Bin Packing: Use Shelf FFDH algorithm for optimal packing
	w, h := binPack(sprites)

	// Log sprite positions
	for _, sprite := range sprites {
		log.Printf("%s, x=%d, y=%d, w=%d, h=%d\n", sprite.GetName(), sprite.X, sprite.Y, sprite.W, sprite.H)
	}

	// Calculate packing efficiency
	totalSpriteArea := 0
	for _, sprite := range sprites {
		totalSpriteArea += sprite.W * sprite.H
	}
	sheetArea := w * h
	efficiency := float64(totalSpriteArea) / float64(sheetArea) * 100.0

	log.Printf("SpriteSheet: Width=%d, Height=%d, Area=%d, SpriteArea=%d, Efficiency=%.2f%%\n",
		w, h, sheetArea, totalSpriteArea, efficiency)

	if err = sprites.ExportSVG(name, w, h); err != nil {
		return err
	}

	// ConvertPNG and ExportMap
	for _, ratio := range []int{1, 2, 3} {
		if err = sprites.ExportMap(name, ratio); err != nil {
			return err
		}
		if err = sprites.ConvertPNG(name, ratio); err != nil {
			return err
		}
	}
	spriteMap := fmt.Sprintf("%s.json", name)
	if err = SaveJsonSpriteMap(spriteMap, sprites); err != nil {
		return err
	}

	//if ratio > 1 {
	//	spriteFile = fmt.Sprintf("%s@%dx.png", name, ratio)
	//	spriteMap = fmt.Sprintf("%s@%dx.json", name, ratio)
	//}
	//
	//// Save sprite map
	//err = SaveJsonSpriteMap(spriteMap, sprites)
	//if err != nil {
	//	log.Fatalf("Failed to save SpriteMap %s: %v", spriteMap, err)
	//}

	return nil
}
