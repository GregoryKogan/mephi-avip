package imgproc

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
)

func OpenPNG(path string) (image.Image, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func SavePNG(img image.Image, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		return err
	}

	return nil
}
