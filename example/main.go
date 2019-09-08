package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/sno6/mnist"

	"image/png"
)

const (
	// Change these to point your image / label files.
	imagesLoc        = "../data/train-images-idx3-ubyte"
	labelsLoc        = "../data/train-labels-idx1-ubyte"
	testingImagesLoc = "../data/t10k-images-idx3-ubyte"
	testingLabelsLoc = "../data/t10k-labels-idx1-ubyte"
)

func main() {
	// Files can be a partial list and this will still work.
	m, err := mnist.New(&mnist.Files{
		TrainingImagesLoc: imagesLoc,
		TrainingLabelsLoc: labelsLoc,
		TestingImagesLoc:  testingImagesLoc,
		TestingLabelsLoc:  testingLabelsLoc,
	})
	if err != nil {
		log.Fatalf("mnist: Error processing MNIST data: %v\n", err)
	}

	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(m.TestingImages.Count)

	img := mnist.ToGrayScale(m.TestingImages.GetImage(i))
	f, err := os.Create(fmt.Sprintf("mnist-image-%v.png", m.TestingLabels.GetLabel(i)))
	if err != nil {
		log.Fatalf("mnist: Error creating file\n", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		log.Fatalf("mnist: Error writing image\n", err)
	}
}
