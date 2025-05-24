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

func LogLightnessAdjust(img image.Image) image.Image {
	bounds := img.Bounds()
	output := image.NewRGBA(bounds)

	maxLightness := 0.0
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			hsl := ToHSL(img.At(x, y))
			maxLightness = max(maxLightness, hsl.L)
		}
	}

	c := 1 / math.Log(1+maxLightness)
	var wg sync.WaitGroup
	wg.Add(bounds.Dx() * bounds.Dy())
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			go func() {
				defer wg.Done()
				hsl := ToHSL(img.At(x, y))
				hsl.L = c * math.Log(1+hsl.L)
				r, g, b := hsl.ToRgb()
				output.SetRGBA(x, y, color.RGBA{r, g, b, 255})
			}()
		}
	}

	return output
}

func (hsl *HSL) ToRgb() (r, g, b uint8) {
	return hsiToRgb(hsl.H, hsl.S, hsl.L)
}

func hsiToRgb(h, s, i float64) (r, g, b uint8) {
	var rf, gf, bf float64 // Intermediate float R, G, B values in [0,1] range

	// Handle achromatic (gray) case separately for precision and simplicity
	if s == 0 {
		rf = i
		gf = i
		bf = i
	} else {
		// Normalize H to be in [0, 360) range. 360 degrees is equivalent to 0 degrees.
		if h >= 360.0 {
			h = 0.0
		} else if h < 0.0 { // Ensure H is not negative, though input spec is [0, 360]
			h = math.Mod(h, 360.0)
			if h < 0.0 {
				h += 360.0
			}
		}

		// Convert H to radians for trigonometric functions
		// The formulas are typically structured for H in degrees, adjusted per sector.
		// sixtyRad is 60 degrees in radians
		sixtyRad := math.Pi / 3.0

		if h >= 0 && h < 120 {
			// RG Sector (Red-Green)
			// H is relative to Red axis (0 degrees)
			hRad := h * math.Pi / 180.0

			bf = i * (1 - s)
			// Check for denominator close to zero, though theoretically it shouldn't be for H in this range.
			// cos(60 - H) where H is 0-119.99... means angle is (-59.99... to 60), cos is > 0.5
			rf = i * (1 + (s*math.Cos(hRad))/math.Cos(sixtyRad-hRad))
			gf = 3*i - (rf + bf)
		} else if h >= 120 && h < 240 {
			// GB Sector (Green-Blue)
			// H is relative to Green axis (120 degrees)
			hPrimeRad := (h - 120.0) * math.Pi / 180.0 // H' for this sector

			rf = i * (1 - s)
			gf = i * (1 + (s*math.Cos(hPrimeRad))/math.Cos(sixtyRad-hPrimeRad))
			bf = 3*i - (rf + gf)
		} else { // h >= 240 && h < 360
			// BR Sector (Blue-Red)
			// H is relative to Blue axis (240 degrees)
			hDoublePrimeRad := (h - 240.0) * math.Pi / 180.0 // H'' for this sector

			gf = i * (1 - s)
			bf = i * (1 + (s*math.Cos(hDoublePrimeRad))/math.Cos(sixtyRad-hDoublePrimeRad))
			rf = 3*i - (gf + bf)
		}
	}

	// Clamp RGB values to [0, 1] range before scaling
	rf = math.Max(0, math.Min(1, rf))
	gf = math.Max(0, math.Min(1, gf))
	bf = math.Max(0, math.Min(1, bf))

	// Scale to [0, 255] and convert to uint8
	// math.Round is used for proper rounding instead of truncation
	r = uint8(math.Round(rf * 255.0))
	g = uint8(math.Round(gf * 255.0))
	b = uint8(math.Round(bf * 255.0))

	return
}
