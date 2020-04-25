package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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

	var bytes int64
	bytes = stat.Size()

	return bytes, nil
}

func backupVideo(source string, info os.FileInfo, prefix string) error {

	// Open file on disk.
	f, err := os.Open(source)
	if err != nil {
		return err
	}

	// Create a Reader and use ReadAll to get all the bytes from the file.
	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	defer f.Close()

	zipFileName := fmt.Sprintf("%s.zip", info.Name())

	zipFilePath := strings.Replace(source, info.Name(), zipFileName, -1)

	// check if exists
	_, err = os.Stat(zipFilePath)
	if !os.IsNotExist(err) {
		return errors.New("File exits")
	}

	// Open file for writing.
	fileWriter, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}

	defer fileWriter.Close()

	// Write compressed data.
	zw := gzip.NewWriter(fileWriter)

	defer zw.Close()

	_, err = zw.Write(content)
	if err != nil {
		return err
	}

	err = os.Rename(source, strings.Replace(source, info.Name(), prefix+info.Name(), -1))
	if err != nil {
		return err
	}

	return nil
}

func processVideo(file, prefix string) error {
	src := file
	dst := strings.Replace(file, prefix, "", -1)

	cmd := exec.Command("HandBrakeCLI", "-i", src, "-o", dst, "-e", "x264", "-q", "21", "--preset", "Gmail Medium 5 Minutes 480p30")
	std, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Println(string(std))

	return nil
}

func main() {
	prefix := "_"
	count := 0

	record := models.Record{}

	f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		// is a file mp4
		if !info.IsDir() && filepath.Ext(path) == ".mp4" {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			err = backupVideo(absPath, info, prefix)
			if err != nil {
				return err
			}

			model := models.Video{}
			model.Original.FileName = prefix + info.Name()
			model.Original.Path = strings.Replace(absPath, info.Name(), prefix+info.Name(), -1)
			model.Original.Megabytes = (float64(info.Size()) / float64(1024)) / float64(1024)

			model.Reduced.FileName = info.Name()
			model.Reduced.Path = absPath

			record.Videos = append(record.Videos, model)

		}
		if err != nil {
			fmt.Println("ERROR:", err)
		}
		return nil
	})

	if len(record.Videos) == 0 {
		fmt.Println("No files found")
		os.Exit(1)
	}

	fmt.Println("The videos were compressed correctly, do you want to continue? the next process is to reduce the size of the video file.")
	fmt.Println("Press 'Enter' to continue...")
	fmt.Scanln()

	var wg sync.WaitGroup
	wg.Add(len(record.Videos))

	semaphore := make(chan int, 2)

	fmt.Println("starting...")

	for videoIndex := range record.Videos {
		go func(videoIndex int) {
			semaphore <- 1
			if err := processVideo(record.Videos[videoIndex].Original.Path, prefix); err != nil {
				fmt.Println("Process video path:", record.Videos[videoIndex].Original.Path, "error:", err)
				log.Println("Process video path:", record.Videos[videoIndex].Original.Path, "error:", err)
			}

			bytes, err := fileSizeInBytes(record.Videos[videoIndex].Reduced.Path)
			if err != nil {
				fmt.Println("fileSizeInBytes path:", record.Videos[videoIndex].Reduced.Path, "error:", err)
				log.Println("fileSizeInBytes path:", record.Videos[videoIndex].Reduced.Path, "error:", err)
			}
			record.Videos[videoIndex].Reduced.Megabytes = (float64(bytes) / float64(1024)) / float64(1024)
			record.Videos[videoIndex].ReducedMegabytes = record.Videos[videoIndex].Original.Megabytes - record.Videos[videoIndex].Reduced.Megabytes

			// total sum
			record.ReducedMegabytes += record.Videos[videoIndex].ReducedMegabytes

			wg.Done()
			<-semaphore
			count++
			fmt.Println("step", count, "of", len(record.Videos), "finished")
			log.Println("step", count, "of", len(record.Videos), "finished")
		}(videoIndex)
	}

	wg.Wait()

	file, _ := json.MarshalIndent(record, "", " ")

	_ = ioutil.WriteFile("records.json", file, 0644)
}
