package main

import (
	//"os"
	//"io"
	"bytes"
	"log"
	"github.com/alexmullins/zip"
	"flag"
	"io/ioutil"
)

var zipFile string
var passwd string

func init() {
	flag.StringVar(&zipFile, "f", "", "zip file to unpack")
	flag.StringVar(&passwd, "p", "infected", "password to use")
}

func main() {
	flag.Parse()
	z, err := ioutil.ReadFile(zipFile)
	if err != nil {
		log.Fatal("Could not read zipFile.")
	}
	// read the password zip
	zipr, err := zip.NewReader(bytes.NewReader(z), int64(len(z)))
	if err != nil {
		log.Fatal(err)
	}
	for _, z := range zipr.File {
		log.Printf("file: %s, flags: %v, size: %d bytes\n", z.Name, z.Flags, z.UncompressedSize64)
		z.SetPassword(passwd)
		rr, err := z.Open()
		if err != nil {
			log.Fatal(err)
		}
		contents, err := ioutil.ReadAll(rr)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(z.Name, contents, 0644)
		/*
		_, err = io.Copy(os.Stdout, rr)
		*/
		if err != nil {
			log.Fatal(err)
		}
		rr.Close()
	}

}
