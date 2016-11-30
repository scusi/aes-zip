// malware extract - extracts zip files encrypted with the password 'infected',
// which is basically industry standard in the anti-malware industry.
//
package main

import (
	"flag"
	"fmt"
	"github.com/alexmullins/zip" // hacked zip lib, supports AES encrypted archives
	"github.com/scusi/MultiChecksum"
	"github.com/scusi/magic"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	if zipFile == "" {
		for _, zipFile = range os.Args[1:] {
			log.Printf("zipFile is: %s\n", zipFile)
			ok, err := checkType(zipFile)
			if err != nil {
				panic(err)
			}
			log.Printf("file OK? : %t\n", ok)
			if ok {
				log.Printf("zipFile is of correct version\n")
				err := unzip(zipFile, target)
				if err != nil {
					panic(err)
				}
			}
		}
	} else {
		log.Printf("zipFile is: %s\n", zipFile)
		ok, err := checkType(zipFile)
		if err != nil {
			panic(err)
		}
		log.Printf("file OK? : %t\n", ok)
		if ok {
			log.Printf("zipFile is of correct version\n")
			err := unzip(zipFile, target)
			if err != nil {
				panic(err)
			}
		}
	}
}

// check if it a zip file with the right version...
func checkType(fileName string) (ok bool, err error) {
	ok = false
	reader, err := os.Open(fileName)
	if err != nil {
		return ok, err
	}
	mime, _ := magic.MIMETypeFromReader(reader)
	if mime == "Zip Archive Version 5.1, AES Encrypted" {
		log.Printf("mime type is correct: '%s'\n", mime)
		return true, nil
	}

	if mime == "Zip Archive Version 5.1, AES Encrypted, UTF-8" {
		log.Printf("mime type is correct: '%s'\n", mime)
		return true, nil
	} else {
		log.Printf("MIME NOT OK, MIMEType was: %s\n", mime)
		return ok, nil
	}
}

// unzip and decrypt the archive
func unzip(archive, target string) error {
	log.Printf("going to read archive: %s\n", archive)
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer reader.Close()

	if err := os.MkdirAll(target, 0765); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "extracting %d files...\n", len(reader.File))
	//for i, file := range reader.File {
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

		l, err := io.Copy(targetFile, fileReader)
		if err != nil {
			return err
		}
		fileData, err := ioutil.ReadAll(fileReader)
		if err != nil {
			return err
		}
		chksums := multichecksum.CalcChecksums(path, fileData)
		fmt.Println("")
		fmt.Printf("Filename: '%s', Size: '%d'\n%s", path, l, chksums.String())
	}
	return nil
}
