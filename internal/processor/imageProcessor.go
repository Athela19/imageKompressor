package processor

import (
	"image"

	"github.com/disintegration/imaging"
	"imaging-service/internal/parser"
)

func ProcessImage(img image.Image, opts parser.Options) (image.Image, error) {
	result := img

	// Resize
	if opts.Width > 0 || opts.Height > 0 {
		result = imaging.Resize(result, opts.Width, opts.Height, imaging.Lanczos)
	}

	// Flip horizontal
	if opts.Flip {
		result = imaging.FlipH(result)
	}

	// Filters
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
