package imagequant

/*
Copyright (C) 2022 William K. Bae

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"image"
	"image/draw"

	"github.com/WKBae/go-imagequant/internal/cgo"
)

type Options struct {
	MinQuality     int
	MaxQuality     int
	Speed          int
	DitheringLevel float32
	Gamma          float64
}

//go:generate sh ./copy_source.sh

func QuantizeRGBA(img *image.RGBA, opts *Options) (*image.Paletted, error) {
	if opts == nil {
		opts = &Options{
			MinQuality:     0,
			MaxQuality:     100,
			Speed:          4,
			DitheringLevel: 0,
			Gamma:          0,
		}
	}
	return cgo.Quantize(img, opts.MinQuality, opts.MaxQuality, opts.Speed, opts.DitheringLevel, opts.Gamma)
}

func Quantize(img image.Image, opts *Options) (*image.Paletted, error) {
	rgba, ok := img.(*image.RGBA)
	if !ok {
		rect := image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy())
		rgba = image.NewRGBA(rect)
		draw.Draw(rgba, rect, img, img.Bounds().Min, draw.Src)
	}
	return QuantizeRGBA(rgba, opts)
}
