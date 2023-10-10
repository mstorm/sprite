package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

type Sprites []*Sprite

func (s *Sprites) ToMap() map[string]*Sprite {
	m := make(map[string]*Sprite)
	for _, sprite := range *s {
		m[sprite.GetName()] = sprite
	}

	return m
}

func (s *Sprites) GetMaxWidth() int {
	w := 0
	for _, sprite := range *s {
		if sprite.W > w {
			w = sprite.W
		}
	}
	return w
}

func (s *Sprites) GetMaxHeight() int {
	h := 0
	for _, sprite := range *s {
		if sprite.H > h {
			h = sprite.H
		}
	}
	return h
}

func (s *Sprites) GetSumWidth() int {
	w := 0
	for _, sprite := range *s {
		w += sprite.W
	}
	return w
}

func (s *Sprites) GetSumHeight() int {
	h := 0
	for _, sprite := range *s {
		h += sprite.H
	}
	return h
}

type Sprite struct {
	W          int            `json:"width"`
	H          int            `json:"height"`
	X          int            `json:"x"`
	Y          int            `json:"y"`
	PixelRatio float64        `json:"pixelRatio"`
	FilePath   *string        `json:"-"`
	Icon       *oksvg.SvgIcon `json:"-"`
}

func (s *Sprite) GetName() string {
	return strings.TrimSuffix(path.Base(*s.FilePath), path.Ext(*s.FilePath))
}

func NewSprite(filePath string, ratio float64) (*Sprite, error) {
	icon, _ := oksvg.ReadIcon(filePath)
	//if ratio > 1 {
	//	icon.Transform = icon.Transform.Scale(ratio, ratio)
	//}
	return &Sprite{
		W:          int(icon.ViewBox.W * ratio),
		H:          int(icon.ViewBox.H * ratio),
		X:          0,
		Y:          0,
		PixelRatio: ratio,
		FilePath:   &filePath,
		Icon:       icon,
	}, nil
}

func SaveToPNG(filePath string, img image.Image) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	err = png.Encode(b, img)
	if err != nil {
		return err
	}

	err = b.Flush()
	if err != nil {
		return err
	}
	return nil
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

func main() {
	// Arguments
	name := "sprites"
	svgFiles := os.Args[1:]

	for _, ratio := range []int{1, 2, 3} {
		var sprites Sprites
		for _, svgFile := range svgFiles {
			newSprite, err := NewSprite(svgFile, float64(ratio))
			if err != nil {
				log.Fatalf("Failed to Open file %s: %v", svgFile, err)
			}

			sprites = append(sprites, newSprite)
		}

		// Layouts
		sort.Slice(sprites, func(i, j int) bool {
			return sprites[i].H > sprites[j].H
		})

		// TODO: Bin Packing
		w := sprites.GetMaxWidth()
		h := sprites.GetSumHeight()
		offsetY := 0
		for _, sprite := range sprites {
			sprite.Y = offsetY
			offsetY += sprite.H

			log.Printf("%s, x=%d, y=%d\n", sprite.GetName(), sprite.X, sprite.Y)
		}
		log.Printf("Ratio=%d, MaxWidth=%d, TotalHeight=%d\n", ratio, w, h)

		// Combine images
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		scannerGV := rasterx.NewScannerGV(w, h, img, img.Bounds())
		dasher := rasterx.NewDasher(w, h, scannerGV)

		for _, sprite := range sprites {
			sprite.Icon.Transform = rasterx.Identity.Translate(float64(sprite.X), float64(sprite.Y)).Scale(float64(ratio), float64(ratio))
			sprite.Icon.Draw(dasher, 1.0)
		}

		spriteFile := fmt.Sprintf("%s.png", name)
		spriteMap := fmt.Sprintf("%s.json", name)
		if ratio > 1 {
			spriteFile = fmt.Sprintf("%s@%dx.png", name, ratio)
			spriteMap = fmt.Sprintf("%s@%dx.json", name, ratio)
		}

		// Save sprite image
		err := SaveToPNG(spriteFile, img)
		if err != nil {
			log.Fatalf("Failed to save SpriteImage %s: %v", spriteFile, err)
		}

		// Save sprite map
		err = SaveJsonSpriteMap(spriteMap, sprites)
		if err != nil {
			log.Fatalf("Failed to save SpriteMap %s: %v", spriteMap, err)
		}
	}
}
