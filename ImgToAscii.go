package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strconv"

	"github.com/nfnt/resize"
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

func main() {
	var (
		grayscale             [8]string = [8]string{"@", "%", "#", "*", "+", "=", ":", "."}
		imagePath, outputPath string
		threshhold            *uint
		Width, Height         int
	)

	switch len(os.Args) {
	case 4:
		i, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		Width = i
		Height, err = strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		imagePath = os.Args[1]
		outputPath = imagePath[0:len(imagePath)-len("png")] + "txt"
	case 3:
		i, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		threshhold = new(uint)
		*threshhold = uint(i)
		fallthrough
	case 2:
		imagePath = os.Args[1]
		outputPath = imagePath[0:len(imagePath)-len("png")] + "txt"
	case 1:
		imagePath = "image.png"
		outputPath = "image.txt"
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

	if threshhold != nil {
		var width, height = uint(img.Bounds().Max.X), uint(img.Bounds().Max.Y)
		if width > *threshhold {
			img = resize.Resize(uint(*threshhold), 0, img, resize.Lanczos3)
		} else if height > *threshhold {
			img = resize.Resize(0, uint(*threshhold), img, resize.Lanczos3)
		}
	} else {
		img = resize.Resize(uint(Width), uint(Height), img, resize.Lanczos3)
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
