package parser

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
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

	optParts := strings.Split(optStr, ":")

	// Resize parsing
	if len(optParts) > 0 {
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

	// Crop parsing
	if len(optParts) > 1 {
		xy1 := strings.Split(optParts[0], "x")
		xy2 := strings.Split(optParts[1], "x")
		if len(xy1) == 2 && len(xy2) == 2 {
			opts.CropRegion[0], _ = strconv.Atoi(xy1[0])
			opts.CropRegion[1], _ = strconv.Atoi(xy1[1])
			opts.CropRegion[2], _ = strconv.Atoi(xy2[0])
			opts.CropRegion[3], _ = strconv.Atoi(xy2[1])
		} else if len(xy2) == 2 {
			// jika hanya ada satu koordinat crop (WxH), anggap x1,y1=0,0
			opts.CropRegion[0], opts.CropRegion[1] = 0, 0
			opts.CropRegion[2], _ = strconv.Atoi(xy2[0])
			opts.CropRegion[3], _ = strconv.Atoi(xy2[1])
		}
	}

	// Smart crop
	if strings.Contains(optStr, "smart") {
		opts.SmartCrop = true
	}

	// Filter parsing
	if strings.Contains(optStr, "filters:") {
		filterStr := strings.SplitN(optStr, "filters:", 2)[1]
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

	// Wa
	if wmURL != "" {
		opts.Watermark = wmURL
	}

	fmt.Println("DEBUG PARSER:")
	fmt.Println("ImageURL:", imageURL)
	fmt.Println("Options:", opts)

	return opts, imageURL, nil
}
