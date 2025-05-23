package imgproc

import (
	"image"
	"image/color"
	"math"
	"sync"
)

func GetNiblackThresholding(img *image.Gray, winSize int, k float64) *image.Gray {
	bounds := img.Bounds()
	output := image.NewGray(bounds)

	var wg sync.WaitGroup
	wg.Add(bounds.Dx() * bounds.Dy())
	for x := range bounds.Max.X {
		for y := range bounds.Max.Y {
			go func() {
				defer wg.Done()
				processPixelNiblack(img, output, k, winSize, x, y)
			}()
		}
	}

	wg.Wait()

	return output
}

func GetThresholding(img *image.Gray, threshold uint8) *image.Gray {
	bounds := img.Bounds()
	output := image.NewGray(bounds)

	var wg sync.WaitGroup
	wg.Add(bounds.Dx() * bounds.Dy())
	for x := range bounds.Max.X {
		for y := range bounds.Max.Y {
			go func() {
				defer wg.Done()
				if img.GrayAt(x, y).Y > threshold {
					output.SetGray(x, y, color.Gray{255})
				} else {
					output.SetGray(x, y, color.Gray{0})
				}
			}()
		}
	}

	wg.Wait()
	return output
}

func processPixelNiblack(src, dest *image.Gray, k float64, winSize, x, y int) {
	m := getAverage(src, winSize, x, y)
	s := getStandardDeviation(src, winSize, x, y, m)
	threshold := m + k*s

	value := float64(src.GrayAt(x, y).Y)
	if value > threshold {
		dest.SetGray(x, y, color.Gray{255})
	} else {
		dest.SetGray(x, y, color.Gray{0})
	}
}

func getAverage(img *image.Gray, winSize, x, y int) float64 {
	bounds := img.Bounds()
	halfSize := winSize / 2
	sx, ex := max(x-halfSize, 0), min(x+halfSize, bounds.Max.X)
	sy, ey := max(y-halfSize, 0), min(y+halfSize, bounds.Max.Y)

	sum := 0.0
	n := 0
	for xi := sx; xi <= ex; xi++ {
		for yi := sy; yi <= ey; yi++ {
			sum += float64(img.GrayAt(xi, yi).Y)
			n++
		}
	}
	return sum / float64(n)
}

func getStandardDeviation(img *image.Gray, winSize, x, y int, avg float64) float64 {
	bounds := img.Bounds()
	halfSize := winSize / 2
	sx, ex := max(x-halfSize, 0), min(x+halfSize, bounds.Max.X)
	sy, ey := max(y-halfSize, 0), min(y+halfSize, bounds.Max.Y)

	sqSum := 0.0
	n := 0
	for xi := sx; xi <= ex; xi++ {
		for yi := sy; yi <= ey; yi++ {
			sqSum += math.Pow(float64(img.GrayAt(xi, yi).Y)-avg, 2.0)
			n++
		}
	}
	return math.Sqrt(sqSum / float64(n))
}
