package main

import (
	"os"
	"io"
	"log"
	"io/ioutil"
	"bytes"
	"github.com/alexmullins/zip"
)

var zipFile string
var files  []string
var passwd  string

func main() {
	log.Printf("os.Args: '%v'\n", os.Args)
	if len(os.Args) < 3 {
		log.Fatal("USAGE: zip archive.zip infected ./foo.exe ./bar.pdf ./baz.txt")
	}
	zipFile = os.Args[1]
	if zipFile == "" {
		log.Fatal("No zipFile given, first argument must be a a zipFile to create.")
	}
	passwd = os.Args[2]
	if passwd == "" {
		//log.Fatal("No password given, second argument need to be a password to be used.")
		passwd = "infected"
	}
	files = os.Args[3:]
	if len(files)  < 1 {
		log.Fatal("No files to archive given, third and following arguments must be filenames to be zipped/archived.")
	}
	// write a password zip
	raw := new(bytes.Buffer)
	zipw := zip.NewWriter(raw)
	for _, f := range files {
		contents, err := ioutil.ReadFile(f)
		if err != nil {
			log.Printf("Could not read file '%s'. Error: %s\n", f, err.Error)
			continue
		}
		w, err := zipw.Encrypt(f, passwd)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(w, bytes.NewReader(contents))
		if err != nil {
			log.Fatal(err)
		}
	}
	zipw.Close()
	err := ioutil .WriteFile(zipFile, raw.Bytes(), 0644)
	if err != nil {
		log.Fatal("Could not write zipFile '%s'. Erorr: %s\n", zipFile, err.Error)
	}
}
