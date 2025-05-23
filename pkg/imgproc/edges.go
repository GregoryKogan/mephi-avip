package imgproc

import (
	"image"
	"image/color"
	"math"
	"slices"
	"sync"
)

func GetScharrEdges(img *image.Gray) (*image.Gray, *image.Gray, *image.Gray) {
	bounds := img.Bounds()

	unnormalizedGx := NewUnnormalizedGray(bounds.Dx(), bounds.Dy())
	unnormalizedGy := NewUnnormalizedGray(bounds.Dx(), bounds.Dy())

	var wg sync.WaitGroup
	wg.Add(bounds.Dx() * bounds.Dy() * 2)
	for x := range bounds.Max.X {
		for y := range bounds.Max.Y {
			index := y*unnormalizedGx.width + x
			go func() {
				defer wg.Done()
				unnormalizedGx.values[index] = Pixel3x3Convolution(img, x, y, [3][3]int{
					{3, 10, 3},
					{0, 0, 0},
					{-3, -10, -3},
				})
			}()
			go func() {
				defer wg.Done()
				unnormalizedGy.values[index] = Pixel3x3Convolution(img, x, y, [3][3]int{
					{3, 0, -3},
					{10, 0, -10},
					{3, 0, -3},
				})
			}()
		}
	}

	wg.Wait()

	unnormalizedG := UnnormalizedGrayAbsSum(unnormalizedGx, unnormalizedGy)
	return unnormalizedGx.Normalize(), unnormalizedGy.Normalize(), unnormalizedG.Normalize()
}

func Pixel3x3Convolution(img *image.Gray, x, y int, weights [3][3]int) int {
	bounds := img.Bounds()
	sx, ex := max(x-1, 0), min(x+1, bounds.Max.X)
	sy, ey := max(y-1, 0), min(y+1, bounds.Max.Y)

	weightedSum := 0
	for xi := sx; xi <= ex; xi++ {
		for yi := sy; yi <= ey; yi++ {
			weight := weights[yi-y+1][xi-x+1]
			weightedSum += int(img.GrayAt(xi, yi).Y) * weight
		}
	}
	return weightedSum
}

type UnnormalizedGray struct {
	width, height int
	values        []int
}

func NewUnnormalizedGray(w, h int) UnnormalizedGray {
	return UnnormalizedGray{width: w, height: h, values: make([]int, w*h)}
}

func UnnormalizedGrayAbsSum(a, b UnnormalizedGray) UnnormalizedGray {
	if a.width != b.width || a.height != b.height {
		panic("UnnormalizedGraySum: sizes do not match")
	}
	output := NewUnnormalizedGray(a.width, a.height)
	for i := range a.values {
		output.values[i] = int(math.Abs(float64(a.values[i])) + math.Abs(float64(b.values[i])))
	}
	return output
}

func (u *UnnormalizedGray) Normalize() *image.Gray {
	minValue := slices.Min(u.values)
	maxValue := slices.Max(u.values)

	output := image.NewGray(image.Rect(0, 0, u.width, u.height))
	for x := range u.width {
		for y := range u.height {
			index := y*u.width + x
			normValue := uint8(255.0 * float64(u.values[index]-minValue) / float64(maxValue-minValue))
			output.SetGray(x, y, color.Gray{normValue})
		}
	}

	return output
}
