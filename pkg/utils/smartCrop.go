package utils

import (
	"image"
)

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

func FindDominantObjectRegion(mapBright [][]float64) image.Rectangle {
	height := len(mapBright)
	width := len(mapBright[0])

	grad := make([][]float64, height)
	for y := range grad {
		grad[y] = make([]float64, width)
		for x := range grad[y] {
			if x == 0 || y == 0 || x == width-1 || y == height-1 {
				continue
			}
			dx := mapBright[y][x+1] - mapBright[y][x-1]
			dy := mapBright[y+1][x] - mapBright[y-1][x]
			grad[y][x] = dx*dx + dy*dy
		}
	}

	maxGrad := 0.0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if grad[y][x] > maxGrad {
				maxGrad = grad[y][x]
			}
		}
	}
	if maxGrad == 0 {
		return image.Rect(0, 0, width, height)
	}
	threshold := maxGrad * 0.3 // ambil 30% dari nilai maksimum

	minX, minY := width, height
	maxX, maxY := 0, 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if grad[y][x] > threshold {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	if minX >= maxX || minY >= maxY || (maxX-minX) < width/20 || (maxY-minY) < height/20 {
		return image.Rect(0, 0, width, height)
	}

	paddingX := int(float64(maxX-minX) * 0.1)
	paddingY := int(float64(maxY-minY) * 0.1)

	minX -= paddingX
	minY -= paddingY
	maxX += paddingX
	maxY += paddingY

	if minX < 0 {
		minX = 0
	}
	if minY < 0 {
		minY = 0
	}
	if maxX > width {
		maxX = width
	}
	if maxY > height {
		maxY = height
	}

	return image.Rect(minX, minY, maxX, maxY)
}

