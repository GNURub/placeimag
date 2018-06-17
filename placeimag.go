package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/urfave/cli"
)

// Listado de ficheros en el argumento
var argsFiles []string

func getFiles(args cli.Args) []string {
	var files []string
	pwd, _ := os.Getwd()
	for i := 0; i < len(args); i++ {
		filepaht := path.Join(pwd, args.Get(i))

		if fileExists(filepaht) {
			argsFiles = append(argsFiles, args.Get(i))
			files = append(files, filepaht)
		} else {
			fs, err := filepath.Glob(filepaht)

			if err == nil && len(fs) > 0 {
				argsFiles = append(argsFiles, args.Get(i))
				files = append(files, fs...)
			}
		}
	}

	return files
}

func getQualities(args cli.Args) []string {
	var qualities []string
	for i := 0; i < len(args); i++ {
		is := false
		for _, f := range argsFiles {
			if f == args.Get(i) {
				is = true
				break
			}
		}

		if !is {
			qualities = append(qualities, args.Get(i))
		}
	}

	return qualities
}

func fileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}

	return false
}

func main() {

	app := cli.NewApp()
	app.Name = "medium"
	app.Version = "0.0.1"
	app.Email = "rubencidlara@gmail.com"
	app.Author = "Rubén Cid Lara - @gnurub"
	app.Usage = "Generate th placeholder images for the loading of images. Inspired by Medium"
	app.Action = func(c *cli.Context) error {

		// Open the test image.
		var files []string
		if len(c.Args()) > 1 {
			files = getFiles(c.Args())
		}

		if len(files) == 0 {
			fmt.Print("File not exists")
			os.Exit(1)
		}

		var qualities = getQualities(c.Args())

		if len(qualities) > 0 {
			for i := 0; i < len(qualities); i++ {

				param := strings.Split(qualities[i], ":")

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
					log.Fatalf("Format not correct <quality (int)>:<blur (float)>[:size]")
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

				for _, file := range files {
					GenerateImage(file, int(quality), float64(blur), int(width))
				}
			}
		} else {
			for _, file := range files {
				GenerateImage(file, 30, 4, 100)
				GenerateImage(file, 70, 0, 0)
			}
		}

		return nil
	}

	app.Run(os.Args)
}

// Generates the optimized image
func GenerateImage(file string, quality int, blur float64, width int) {
	dir, filename := path.Split(file)
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
	err = imaging.Save(src, path.Join(dir, finalName), imaging.JPEGQuality(quality))
	if err != nil {
		log.Fatalf("Save failed: %v", err)
	}

}
