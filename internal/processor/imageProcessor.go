package processor

import (
	"image"

	"github.com/disintegration/imaging"
	"imaging-service/internal/parser"
)

func ProcessImage(img image.Image, opts parser.Options) (image.Image, error) {
	result := img

	if opts.Width > 0 || opts.Height > 0 {
		result = imaging.Resize(result, opts.Width, opts.Height, imaging.Lanczos)
	}

	if opts.Flip {
		result = imaging.FlipH(result)
	}

	for name, val := range opts.Filters {
		switch name {
		case "grayscale":
			result = imaging.Grayscale(result)
		case "blur":
			result = imaging.Blur(result, val)
		case "brightness":
			result = imaging.AdjustBrightness(result, val)
		case "contrast":
			result = imaging.AdjustContrast(result, val)
		}
	}

	return result, nil
}
