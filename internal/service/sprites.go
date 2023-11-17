package service

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mstorm/sprite/external/inkscape"
)

const (
	FileExtJSON = "json"
	FileExtPNG  = "png"
	FileExtSVG  = "svg"
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

func (s *Sprites) GetSvgFileName(name string) string {
	return getFilename(name, 1, FileExtSVG)
}

func (s *Sprites) ExportSVG(name string, w, h int) error {
	// Combine SVG
	svgFile, err := os.Create(s.GetSvgFileName(name))
	if err != nil {
		return err
	}
	defer svgFile.Close()

	b := bufio.NewWriter(svgFile)
	if _, err = fmt.Fprintf(b, `<svg width="%d" height="%d" fill="none" xmlns="http://www.w3.org/2000/svg">`, w, h); err != nil {
		return err
	}
	for _, sprite := range *s {
		if _, err := fmt.Fprintf(b, `<g transform="translate(%d,%d)">%s</g>`, sprite.X, sprite.Y, sprite.Bytes); err != nil {
			return err
		}
	}
	if _, err = fmt.Fprintf(b, `</svg>`); err != nil {
		return err
	}

	err = b.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (s *Sprites) Scale(ratio int) Sprites {
	var sprites Sprites
	for _, sprite := range *s {
		newSprite := &Sprite{
			W:          sprite.W * ratio,
			H:          sprite.H * ratio,
			X:          sprite.X,
			Y:          sprite.Y,
			PixelRatio: float64(ratio),
			FilePath:   sprite.FilePath,
			Bytes:      sprite.Bytes,
		}

		sprites = append(sprites, newSprite)
	}

	return sprites
}

func getFilename(name string, ratio int, ext string) string {
	if ratio == 1 {
		return fmt.Sprintf("%s.%s", name, ext)
	}

	return fmt.Sprintf("%s@%dx.%s", name, ratio, ext)
}

func (s *Sprites) ExportMap(name string, ratio int) error {
	filename := getFilename(name, ratio, FileExtJSON)
	scaledSprite := s.Scale(ratio)

	if err := SaveJsonSpriteMap(filename, scaledSprite); err != nil {
		return err
	}

	return nil
}

func (s *Sprites) ConvertPNG(name string, ratio int) error {
	filename := getFilename(name, ratio, FileExtPNG)

	if _, err := inkscape.Convert(s.GetSvgFileName(name), filename, float64(ratio)); err != nil {
		return err
	}

	return nil
}
