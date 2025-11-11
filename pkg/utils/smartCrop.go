package utils

import (
	"image"
)

// GetBrightnessMap mengubah gambar menjadi peta kecerahan (grayscale)
func GetBrightnessMap(img image.Image) [][]float64 {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	mapBright := make([][]float64, height)
	for y := 0; y < height; y++ {
		mapBright[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			mapBright[y][x] = float64((r>>8 + g>>8 + b>>8) / 3)
		}
	}
	return mapBright
}

// GetLocalVariance menghitung variance brightness di rectangle
func GetLocalVariance(mapBright [][]float64, rect image.Rectangle) float64 {
	sum, sumSq := 0.0, 0.0
	h := rect.Dy()
	w := rect.Dx()
	n := float64(h * w)

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			v := mapBright[y][x]
			sum += v
			sumSq += v * v
		}
	}
	mean := sum / n
	return sumSq/n - mean*mean
}

// FindMostContrastedRegion mencari rectangle dengan variance brightness tertinggi
func FindMostContrastedRegion(mapBright [][]float64, cropWidth, cropHeight int) image.Rectangle {
	height := len(mapBright)
	width := len(mapBright[0])

	var bestRect image.Rectangle
	maxVar := -1.0

	for y := 0; y <= height-cropHeight; y++ {
		for x := 0; x <= width-cropWidth; x++ {
			rect := image.Rect(x, y, x+cropWidth, y+cropHeight)
			variance := GetLocalVariance(mapBright, rect)
			if variance > maxVar {
				maxVar = variance
				bestRect = rect
			}
		}
	}
	return bestRect
}
