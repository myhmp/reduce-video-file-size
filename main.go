package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func renameVideo(path string, info os.FileInfo, prefix string) (fileName string, err error) {
	if info.IsDir() {
		return "", errors.New("Is a directory")
	}

	// source and destionation name
	src := path
	dst := strings.Replace(path, info.Name(), "_"+info.Name(), -1)

	return dst, os.Rename(src, dst)
}

func processVideo(file, prefix string) error {
	src := file
	dst := strings.Replace(file, prefix, "", -1)
	cmd := exec.Command("HandBrakeCLI", "-i", src, "-o", dst, "-e", "x264", "-q", "21", "--preset", "Discord Nitro Small 10-20 Minutes 480p30")
	std, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Println(string(std))

	return nil
}

func main() {
	files := []string{}
	prefix := "_"

	filepath.Walk("./resources", func(path string, info os.FileInfo, err error) error {
		// is a file mp4
		if !info.IsDir() && filepath.Ext(path) == ".mp4" {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			newFileName, err := renameVideo(absPath, info, prefix)
			if err != nil {
				return err
			}

			files = append(files, newFileName)
		}
		if err != nil {
			fmt.Println("ERROR:", err)
		}
		return nil
	})

	if len(files) == 0 {
		fmt.Println("No files found")
		os.Exit(1)
	}

	fmt.Println("The videos were renamed correctly, do you want to continue? the next process is to reduce the size of the video file.")
	fmt.Println("Press 'Enter' to continue...")
	fmt.Scanln()

	var wg sync.WaitGroup
	wg.Add(len(files))

	semaphore := make(chan int, 2)

	for _, file := range files {
		go func(file string) {
			semaphore <- 1
			if err := processVideo(file, prefix); err != nil {
				fmt.Println(file, err)
			}

			wg.Done()
			<-semaphore
		}(file)
	}

	wg.Wait()
}

// HandBrakeCLI -i Vi2.mp4 -o ~/Desktop/Vi22.mp4 -e x264 -q 21 --preset="Discord Nitro Small 10-20 Minutes 480p30"

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"reduce-video-file-size/models"
// 	"strings"
// )

// func main() {
// 	video := models.Reduced{}
// 	videos := []models.Reduced{}

// 	filepath.Walk("./resources", func(path string, info os.FileInfo, err error) error {
// 		if !info.IsDir() && filepath.Ext(path) == ".mp4" {
// 			absPath, err := filepath.Abs(path)
// 			if err != nil {
// 				return err
// 			}

// 			// original data of the video
// 			newFileName := fmt.Sprintf("_%s", info.Name())
// 			dst := strings.Replace(absPath, info.Name(), newFileName, -1)

// 			video.Original.FileName = newFileName
// 			video.Original.Path = dst
// 			video.Original.FileSize = info.Size()

// 			// reduce vídeo file size

// 			videoFile, err := os.Stat(path)
// 			if err != nil {
// 				fmt.Println(err)
// 				return err
// 			}

// 			video.FileName = info.Name()
// 			video.Path = path
// 			video.FileSize = videoFile.Size()
// 			video.ReducedMegabytes = info.Size() - videoFile.Size()

// 			videos = append(videos, video)
// 		}
// 		if err != nil {
// 			fmt.Println("ERROR:", err)
// 		}
// 		return nil
// 	})

// 	file, _ := json.MarshalIndent(videos, "", " ")

// 	_ = ioutil.WriteFile("test.json", file, 0644)

// }

// // HandBrakeCLI -i Vi2.mp4 -o ~/Desktop/Vi22.mp4 -e x264 -q 21 --preset="Discord Nitro Small 10-20 Minutes 480p30"
