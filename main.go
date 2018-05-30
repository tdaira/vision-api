package main

import (
	"fmt"
	"log"
	"os"

	// Imports the Google Cloud Vision API client package.
	"cloud.google.com/go/vision/apiv1"
	"golang.org/x/net/context"
	"io/ioutil"
	"path/filepath"
)

func main() {
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	htmlFile, err := os.Create("./html/breakfast.html")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer htmlFile.Close()

	htmlHead := "<html><head><meta charset=\"UTF-8\"><title></title></head><body>"
	htmlTail := "</body></html>"
	htmlFile.Write(([]byte)(htmlHead))

	// Get picture paths.
	paths := dirPaths("./data/breakfast")
	for _, path := range paths {
	    fmt.Println("Path: " + path)

		file, err := os.Open(path)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		defer file.Close()
		image, err := vision.NewImageFromReader(file)
		if err != nil {
			log.Fatalf("Failed to create image: %v", err)
		}

		labels, err := client.DetectLabels(ctx, image, nil, 10)
		if err != nil {
			log.Fatalf("Failed to detect labels: %v", err)
		}

		fmt.Println("Labels:")
		labelStr := ""
		for _, label := range labels {
			fmt.Println(label.Description)
			labelStr += label.Description + " "
		}

		htmlImage := fmt.Sprintf(
			"<img src=\"../%s\" style=\"width:auto;height:100px;\"><p>%s</p>",
			path, labelStr)
		htmlFile.Write(([]byte)(htmlImage))
	}

	htmlFile.Write(([]byte)(htmlTail))
}

func dirPaths(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}
