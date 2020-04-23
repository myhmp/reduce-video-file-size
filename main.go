package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	err := filepath.Walk("/Users",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fmt.Println(path, info.Name(), info.Size())
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

// HandBrakeCLI -i Vi2.mp4 -o ~/Desktop/Vi22.mp4 -e x264 -q 21 --preset="Discord Nitro Small 10-20 Minutes 480p30"
