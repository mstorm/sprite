package main

import (
	"bufio"
	"encoding/json"
	"image"
	"image/png"
	"log"
	"os"
	"path"
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

func NewSprite(filePath string) (*Sprite, error) {
	icon, _ := oksvg.ReadIcon(filePath)
	return &Sprite{
		W:          int(icon.ViewBox.W),
		H:          int(icon.ViewBox.H),
		X:          0,
		Y:          0,
		PixelRatio: 1,
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
	spriteFile := "sprites.png"
	spriteMap := "sprites.json"
	svgFiles := os.Args[1:]

	var sprites Sprites
	for _, svgFile := range svgFiles {
		newSprite, err := NewSprite(svgFile)
		if err != nil {
			log.Fatalf("Failed to Open file %s: %v", svgFile, err)
		}

		sprites = append(sprites, newSprite)
	}

	// Layouts
	// TODO: Bin Packing
	w := sprites.GetMaxWidth()
	h := sprites.GetSumHeight()
	offsetY := 0
	for _, sprite := range sprites {
		sprite.Y = offsetY
		offsetY += sprite.H

		log.Printf("%s, x=%d, y=%d\n", sprite.GetName(), sprite.X, sprite.Y)
	}
	log.Printf("MaxWidth=%d, TotalHeight=%d\n", w, h)

	// Combine images
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	scannerGV := rasterx.NewScannerGV(w, h, img, img.Bounds())
	dasher := rasterx.NewDasher(w, h, scannerGV)

	for _, sprite := range sprites {
		sprite.Icon.Transform = rasterx.Identity.Translate(float64(sprite.X), float64(sprite.Y))
		sprite.Icon.Draw(dasher, 1.0)
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
