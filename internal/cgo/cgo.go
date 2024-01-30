package cgo

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

/*
#cgo CFLAGS: -DNDEBUG
#cgo LDFLAGS: -lm

#include <stdlib.h>
#include "libimagequant.h"

static liq_image *cgo_liq_image_create_rgba_rows(const liq_attr *attr, void *buffer, int offsets[], int width, int height, double gamma) {
	void **rows = malloc(sizeof(void*) * height);
	for (int i = 0; i < height; i++) {
		rows[i] = buffer + offsets[i];
	}
	liq_image *img = liq_image_create_rgba_rows(attr, rows, width, height, gamma);
	free(rows);
	return img;
}
*/
import "C"
import (
	"fmt"
	"image"
	"image/color"
	"runtime"
	"unsafe"
)

func Quantize(img *image.NRGBA, minQuality, maxQuality, speed int, ditheringLevel float32, gamma float64) (*image.Paletted, error) {
	var liqErr C.liq_error
	liqAttr := C.liq_attr_create()
	liqErr = C.liq_set_quality(liqAttr, C.int(minQuality), C.int(maxQuality))
	if err := wrapLiqError(liqErr); err != nil {
		return nil, err
	}
	liqErr = C.liq_set_speed(liqAttr, C.int(speed))
	if err := wrapLiqError(liqErr); err != nil {
		return nil, err
	}
	var liqImage *C.liq_image
	if img.Rect.Min.X != 0 || img.Rect.Min.Y != 0 || img.Stride != img.Rect.Dx()*4 {
		// img is cropped; Use liq_image_create_rgba_rows to pass the buffer with minimal copy
		offsets := make([]C.int, img.Rect.Dy())
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			offsets[y-img.Rect.Min.Y] = C.int(img.PixOffset(0, y))
		}
		liqImage = C.cgo_liq_image_create_rgba_rows(liqAttr, unsafe.Pointer(&img.Pix[0]), &offsets[0], C.int(img.Rect.Dx()), C.int(img.Rect.Dy()), C.double(gamma))
	} else {
		liqImage = C.liq_image_create_rgba(liqAttr, unsafe.Pointer(&img.Pix[0]), C.int(img.Rect.Dx()), C.int(img.Rect.Dy()), 0)
	}
	var liqResult *C.liq_result
	liqErr = C.liq_image_quantize(liqImage, liqAttr, &liqResult)
	if err := wrapLiqError(liqErr); err != nil {
		return nil, err
	}
	liqErr = C.liq_set_dithering_level(liqResult, C.float(ditheringLevel))
	if err := wrapLiqError(liqErr); err != nil {
		return nil, err
	}
	remapped := make([]byte, img.Rect.Dx()*img.Rect.Dy())
	liqErr = C.liq_write_remapped_image(liqResult, liqImage, unsafe.Pointer(&remapped[0]), C.ulong(len(remapped)))
	if err := wrapLiqError(liqErr); err != nil {
		return nil, err
	}
	liqPalette := C.liq_get_palette(liqResult)
	liqPaletteLen := int(liqPalette.count)
	palette := make(color.Palette, liqPaletteLen)
	for i := 0; i < liqPaletteLen; i++ {
		palette[i] = color.NRGBA{
			R: byte(liqPalette.entries[i].r),
			G: byte(liqPalette.entries[i].g),
			B: byte(liqPalette.entries[i].b),
			A: byte(liqPalette.entries[i].a),
		}
	}
	C.liq_result_destroy(liqResult)
	C.liq_image_destroy(liqImage)
	C.liq_attr_destroy(liqAttr)
	runtime.KeepAlive(img)
	return &image.Paletted{
		Pix:     remapped,
		Stride:  img.Rect.Dx(),
		Rect:    image.Rect(0, 0, img.Rect.Dx(), img.Rect.Dy()),
		Palette: palette,
	}, nil
}

func wrapLiqError(err C.liq_error) error {
	switch err {
	case C.LIQ_OK:
		return nil
	case C.LIQ_QUALITY_TOO_LOW:
		return fmt.Errorf("libimagequant returned error: QUALITY_TOO_LOW")
	case C.LIQ_VALUE_OUT_OF_RANGE:
		return fmt.Errorf("libimagequant returned error: VALUE_OUT_OF_RANGE")
	case C.LIQ_OUT_OF_MEMORY:
		return fmt.Errorf("libimagequant returned error: OUT_OF_MEMORY")
	case C.LIQ_ABORTED:
		return fmt.Errorf("libimagequant returned error: ABORTED")
	case C.LIQ_BITMAP_NOT_AVAILABLE:
		return fmt.Errorf("libimagequant returned error: BITMAP_NOT_AVAILABLE")
	case C.LIQ_BUFFER_TOO_SMALL:
		return fmt.Errorf("libimagequant returned error: BUFFER_TOO_SMALL")
	case C.LIQ_INVALID_POINTER:
		return fmt.Errorf("libimagequant returned error: INVALID_POINTER")
	default:
		return fmt.Errorf("libimagequant returned unknown error: %d", int(err))
	}
}
