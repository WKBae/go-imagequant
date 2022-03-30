# go-imagequant

Optimize image size on Go with [libimagequant](https://github.com/ImageOptim/libimagequant) using Go-native [image.Image](https://pkg.go.dev/image#Image) types.

## Example

```go
package main

import (
	"image/png"
	"os"

	"github.com/WKBae/go-imagequant"
)

func main() {
	input, _ := os.Open("input.png")
	defer input.Close()
	srcImg, _ := png.Decode(input)
	resImg, _ := imagequant.Quantize(srcImg, &imagequant.Options{
		MinQuality:     0,
		MaxQuality:     100,
		Speed:          5,
		DitheringLevel: 0.5,
		Gamma:          0,
	})
	output, _ := os.OpenFile("output.png", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	defer output.Close()
	_ = png.Encode(output, resImg)
}
```

