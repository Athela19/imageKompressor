package processor

import (
	"fmt"
	"image"
	"strings"

	"github.com/disintegration/imaging"
	"imaging-service/internal/parser"
	"imaging-service/pkg/utils"
)

func ProcessImage(img image.Image, opts parser.Options) (image.Image, error) {
	fmt.Println("DEBUG PROCESSOR: Start processing image")
	fmt.Println("Width:", opts.Width, "Height:", opts.Height, "Flip:", opts.Flip, "Filters:", opts.Filters)
	fmt.Println("Watermark URL raw:", opts.Watermark)
		fmt.Println("CropRegion:", opts.CropRegion)

	result := img
	if opts.CropRegion != [4]int{} {
		x0, y0, x1, y1 := opts.CropRegion[0], opts.CropRegion[1], opts.CropRegion[2], opts.CropRegion[3]
		if x1 > x0 && y1 > y0 && x1 <= result.Bounds().Dx() && y1 <= result.Bounds().Dy() {
			result = imaging.Crop(result, image.Rect(x0, y0, x1, y1))
			fmt.Println("Image cropped:", x0, y0, x1, y1)
		} else {
			fmt.Println("Invalid crop region, skipping crop")
		}
	}

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

	if opts.Watermark != "" {
		wmURL := strings.Trim(opts.Watermark, `"`)
		fmt.Println("Fetching watermark from:", wmURL)
		wm, err := utils.FetchImage(wmURL)
		if err != nil {
			fmt.Println("Failed to fetch watermark:", err)
		} else {
			fmt.Println("Watermark fetched successfully")
			wm = imaging.Resize(wm, result.Bounds().Dx()/5, 0, imaging.Lanczos)
			offset := image.Pt(result.Bounds().Dx()-wm.Bounds().Dx()-10, result.Bounds().Dy()-wm.Bounds().Dy()-10)
			result = imaging.Overlay(result, wm, offset, 0.5)
		}
	}

	fmt.Println("DEBUG PROCESSOR: Done processing image")
	return result, nil
}
