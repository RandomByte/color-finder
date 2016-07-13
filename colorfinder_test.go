package colorfinder

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"os"
	"testing"
)

func TestFind(t *testing.T) {
	file, err := os.Open("pic.jpg")
	if err != nil {
		t.Fatal("Testpic missing", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	bound := img.Bounds()
	m := image.NewRGBA(bound)
	draw.Draw(m, bound, img, bound.Min, draw.Src)

	colors := find(m)
	fmt.Println(colors)
	if colors == (color.RGBA{}) {
		t.Error("No colors returned")
	}
	if colors.R == 0 && colors.G == 0 && colors.B == 0 {
		t.Error("Color is black. Expected something colorfull")
	}
}
