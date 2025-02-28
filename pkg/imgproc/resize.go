package imgproc

import (
	"image"
	"sync"
)

type ResizeRatio struct {
	M int
	N int
}

func NewResizeRatio(m int, n int) ResizeRatio {
	return ResizeRatio{M: m, N: n}
}

func Resize(img image.Image, ratio ResizeRatio) image.Image {
	bounds := img.Bounds()

	newWidth := bounds.Dx() * ratio.M / ratio.N
	newHeight := bounds.Dy() * ratio.M / ratio.N

	output := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	var wg sync.WaitGroup
	wg.Add(newWidth * newHeight)
	for x := range newWidth {
		for y := range newHeight {
			go func() {
				defer wg.Done()
				originalX := x * ratio.N / ratio.M
				originalY := y * ratio.N / ratio.M
				output.Set(x, y, img.At(originalX, originalY))
			}()
		}
	}

	wg.Wait()
	return output
}
