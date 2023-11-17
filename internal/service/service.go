package service

import (
	"encoding/json"
	"fmt"
	"log"
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
