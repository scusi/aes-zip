// malware extract - extracts zip files encrypted with the password 'infected',
// which is basically industry standard in the anti-malware industry.
//
package main

import (
	"github.com/alexmullins/zip"    // hacked zip lib, supports AES encrypted archives
	"io"
	"os"
	"path/filepath"
    "flag"
)

// define local variables
var zipFile string
var passwd string
var target string

// initialize variables
func init() {
	flag.StringVar(&zipFile, "f", "", "zip file to unpack")
	flag.StringVar(&passwd, "p", "infected", "password to use")
	flag.StringVar(&target, "d", "./", "target dir to write output to")
}

func main() {
	flag.Parse()
	err := unzip(zipFile, target)
	if err != nil {
		panic(err)
	}
}

// unzip and decrypt the archive
func unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer reader.Close()

	if err := os.MkdirAll(target, 0765); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}
		file.SetPassword(passwd)
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}
