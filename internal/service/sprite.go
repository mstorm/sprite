package service

import (
	"path"
	"strings"
)

type Sprite struct {
	W          int     `json:"width"`
	H          int     `json:"height"`
	X          int     `json:"x"`
	Y          int     `json:"y"`
	PixelRatio float64 `json:"pixelRatio"`
	FilePath   *string `json:"-"`
	Bytes      []byte  `json:"-"`
}

func (s *Sprite) GetName() string {
	return strings.TrimSuffix(path.Base(*s.FilePath), path.Ext(*s.FilePath))
}
