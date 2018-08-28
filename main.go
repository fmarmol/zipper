package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func main() {
	directory := flag.String("d", "", "directory to zip")
	flag.Parse()

	if *directory == "" {
		flag.Usage()
		return
	}
	*directory = filepath.Clean(*directory)
	zipName := fmt.Sprintf("%v.zip", *directory)
	zipFile, err := os.Create(zipName)
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	writer := zip.NewWriter(zipFile)
	defer writer.Close()

	filesInfo, err := ioutil.ReadDir(*directory)
	if err != nil {
		panic(err)
	}
	bar := pb.StartNew(len(filesInfo))
	for _, fileInfo := range filesInfo {
		if fileInfo.IsDir() {
			bar.Increment()
			continue
		}

		srcFile, err := os.Open(filepath.Join(*directory, fileInfo.Name()))
		if err != nil {
			log.Println(err)
			bar.Increment()
			continue
		}
		dstFile, err := writer.Create(srcFile.Name())
		if err != nil {
			srcFile.Close()
			log.Println(err)
			bar.Increment()
			continue
		}
		io.Copy(dstFile, srcFile)
		srcFile.Close()
		bar.Increment()
	}
	bar.Finish()
}
