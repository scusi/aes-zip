package main

import (
	"github.com/alexmullins/zip"    // hacked zip lib, supports AES encrypted archives
	"io"
	"os"
	"path/filepath"
	"strings"
    "flag"
)

var zipFile string
var passwd string
var target string
var needBaseDir bool

func init() {
	flag.StringVar(&zipFile, "f", "", "zip file to create")
	flag.StringVar(&passwd, "p", "infected", "password to use")
    flag.StringVar(&target, "d", "./", "file/dir to add to zip")
    flag.BoolVar(&needBaseDir, "includeBaseDir", true, "per default includes the dir structure")
}

func main() {
    flag.Parse()
    //err := zipit(target, zipFile, false)
    err := zipit(target, zipFile, needBaseDir)
    if err != nil {
        panic(err)
    }
}

func zipit(source, target string, needBaseDir bool) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()
    //archive.SetPassword(passwd)

	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			if needBaseDir {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			} else {
				path := strings.TrimPrefix(path, source)
				if len(path) > 0 && (path[0] == '/' || path[0] == '\\') {
					path = path[1:]
				}
				if len(path) == 0 {
					return nil
				}
				header.Name = path
			}
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
        header.SetPassword(passwd)
		writer, err := archive.CreateHeader(header)

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
