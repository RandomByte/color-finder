package colorfinder

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"
)

type PixelData struct {
	R, G, B, Count, Weight, Degrade int
}

func find(img *image.RGBA) color.RGBA {

	pixels := getImageData(img)

	rgb := PixelData{}

	rgb = getMostProminentRGBImpl(pixels, 6, rgb)
	rgb = getMostProminentRGBImpl(pixels, 4, rgb)
	rgb = getMostProminentRGBImpl(pixels, 2, rgb)
	rgb = getMostProminentRGBImpl(pixels, 0, rgb)

	return color.RGBA{
		R: uint8(rgb.R),
		G: uint8(rgb.G),
		B: uint8(rgb.B),
	}
}

func getImageData(img *image.RGBA) map[string]PixelData {

	pixels := map[string]PixelData{}

	length := len(img.Pix)

	factor := int(math.Max(1, math.Floor(float64(length)/5000+0.5)))

	for i := 4 * factor; i < length; {
		if img.Pix[i+3] > 32 {
			var buffer bytes.Buffer
			buffer.WriteString(strconv.Itoa(int(img.Pix[i])))
			buffer.WriteString(",")
			buffer.WriteString(strconv.Itoa(int(img.Pix[i+1])))
			buffer.WriteString(",")
			buffer.WriteString(strconv.Itoa(int(img.Pix[i+2])))

			key := buffer.String()
			fmt.Println(key)

			pixel, ok := pixels[key]
			if ok {
				pixel.Count++
			} else {
				// create new entry
				newPixel := PixelData{
					R:     int(img.Pix[i]),
					G:     int(img.Pix[i+1]),
					B:     int(img.Pix[i+2]),
					Count: 1,
				}
				newPixel.Weight = favorBright(newPixel.R, newPixel.G, newPixel.B)

				pixels[key] = newPixel
			}
		}
		i += 4 * factor
	}

	return pixels
}

func getMostProminentRGBImpl(pixels map[string]PixelData, degrade uint, rgbMatch PixelData) PixelData {

	rgb := PixelData{
		Degrade: int(degrade),
	}

	db := map[string]int{}

	for _, pixel := range pixels {
		totalWeight := pixel.Weight * pixel.Count
		if doesRgbMatch(rgbMatch, pixel) {
			var buffer bytes.Buffer
			buffer.WriteString(strconv.Itoa(pixel.R >> degrade))
			buffer.WriteString(",")
			buffer.WriteString(strconv.Itoa(pixel.G >> degrade))
			buffer.WriteString(",")
			buffer.WriteString(strconv.Itoa(pixel.B >> degrade))

			pixelGroupKey := buffer.String()
			fmt.Println(pixelGroupKey)

			group, ok := db[pixelGroupKey]
			if ok {
				group += totalWeight
			} else {
				db[pixelGroupKey] = totalWeight
			}
		}
	}

	for key, group := range db {
		rgbs := strings.Split(key, ",")
		r := rgbs[0]
		g := rgbs[1]
		b := rgbs[2]

		count := group

		if count > rgb.Count {
			rgb.Count = count
			rgb.R, _ = strconv.Atoi(r)
			rgb.G, _ = strconv.Atoi(g)
			rgb.B, _ = strconv.Atoi(b)
		}
	}
	return rgb
}

func doesRgbMatch(rgb, pixel PixelData) bool {
	if rgb == (PixelData{}) {
		return true
	}
	degrade := uint(rgb.Degrade)

	r := pixel.R >> degrade
	g := pixel.G >> degrade
	b := pixel.B >> degrade
	return rgb.R == r && rgb.G == g && rgb.B == b
}

func favorBright(r, g, b int) int {
	return r + g + b + 1
}
