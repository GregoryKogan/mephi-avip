package imgproc

import (
	"image"
	"image/color"
)

type segment struct {
	start, end int
}

func SegmentLetters(img *image.Gray) []image.Rectangle {
	xSegments := horizontalSegments(img)

	rectangles := make([]image.Rectangle, 0, len(xSegments))
	for _, xSeg := range xSegments {
		minY, maxY := verticalBounds(img, xSeg)
		r := image.Rect(xSeg.start, minY, xSeg.end, maxY)
		rectangles = append(rectangles, r)
	}
	return rectangles
}

func verticalBounds(img *image.Gray, xSeg segment) (int, int) {
	bounds := img.Bounds()
	minY, maxY := -1, -1
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := xSeg.start; x < xSeg.end; x++ {
			f := img.GrayAt(x, y).Y
			if f == 0 {
				continue
			}
			if minY == -1 {
				minY = y
			}
			maxY = y
		}
	}
	return minY, maxY
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
			curSegment.end = i
			segments = append(segments, curSegment)
			curSegment.start = -1
		}
	}
	return segments
}

func DrawRectangles(img *image.Gray, rectangles []image.Rectangle) image.Image {
	bounds := img.Bounds()
	output := image.NewRGBA(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
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
