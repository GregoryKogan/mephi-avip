package imgproc

import (
	"image"
	"math"
)

func QuartersBlackWeight(img *image.Gray) (float64, float64, float64, float64) {
	bounds := img.Bounds()

	minX, maxX := bounds.Min.X, bounds.Max.X
	minY, maxY := bounds.Min.Y, bounds.Max.Y
	cenX := bounds.Min.X + bounds.Dx()/2
	cenY := bounds.Min.Y + bounds.Dy()/2

	q1 := BlackWeight(img, image.Rect(minX, minY, cenX, cenY))
	q2 := BlackWeight(img, image.Rect(cenX, minY, maxX, cenY))
	q3 := BlackWeight(img, image.Rect(minX, cenY, cenX, maxY))
	q4 := BlackWeight(img, image.Rect(cenX, cenY, maxX, maxY))
	return q1, q2, q3, q4
}

func QuartersRelativeBlackWeight(img *image.Gray) (float64, float64, float64, float64) {
	bounds := img.Bounds()
	area := float64(bounds.Dx() / 2 * bounds.Dy() / 2)
	q1, q2, q3, q4 := QuartersBlackWeight(img)
	return q1 / area, q2 / area, q3 / area, q4 / area
}

func BlackWeight(img *image.Gray, bounds image.Rectangle) float64 {
	sum := 0.0
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			sum += float64(img.GrayAt(x, y).Y) / 255.0
		}
	}
	return sum
}

func CenterOfMass(img *image.Gray) (float64, float64) {
	bounds := img.Bounds()

	xSum, ySum := 0.0, 0.0
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			f := float64(img.GrayAt(x, y).Y) / 255.0
			xSum += float64(x) * f
			ySum += float64(y) * f
		}
	}
	weight := BlackWeight(img, bounds)
	return xSum / weight, ySum / weight
}

func RelativeCenterOfMass(img *image.Gray) (float64, float64) {
	bounds := img.Bounds()
	cX, cY := CenterOfMass(img)
	return (cX - 1.0) / (float64(bounds.Dx()) - 1.0), (cY - 1.0) / (float64(bounds.Dy()) - 1.0)
}

func Inertia(img *image.Gray) (float64, float64) {
	Ix, Iy := 0.0, 0.0
	cX, cY := CenterOfMass(img)

	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			f := float64(img.GrayAt(x, y).Y) / 255.0
			Ix += math.Pow(float64(y)-cY, 2.0) * f
			Iy += math.Pow(float64(x)-cX, 2.0) * f
		}
	}

	return Ix, Iy
}

func RelativeInertia(img *image.Gray) (float64, float64) {
	Ix, Iy := Inertia(img)
	bounds := img.Bounds()
	denominator := float64(bounds.Dx() * bounds.Dx() * bounds.Dy() * bounds.Dy())
	return Ix / denominator, Iy / denominator
}

func HorizontalProfile(img *image.Gray) []float64 {
	bounds := img.Bounds()
	lines := make([]float64, 0, bounds.Dy())
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		lineSum := 0.0
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			lineSum += float64(img.GrayAt(x, y).Y) / 255.0
		}
		lines = append(lines, lineSum)
	}

	return lines
}

func VerticalProfile(img *image.Gray) []float64 {
	bounds := img.Bounds()
	columns := make([]float64, 0, bounds.Dx())
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		columnSum := 0.0
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			columnSum += float64(img.GrayAt(x, y).Y) / 255.0
		}
		columns = append(columns, columnSum)
	}

	return columns
}
