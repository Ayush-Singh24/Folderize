package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type FolderConfig struct {
	image   bool
	pdf     bool
	video   bool
	program bool
	zip     bool
	doc     bool
	others  string
}

type Extensions struct {
	imageExts   []string
	videoExts   []string
	pdfExts     []string
	docExts     []string
	zipExts     []string
	programExts []string
}

// pre-defined extensions

var extensions = Extensions{
	imageExts:   []string{".jpg", ".jpeg", ".png", ".webp", ".ico", ".gif", ".svg"},
	videoExts:   []string{".mp4", ".mkv"},
	pdfExts:     []string{".pdf", ".epub"},
	docExts:     []string{".docx", ".pptx", ".txt", ".xlsx"},
	zipExts:     []string{".zip", ".rar", ".7zip"},
	programExts: []string{".msi", ".exe", ".bat"},
}

var folderConfig FolderConfig

var numberOfFiles = 0

func matchExtension(ext string, exts []string) bool {
	for _, v := range exts {
		if v == ext {
			return true
		}
	}
	return false
}

var rootCmd = &cobra.Command{
	Use:   "folderize \"path/to/src/directory\" \"path/to/dest/directory\" [flags] \n\nCAUTION:\n  If you have same folders as (dest/images, dest/videos, dest/docs, dest/programs, dest/pdfs, dest/zips, dest/others), Folderize would overwrite the files already present at destination with the same name",
	Args:  cobra.RangeArgs(1, 2),
	Short: "“Folderize” is a CLI tool that organizes files based on their type into separate folders.",
	RunE: func(cmd *cobra.Command, args []string) error {
		src := args[0]
		var dest string
		if len(args) == 1 {
			dest = src
		} else {
			dest = args[1]
		}
		if _, err := os.Stat(src); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("path does not exist")
		}

		if err := os.MkdirAll(dest, os.ModePerm); err != nil {
			return fmt.Errorf("could not create the destination directory")
		}

		list, err := os.ReadDir(src)
		if err != nil {
			return fmt.Errorf("could not read files at source")
		}

		if err := os.Chdir(dest); err != nil {
			return fmt.Errorf("could not visit destination")
		}

		if !folderConfig.image && !folderConfig.pdf && !folderConfig.doc && !folderConfig.video && !folderConfig.zip && !folderConfig.program && len(folderConfig.others) == 0 {
			folderConfig.image = true
			folderConfig.pdf = true
			folderConfig.video = true
			folderConfig.program = true
			folderConfig.zip = true
			folderConfig.doc = true
			folderConfig.others = "others"
		}

		switch {
		case folderConfig.image:
			os.MkdirAll("images", os.ModePerm)
		case folderConfig.pdf:
			os.MkdirAll("pdfs", os.ModePerm)
		case folderConfig.zip:
			os.MkdirAll("zips", os.ModePerm)
		case folderConfig.video:
			os.MkdirAll("videos", os.ModePerm)
		case folderConfig.program:
			os.MkdirAll("programs", os.ModePerm)
		case folderConfig.doc:
			os.MkdirAll("docs", os.ModePerm)
		case len(folderConfig.others) > 0:
			os.MkdirAll(folderConfig.others, os.ModePerm)
		}

		for _, v := range list {
			if v.IsDir() {
				continue
			}
			a, err := v.Info()
			if err != nil {
				fmt.Println(a.Name(), "can't read this file")
				continue
			}
			var moveTo string
			ext := filepath.Ext(a.Name())
			switch {
			case matchExtension(ext, extensions.imageExts) && folderConfig.image:
				moveTo = "images"
			case matchExtension(ext, extensions.pdfExts) && folderConfig.pdf:
				moveTo = "pdfs"
			case matchExtension(ext, extensions.docExts) && folderConfig.doc:
				moveTo = "docs"
			case matchExtension(ext, extensions.zipExts) && folderConfig.zip:
				moveTo = "zips"
			case matchExtension(ext, extensions.videoExts) && folderConfig.video:
				moveTo = "videos"
			case matchExtension(ext, extensions.programExts) && folderConfig.program:
				moveTo = "programs"
			case len(folderConfig.others) > 0:
				moveTo = folderConfig.others
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
		fmt.Println(numberOfFiles, "moved")
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// flags

func init() {
	rootCmd.Flags().BoolVarP(&(folderConfig.image), "image", "i", false, "folderize images")
	rootCmd.Flags().BoolVarP(&(folderConfig.pdf), "pdf", "p", false, "folderize pdfs")
	rootCmd.Flags().BoolVarP(&(folderConfig.zip), "zip", "z", false, "folderize zips")
	rootCmd.Flags().BoolVarP(&(folderConfig.video), "video", "v", false, "folderize videos")
	rootCmd.Flags().BoolVarP(&(folderConfig.program), "program", "x", false, "folderize programs")
	rootCmd.Flags().BoolVarP(&(folderConfig.doc), "doc", "d", false, "folderize docs")
	rootCmd.Flags().StringVarP(&(folderConfig.others), "others", "o", "", "transfer rest of the files into seperate folder --others=\"folder_name\"")

	rootCmd.Flags().Lookup("others").NoOptDefVal = "others"
}
