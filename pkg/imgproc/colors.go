package imgproc

import (
	"image"
	"image/color"
	"math"
	"sync"
)

type HSL struct {
	H, S, L float64
}

func To24bit(c color.Color) (uint8, uint8, uint8) {
	r, g, b := ToNormalized(c)
	return uint8(r * 255.0), uint8(g * 255.0), uint8(b * 255.0)
}

func ToNormalized(c color.Color) (float64, float64, float64) {
	uir, uig, uib, uia := c.RGBA()
	return float64(uir) / float64(uia), float64(uig) / float64(uia), float64(uib) / float64(uia)
}

func ToHSL(c color.Color) HSL {
	r, g, b := ToNormalized(c)

	hsl := HSL{}

	thetaRadians := math.Acos((((r - g) + (r - b)) / 2.0) / (math.Sqrt((r-g)*(r-g) + (r-b)*(g-b))))
	theta := thetaRadians * 180.0 / math.Pi
	if b <= g {
		hsl.H = theta
	} else {
		hsl.H = 360.0 - theta
	}

	hsl.S = 1.0 - 3.0/(r+g+b)*min(r, g, b)
	hsl.L = (r + g + b) / 3.0

	return hsl
}

func GetColorComponents(img image.Image) (image.Image, image.Image, image.Image) {
	bounds := img.Bounds()
	rImg := image.NewRGBA(bounds)
	bImg := image.NewRGBA(bounds)
	gImg := image.NewRGBA(bounds)

	var wg sync.WaitGroup
	wg.Add(bounds.Dx() * bounds.Dy())
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			go func() {
				defer wg.Done()
				r, g, b := To24bit(img.At(x, y))
				rImg.SetRGBA(x, y, color.RGBA{uint8(r), 0, 0, 255})
				gImg.SetRGBA(x, y, color.RGBA{0, uint8(g), 0, 255})
				bImg.SetRGBA(x, y, color.RGBA{0, 0, uint8(b), 255})
			}()
		}
	}

	wg.Wait()

	return rImg, gImg, bImg
}

func GetLightness(img image.Image) *image.Gray {
	bounds := img.Bounds()
	output := image.NewGray(bounds)

	var wg sync.WaitGroup
	wg.Add(bounds.Dx() * bounds.Dy())
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			go func() {
				defer wg.Done()
				hsl := ToHSL(img.At(x, y))
				output.SetGray(x, y, color.Gray{uint8(math.Round(hsl.L * 255.0))})
			}()
		}
	}

	wg.Wait()
	return output
}

func GetHalfTone(img image.Image) *image.Gray {
	bounds := img.Bounds()
	output := image.NewGray(bounds)

	var wg sync.WaitGroup
	wg.Add(bounds.Dx() * bounds.Dy())
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			go func() {
				defer wg.Done()
				r, g, b := ToNormalized(img.At(x, y))
				ht := 0.3*r + 0.59*g + 0.11*b
				output.SetGray(x, y, color.Gray{uint8(ht * 255.0)})
			}()
		}
	}

	wg.Wait()
	return output
}

func InvertColors(img image.Image) image.Image {
	bounds := img.Bounds()
	output := image.NewRGBA(bounds)

	var wg sync.WaitGroup
	wg.Add(bounds.Dx() * bounds.Dy())
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			go func() {
				defer wg.Done()
				r, g, b := To24bit(img.At(x, y))
				output.Set(x, y, color.RGBA{255 - r, 255 - g, 255 - b, 255})
			}()
		}
	}

	wg.Wait()
	return output
}

func FillBackground(img image.Image, fillColor color.Color) image.Image {
	bounds := img.Bounds()
	output := image.NewRGBA(bounds)

	var wg sync.WaitGroup
	wg.Add(bounds.Dx() * bounds.Dy())
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			go func() {
				defer wg.Done()
				output.Set(x, y, img.At(x, y))
				_, _, _, a := img.At(x, y).RGBA()
				if a == 0 {
					output.Set(x, y, fillColor)
				}
			}()
		}
	}

	wg.Wait()
	return output
}
