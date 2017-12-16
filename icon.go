// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import (
	"image"
	"log"
	"os"

	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/nfnt/resize"
)

// Icon represents a icon on the map
type Icon struct {
	MapObject
	Position s2.LatLng
	FilePath string
	Width    uint
	Height   uint
	Offset   []float64
}

// NewIcon creates a new Icon
func NewIcon(pos s2.LatLng, filePath string, size uint) *Icon {
	i := new(Icon)
	i.Position = pos
	i.FilePath = filePath
	i.Width = size
	i.Height = size
	i.Offset = []float64{0, 0}

	return i
}

func (i *Icon) SetWidth(width uint) {
	i.Width = width
}

func (i *Icon) SetHeight(height uint) {
	i.Height = height
}

func (i *Icon) SetOffset(x, y float64) {
	ox := x / float64(i.Width)
	oy := y / float64(i.Height)
	i.Offset = []float64{ox, oy}
}

func (i *Icon) bounds() s2.Rect {
	r := s2.EmptyRect()
	r = r.AddPoint(i.Position)
	return r
}

func (i *Icon) draw(gc *gg.Context, trans *transformer) {
	if !CanDisplay(i.Position) {
		log.Printf("Marker coordinates not displayable: %f/%f", i.Position.Lat.Degrees(), i.Position.Lng.Degrees())
		return
	}

	file, err := os.Open(i.FilePath)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	img = resize.Resize(i.Width, i.Height, img, resize.Lanczos3)

	x, y := trans.ll2p(i.Position)
	ix := int(x) - int(i.Width/2)
	iy := int(y) - int(i.Height/2)

	gc.DrawImageAnchored(img, ix, iy, i.Offset[0], i.Offset[1])
}
