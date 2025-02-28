package imgproc

import (
	"image"
	"image/color"
	"math"
	"sort"
	"sync"
)

func MedianFilter3x3(img *image.Gray, aperture [3][3]int, k int) *image.Gray {
	bounds := img.Bounds()
	output := image.NewGray(bounds)

	var wg sync.WaitGroup
	wg.Add(bounds.Dx() * bounds.Dy())
	for x := range bounds.Max.X {
		for y := range bounds.Max.Y {
			go func() {
				defer wg.Done()
				processPixelMedianFilter3x3(img, output, aperture, k, x, y)
			}()
		}
	}

	wg.Wait()
	return output
}

func GetDifference(img1, img2 *image.Gray) *image.Gray {
	if img1.Bounds() != img2.Bounds() {
		panic("images sizes do not match")
	}

	bounds := img1.Bounds()
	output := image.NewGray(bounds)

	var wg sync.WaitGroup
	wg.Add(bounds.Dx() * bounds.Dy())
	for x := range bounds.Max.X {
		for y := range bounds.Max.Y {
			go func() {
				defer wg.Done()
				diff := math.Abs(float64(img1.GrayAt(x, y).Y - img2.GrayAt(x, y).Y))
				output.SetGray(x, y, color.Gray{uint8(diff)})
			}()
		}
	}

	wg.Wait()
	return output
}

func processPixelMedianFilter3x3(src, dest *image.Gray, aperture [3][3]int, k, x, y int) {
	bounds := src.Bounds()
	sx, ex := max(x-1, 0), min(x+1, bounds.Max.X)
	sy, ey := max(y-1, 0), min(y+1, bounds.Max.Y)

	values := make([]int, 0, 9)
	for xi := sx; xi <= ex; xi++ {
		for yi := sy; yi <= ey; yi++ {
			ai := yi - y + 1
			aj := xi - x + 1
			values = append(values, int(src.GrayAt(xi, yi).Y)*aperture[ai][aj])
		}
	}
	sort.Ints(values)
	rangMedian := values[len(values)-k-1]
	dest.SetGray(x, y, color.Gray{uint8(rangMedian)})
}
