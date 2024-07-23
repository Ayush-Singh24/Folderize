package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var numberOfFiles int = 0

func matchExtension(ext string, exts []string) bool {
	for _, v := range exts {
		if v == ext {
			return true
		}
	}
	return false
}

var imageExts = []string{".jpg", ".jpeg", ".png", ".webp", ".ico", ".gif", ".svg"}
var videoExts = []string{".mp4", ".mkv"}
var pdfExts = []string{".pdf", ".epub"}
var docExts = []string{".docx", ".pptx", ".txt", ".xlsx"}
var zipExts = []string{".zip", ".rar", ".7zip"}
var programExts = []string{".msi", ".exe", ".bat"}

func main() {
	var src, dest string
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("No arguments provided")
		return
	}
	if len(args) > 2 {
		log.Fatal("Too many arguments")
		return
	}

	src = args[0]
	if len(args) == 1 {
		dest = src
	} else {
		dest = args[1]
	}

	if _, err := os.Stat(src); errors.Is(err, os.ErrNotExist) {
		log.Fatal("Path does not exist")
		return
	}

	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		log.Fatal("Could not create the destination directory.")
		return
	}

	list, err := os.ReadDir(src)
	if err != nil {
		log.Fatal("Could not read files at source")
		return
	}

	if err := os.Chdir(dest); err != nil {
		log.Fatal("Could not visit destination")
		return
	}
	os.MkdirAll("images", os.ModePerm)
	os.MkdirAll("pdfs", os.ModePerm)
	os.MkdirAll("docs", os.ModePerm)
	os.MkdirAll("zips", os.ModePerm)
	os.MkdirAll("programs", os.ModePerm)
	os.MkdirAll("videos", os.ModePerm)
	os.MkdirAll("others", os.ModePerm)

	for _, v := range list {
		if v.IsDir() {
			continue
		}
		a, err := v.Info()
		if err != nil {
			fmt.Println(a.Name(), "error here")
		}
		var moveTo string
		ext := filepath.Ext(a.Name())
		switch {
		case matchExtension(ext, imageExts):
			moveTo = "images"
		case matchExtension(ext, pdfExts):
			moveTo = "pdfs"
		case matchExtension(ext, docExts):
			moveTo = "docs"
		case matchExtension(ext, zipExts):
			moveTo = "zips"
		case matchExtension(ext, videoExts):
			moveTo = "images"
		case matchExtension(ext, programExts):
			moveTo = "programs"
		default:
			moveTo = "others"
		}
		if len(moveTo) == 0 {
			continue
		}
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		err = os.Rename(filepath.Join(src, v.Name()), filepath.Join(".", moveTo, v.Name()))
		if err != nil {
			fmt.Println(v.Name(), "error while moving:", err.Error())
			continue
		}
		fmt.Println(v.Name(), "moved")
		numberOfFiles++
	}
	fmt.Println("Number of files moved:", numberOfFiles)
}
