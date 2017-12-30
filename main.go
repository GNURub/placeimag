package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "medium"
	app.Version = "0.0.1"
	app.Email = "rubencidlara@gmail.com"
	app.Author = "Rubén Cid Lara - @gnurub"
	app.Usage = "Generate the images for the loading of images. Inspired by Medium"
	app.Action = func(c *cli.Context) error {

		// Open the test image.
		filename := c.Args().Get(0)

		if len(c.Args()) > 1 {
			for i := 1; i < len(c.Args()); i++ {
				param := strings.Split(c.Args().Get(i), ":")

				quality, errQ := strconv.ParseInt(param[0], 10, 64)
				blur, errB := strconv.ParseFloat(param[1], 64)
				var width int64 = 0
				var errW error

				if len(param) > 2 && len(param[2]) > 0 {
					width, errW = strconv.ParseInt(param[2], 10, 64)

					if errW != nil {
						width = 0
					}
				}

				if errQ != nil || errB != nil {
					log.Fatalf("Format not correct <quality (int)>:<blur (float)>")
				}

				if quality > 100 {
					quality = 100
				} else if quality < 0 {
					quality = 0
				}

				if blur > 100 {
					blur = 100
				} else if blur < 0 {
					blur = 0
				}

				GenerateImage(filename, int(quality), float64(blur), int(width))
			}
		} else {
			GenerateImage(filename, 30, 4, 100)
			GenerateImage(filename, 70, 0, 0)
		}

		return nil
	}

	app.Run(os.Args)
}

// Generates the optimized image
func GenerateImage(filename string, quality int, blur float64, width int) {

	pwd, _ := os.Getwd()

	file := path.Join(pwd, filename)

	extension := filepath.Ext(filename)
	name := filename[0 : len(filename)-len(extension)]

	finalName := name + "-" + strconv.Itoa(quality) + extension

	src, err := imaging.Open(file)

	if err != nil {
		log.Fatalf("Open failed: %v", err)
	}

	if width > 0 {
		src = imaging.Resize(src, width, 0, imaging.Lanczos)
	}

	// Create a blurred version of the image.
	src = imaging.Blur(src, blur)

	// Save the resulting image using JPEG format.
	err = imaging.Save(src, path.Join(pwd, finalName), imaging.JPEGQuality(quality))
	if err != nil {
		log.Fatalf("Save failed: %v", err)
	}

}
