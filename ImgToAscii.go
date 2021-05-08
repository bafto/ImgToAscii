package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func mapRange(x, in_min, in_max, out_min, out_max int) int {
	return (x-in_min)*(out_max-out_min)/(in_max-in_min) + out_min
}

func getImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(file)
	file.Close()
	return img, err
}

/*func scaleImage(maxWidth, maxHeight int, img *image.Image) error {
	var image image.Image
	image.Bounds().
}*/

func main() {
	var grayscale [8]string = [8]string{"@", "%", "#", "*", "+", "=", ":", "."}

	var imagePath, outputPath string

	switch len(os.Args) {
	case 1:
		imagePath = "image.png"
		outputPath = "image.txt"
	case 2:
		imagePath = os.Args[1]
		outputPath = imagePath[0:len(imagePath)-len("png")] + "txt"
	case 3:
		imagePath = os.Args[1]
		outputPath = os.Args[2]
	default:
		log.Fatal("Bad path arguments")
	}

	img, err := getImage(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixel := (r + g + b) / 3
			fmt.Fprintf(outputFile, "%v", grayscale[mapRange(int(pixel), 0, 65535, 0, 7)])
		}
		fmt.Fprintf(outputFile, "\n")
	}
}
