package imgproc

import (
	"image"
	"math"
	"sync"
)

func BrightnessHistogram(img image.Image) []float64 {
	bins := make([]float64, 256)

	bounds := img.Bounds()
	var wg sync.WaitGroup
	wg.Add(bounds.Dx() * bounds.Dy())
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			l := ToHSL(img.At(x, y)).L
			bins[int(l*255)]++
		}
	}

	return bins
}

func HOG(img *image.Gray) []float64 {
	bounds := img.Bounds()
	gx, gy, g := GetScharrEdges(img)
	cellWidth := 8

	histSum := [][9]float64{}
	for cx := bounds.Min.X; cx < bounds.Max.X-cellWidth; cx += cellWidth {
		for cy := bounds.Min.Y; cy < bounds.Max.Y-cellWidth; cy += cellWidth {
			cellHist := [9]float64{}
			for xi := range cellWidth {
				for yi := range cellWidth {
					x := cx + xi
					y := cy + yi
					thetaRad := math.Atan(float64(gy.GrayAt(x, y).Y) / float64(gx.GrayAt(x, y).Y))
					theta := thetaRad * 180.0 / math.Pi
					k := int(theta/20 - 0.00001)
					cellHist[k] += float64(g.GrayAt(x, y).Y) * (1 - (theta-20.0*float64(k))/20.0)
					if k+1 < 9 {
						cellHist[k+1] += float64(g.GrayAt(x, y).Y) * ((theta - 20.0*float64(k)) / 20.0)
					}
				}
			}
			histSum = append(histSum, cellHist)
		}
	}

	n := len(histSum)
	hist := make([]float64, 9)
	for i := range 9 {
		for _, ch := range histSum {
			hist[i] += ch[i]
		}
		hist[i] /= float64(n)
	}

	return hist
}

func HNorm(hog []float64) []float64 {
	sum := 0.0
	for _, x := range hog {
		sum += x * x
	}
	sum += 0.000001
	s := math.Sqrt(sum)
	norm := make([]float64, len(hog))
	for i := range len(hog) {
		norm[i] = hog[i] / s
	}
	return norm
}
