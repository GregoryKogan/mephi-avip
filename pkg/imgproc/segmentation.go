package imgproc

import (
	"image"
	"image/color"
)

type segment struct {
	start, end int
}

func SegmentLetters(img *image.Gray) []image.Rectangle {
	yMin, yMax := verticalBounds(img)
	xSegments := horizontalSegments(img)

	rectangles := make([]image.Rectangle, 0, len(xSegments))
	for _, xSeg := range xSegments {
		r := image.Rect(xSeg.start, yMin, xSeg.end, yMax)
		rectangles = append(rectangles, r)
	}
	return rectangles
}

func verticalBounds(img *image.Gray) (int, int) {
	horProfile := HorizontalProfile(img)
	minInd, maxInd := -1, -1
	for i, v := range horProfile {
		if v > 0 {
			if minInd == -1 {
				minInd = i
			}
			maxInd = i
		}
	}
	return minInd, maxInd
}

func horizontalSegments(img *image.Gray) []segment {
	segments := make([]segment, 0)
	verProfile := VerticalProfile(img)

	curSegment := segment{start: -1}
	for i, v := range verProfile {
		if v > 0 && curSegment.start == -1 {
			curSegment.start = i
		}
		if v == 0 && curSegment.start != -1 {
			curSegment.end = i - 1
			segments = append(segments, curSegment)
			curSegment.start = -1
		}
	}
	return segments
}

func DrawRectangles(img *image.Gray, rectangles []image.Rectangle) image.Image {
	bounds := img.Bounds()
	output := image.NewRGBA(bounds)

	for x := range bounds.Max.X {
		for y := range bounds.Max.Y {
			f := img.GrayAt(x, y).Y
			output.Set(x, y, color.RGBA{f, f, f, 255})
		}
	}

	red := color.RGBA{255, 0, 0, 255}
	for _, r := range rectangles {
		// Top line
		for x := r.Min.X; x <= r.Max.X; x++ {
			output.Set(x, r.Min.Y, red)
		}
		// Bottom line
		for x := r.Min.X; x <= r.Max.X; x++ {
			output.Set(x, r.Max.Y, red)
		}
		// Left line
		for y := r.Min.Y; y <= r.Max.Y; y++ {
			output.Set(r.Min.X, y, red)
		}
		// Right line
		for y := r.Min.Y; y <= r.Max.Y; y++ {
			output.Set(r.Max.X, y, red)
		}
	}
	return output
}
