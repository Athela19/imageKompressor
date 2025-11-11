package parser

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"imaging-service/pkg/utils"
)

type Options struct {
	Width      int
	Height     int
	SmartCrop  bool
	Flip       bool
	CropRegion [4]int
	Filters    map[string]float64
	Format     string
	Quality    int
	Watermark  string
}

func ParseOptions(path string) (Options, string, error) {
	if path == "" {
		return Options{}, "", errors.New("empty path")
	}

	wmRe := regexp.MustCompile(`watermark\(([^)]+)\)`)
	wmMatch := wmRe.FindStringSubmatch(path)
	var wmURL string
	if len(wmMatch) == 2 {
		wmURL = wmMatch[1]
		path = strings.Replace(path, wmMatch[0], "", 1)
	}

	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return Options{}, "", errors.New("invalid URL format, must be /OPTIONS/ENCODED_URL")
	}

	optStr := parts[0]
	imageURL := parts[1]

	if decoded, err := url.PathUnescape(imageURL); err == nil {
		imageURL = decoded
	}

	imageURL = strings.TrimSpace(imageURL)
	if !strings.HasPrefix(imageURL, "http://") && !strings.HasPrefix(imageURL, "https://") {
		imageURL = "https://" + imageURL
	}

	opts := Options{
		Quality: 75,
		Filters: make(map[string]float64),
	}

	// Pisahkan filter dari optStr jika ada
	filterIndex := strings.Index(optStr, "filters:")
	filterStr := ""
	if filterIndex != -1 {
		filterStr = optStr[filterIndex+len("filters:"):]
		optStr = strings.TrimSuffix(optStr[:filterIndex], ":")
	}

	optParts := strings.Split(optStr, ":")

	// Resize / Crop parsing
	if len(optParts) == 1 {
		// Resize
		size := strings.Split(optParts[0], "x")
		if len(size) == 2 {
			opts.Width, _ = strconv.Atoi(size[0])
			opts.Height, _ = strconv.Atoi(size[1])
			if opts.Width < 0 {
				opts.Flip = true
				opts.Width = -opts.Width
			}
		}
	} else if len(optParts) == 2 {
		// Crop: pastikan kedua bagian valid angka
		isCrop := true
		for _, part := range optParts {
			if !strings.Contains(part, "x") {
				isCrop = false
				break
			}
			nums := strings.Split(part, "x")
			if len(nums) != 2 {
				isCrop = false
				break
			}
			if _, err1 := strconv.Atoi(nums[0]); err1 != nil {
				isCrop = false
				break
			}
			if _, err2 := strconv.Atoi(nums[1]); err2 != nil {
				isCrop = false
				break
			}
		}
		if isCrop {
			xy1 := strings.Split(optParts[0], "x")
			xy2 := strings.Split(optParts[1], "x")
			opts.CropRegion[0], _ = strconv.Atoi(xy1[0])
			opts.CropRegion[1], _ = strconv.Atoi(xy1[1])
			opts.CropRegion[2], _ = strconv.Atoi(xy2[0])
			opts.CropRegion[3], _ = strconv.Atoi(xy2[1])
		} else {
			// Jika bukan crop, treat sebagai resize dengan height = 0
			size := strings.Split(optParts[0], "x")
			if len(size) == 2 {
				opts.Width, _ = strconv.Atoi(size[0])
				opts.Height, _ = strconv.Atoi(size[1])
				if opts.Width < 0 {
					opts.Flip = true
					opts.Width = -opts.Width
				}
			}
		}
	}

	// Smart crop
	if strings.Contains(optStr, "smart") {
		opts.SmartCrop = true
	}

	if opts.SmartCrop {
		fmt.Println("DEBUG: Smart crop is enabled")
		fmt.Println("DEBUG: Fetching image from URL:", imageURL)

		img, err := utils.FetchImage(imageURL)
		if err != nil {
			fmt.Println("DEBUG: Failed to fetch image for smart crop:", err)
		} else {
			fmt.Println("DEBUG: Image fetched successfully")

			cropW, cropH := opts.Width, opts.Height
			if cropW == 0 || cropH == 0 {
				cropW, cropH = 200, 200
			}
			fmt.Println("DEBUG: Crop size:", cropW, "x", cropH)

			mapBright := utils.GetBrightnessMap(img)
			rect := utils.FindMostContrastedRegion(mapBright, cropW, cropH)

			fmt.Println("DEBUG: Smart crop rectangle (most contrasted):", rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)

			opts.CropRegion[0] = rect.Min.X
			opts.CropRegion[1] = rect.Min.Y
			opts.CropRegion[2] = rect.Max.X
			opts.CropRegion[3] = rect.Max.Y
		}
	}

	// Filter parsing
	if filterStr != "" {
		filterParts := strings.Split(filterStr, ":")
		for _, f := range filterParts {
			if f == "" {
				continue
			}

			name := strings.Split(f, "(")[0]
			param := ""
			value := 1.0

			valStr := regexp.MustCompile(`\((.*?)\)`).FindStringSubmatch(f)
			if len(valStr) == 2 {
				param = valStr[1]
				if v, err := strconv.ParseFloat(param, 64); err == nil {
					value = v
				}
			}

			switch name {
			case "format":
				opts.Format = strings.ToLower(param)
			case "quality":
				opts.Quality = int(value)
			case "watermark":
				opts.Watermark = param
			default:
				opts.Filters[name] = value
			}
		}
	}

	// Watermark override
	if wmURL != "" {
		opts.Watermark = wmURL
	}

	fmt.Println("DEBUG PARSER:")
	fmt.Println("ImageURL:", imageURL)
	fmt.Println("Options:", opts)

	return opts, imageURL, nil
}
