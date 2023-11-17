package image

import (
	"encoding/xml"
	"io"
	"os"
	"strconv"
	"strings"
)

type SVGAttr struct {
	W       int    `xml:"width,attr"`
	H       int    `xml:"height,attr"`
	ViewBox string `xml:"viewBox,attr"`
}

func (s *SVGAttr) v() []int {
	if len(s.ViewBox) <= 0 {
		return nil
	}

	b := strings.Split(s.ViewBox, " ")
	if len(b) <= 0 {
		return nil
	}

	var r []int
	for _, i := range b {
		n, err := strconv.Atoi(i)
		if err != nil {
			return nil
		}
		r = append(r, n)
	}
	return r
}

func (s *SVGAttr) Width() int {
	v := s.v()
	if v != nil && len(v) > 0 {
		return v[2]
	}

	return s.W
}

func (s *SVGAttr) Height() int {
	v := s.v()
	if v != nil && len(v) > 0 {
		return v[3]
	}

	return s.H
}

func Parse(filename string) ([]byte, *SVGAttr, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}

	var svg SVGAttr
	if err := xml.Unmarshal(bytes, &svg); err != nil {
		return nil, nil, err
	}

	return bytes, &svg, nil
}
