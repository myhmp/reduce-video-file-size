package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reduce-video-file-size/models"
	"strings"
	"sync"
)

func fileSizeInBytes(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}

	fmt.Println("fileSizeInBytes", stat.Size())
	fmt.Println("path", path)

	var bytes int64
	bytes = stat.Size()

	return bytes, nil
}

func renameVideo(path string, info os.FileInfo, prefix string) (*models.Video, string, error) {
	model := models.Video{}
	if info.IsDir() {
		return &model, "", errors.New("Is a directory")
	}

	newFileName := fmt.Sprintf("_%s", info.Name())

	// source and destionation name
	src := path
	dst := strings.Replace(path, info.Name(), newFileName, -1)

	model.Original.FileName = newFileName
	model.Original.Path = dst
	model.Original.Bytes = info.Size()

	model.Reduced.FileName = info.Name()
	model.Reduced.Path = path

	return &model, dst, os.Rename(src, dst)
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

	videos := []models.Video{}
	video := models.Video{}

	filepath.Walk("./resources", func(path string, info os.FileInfo, err error) error {
		// is a file mp4
		if !info.IsDir() && filepath.Ext(path) == ".mp4" {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			_video, file, err := renameVideo(absPath, info, prefix)
			if err != nil {
				return err
			}

			video = *_video

			files = append(files, file)
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

			bytes, err := fileSizeInBytes(video.Reduced.Path)
			if err != nil {
				fmt.Println(file, err)
			}
			video.Reduced.Bytes = bytes
			video.ReducedBytes = video.Original.Bytes - video.Reduced.Bytes
			fmt.Println("bytes", video.Reduced.Bytes)

			wg.Done()
			<-semaphore
			videos = append(videos, video)

		}(file)
	}

	wg.Wait()

	file, _ := json.MarshalIndent(videos, "", " ")

	_ = ioutil.WriteFile("records.json", file, 0644)
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
