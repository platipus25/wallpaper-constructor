package wallpaperconstructor

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/disintegration/imaging"
)

// ProcessImg layers the input and resizes it to make an image that will fit within the width and height parameters
/*
	Layers:
		1. Bottom layer is semi-opaque solid color background using prominent color
		2. Next layer is a blurred and stretched copy of the image
		3. The foreground is fitted to the frame and the aspect ratio is preserved, then is streched a small ammount
*/
func ProcessImg(width, height int, img image.Image, blurRadius int) image.Image {
	blur := make(chan image.Image)
	foreground := make(chan image.Image)
	background := make(chan image.Image)

	go func(image image.Image) { // blur
		blurred := imaging.Fill(image, width, height, imaging.Center, imaging.Linear)
		blurred = imaging.Blur(blurred, float64(blurRadius))
		fmt.Println("Image Blurred")
		blur <- blurred
	}(img)

	go func(image image.Image) { // foreground
		resized := imaging.Fit(image, width, height, imaging.Lanczos)
		resized = resizeToIsh(resized, width, height, imaging.Lanczos)
		fmt.Println("Image resized")
		foreground <- resized
	}(img)

	go func(image image.Image) { // background
		promColors, err := prominentcolor.Kmeans(img)
		promColor := promColors[0].Color
		if err != nil {
			log.Fatalf("prominentcolor failed: %v", err)
		}
		fmt.Printf("Prominent colors: %v\n", promColors)
		background <- imaging.New(width, height, color.NRGBA{uint8(promColor.R), uint8(promColor.G), uint8(promColor.B), 200})
	}(img)

	// put it all together
	out := <-background
	out = imaging.OverlayCenter(out, <-blur, 0.5)
	out = imaging.PasteCenter(out, <-foreground)
	fmt.Println("Composite compiled")

	return out
}

func resizeToIsh(img image.Image, targetWidth int, targetHeight int, filter imaging.ResampleFilter) *image.NRGBA {
	bounds := img.Bounds().Size()
	width, height := bounds.X, bounds.Y
	aspect := float64(width) / float64(height) // width:height 1:aspect

	fmt.Printf("\tOriginal Aspect Ratio: %v:1\n", aspect)

	deltaWidth := (targetWidth - width) / 4
	deltaHeight := (targetHeight - height) / 5

	finalWidth := width + deltaWidth
	finalHeight := height + deltaHeight

	aspect = float64(finalWidth) / float64(finalHeight)
	fmt.Printf("\tResult Aspect Ratio: %v:1\n", aspect)

	return imaging.Resize(img, finalWidth, finalHeight, filter)
}
