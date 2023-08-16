// mwextract - malware extract - extracts zip files encrypted with the password 'infected',
// which is basically industry standard in the anti-malware industry.
// mwextract does also unpack normal zip archives that are not protected by a password.
//
// mwextract does show checksums of files extracted
//
// Author: Florian 'scusi' Walther <flw@posteo.de>
//
// Usage:
//
// unpack archive 'archive.zip', use password 'infected815'
// mwextract -f archive.zip -p infected0815
//
// unpack archies '1.zip', '2.zip' and '3.zip'
// mwextract 1.zip 2.zip 3.zip
//
package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	aesZip "github.com/alexmullins/zip" // hacked zip lib, supports AES encrypted archives
	"github.com/scusi/MultiChecksum"
	"gitlab.scusi.io/flow/aeszip/magic"
	"io"
	"log"
	"os"
	"path/filepath"
)

// define local variables
var zipFile string // zipFile to unpack
var passwd string  // password to be used for decryption of zipFile
var target string  // file/path where to write output to
var verbose bool   // be verbose or not
var force bool		// force AES unzip

// init - initialize variables, set defaults
func init() {
	flag.StringVar(&zipFile, "f", "", "zip file to unpack")
	flag.StringVar(&passwd, "p", "infected", "password to use")
	flag.StringVar(&target, "d", "./", "target dir to write output to")
	flag.BoolVar(&verbose, "v", false, "be verbose if true")
	flag.BoolVar(&force, "force", true, "tries to AES decrypt anyway")
}

// main - main program loop
func main() {
	flag.Parse()
	if verbose == true {
		log.Printf("verbose == %t\n", verbose)
	}
	// if zipFile variale is empty
	if zipFile == "" {
		// range over arguments
		for _, zipFile = range os.Args[1:] {
			if verbose == true {
				log.Printf("zipFile is: %s\n", zipFile)
			}
			// check if file has a signature that we can handle
			ok, typ, err := checkType(zipFile)
			if err != nil {
				panic(err)
			}
			if verbose == true {
				log.Printf("file OK? : %t\n", ok)
			}
			if force == true {
				ok = true
				typ = "aes"
			}
			// unzip and decrypt if it is AES encrypted
			if ok == true && typ == "aes" {
				if verbose == true {
					log.Printf("zipFile is of correct version\n")
				}
				err := aesUnzip(zipFile, target)
				if err != nil {
					panic(err)
				}
			}
			// just unzip if it is a plain zip archive
			if ok == true && typ == "plain" {
				if verbose == true {
					log.Printf("unzip as non AES zip archive\n")
				}
				err := unzip(zipFile, target)
				if err != nil {
					panic(err)
				}
			}
		}
	} else {
		// operate on given zipFile
		if verbose == true {
			log.Printf("zipFile is: %s\n", zipFile)
		}
		ok, typ, err := checkType(zipFile)
		if err != nil {
			panic(err)
		}
		if verbose == true {
			log.Printf("file OK? : %t\n", ok)
		}
		if force == true {
			ok = true
			typ = "aes"
		}
		// unzip and decrypt if it is AES encrypted
		if ok == true && typ == "aes" {
			if verbose == true {
				log.Printf("zipFile is of correct version\n")
			}
			err := aesUnzip(zipFile, target)
			if err != nil {
				panic(err)
			}
		}
		// just unzip if it is a plain zip archive
		if ok == true && typ == "plain" {
			if verbose == true {
				log.Printf("unzip as non AES zip archive\n")
			}
			err := unzip(zipFile, target)
			if err != nil {
				panic(err)
			}
		}
	}
}

// checkType - checks if given file is a zip file we can handle
func checkType(fileName string) (ok bool, typ string, err error) {
	ok = false
	reader, err := os.Open(fileName)
	if err != nil {
		return ok, typ, err
	}
	mime, _ := magic.MIMETypeFromReader(reader)
	if mime == "Zip Archive Version 5.1, AES Encrypted" {
		typ = "aes"
		if verbose == true {
			log.Printf("mime type is correct: '%s'\n", mime)
		}
		return true, typ, nil
	}

	if mime == "Zip Archive Version 5.1, AES Encrypted, UTF-8" {
		typ = "aes"
		if verbose == true {
			log.Printf("mime type is correct: '%s'\n", mime)
		}
		return true, typ, nil
	}
	if mime == "Zip Archive" {
		typ = "plain"
		if verbose == true {
			log.Printf("mime type is correct: '%s'\n", mime)
		}
		return true, typ, nil
	} else {
		if verbose == true {
			log.Printf("MIME NOT OK, MIMEType was: %s\n", mime)
		}
		return ok, typ, nil
	}
}

// aesUnzip - unzip and decrypt the archive
func aesUnzip(archive, target string) error {
	if verbose == true {
		log.Printf("going to read archive: %s\n", archive)
	}
	reader, err := aesZip.OpenReader(archive)
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
		// use a T-reader
		var buf bytes.Buffer
		teeReader := io.TeeReader(fileReader, &buf)
		// copy unziped data to target file
		l, err := io.Copy(targetFile, teeReader)
		if err != nil {
			return err
		}
		// calculate checksums for each file written
		chksums := multichecksum.CalcChecksums(path, buf.Bytes())
		fmt.Println("")
		fmt.Printf("Filename: '%s', Size: '%d'\n%s", path, l, chksums.Filter("MD5", "SHA1", "SHA2", "Blake2b"))
	}
	return nil
}

// unzip - just use normal unzip procedure (without any encryption)
func unzip(archive, target string) error {
	if verbose == true {
		log.Printf("going to read archive: %s\n", archive)
	}
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
		//file.SetPassword(passwd)
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

		// use a T-reader
		var buf bytes.Buffer
		teeReader := io.TeeReader(fileReader, &buf)
		// copy unziped data to target file
		l, err := io.Copy(targetFile, teeReader)
		if err != nil {
			return err
		}
		// calculate checksums for each file written
		chksums := multichecksum.CalcChecksums(path, buf.Bytes())
		fmt.Println("")
		//fmt.Printf("Filename: '%s', Size: '%d'\n%s", path, l, chksums.String())
		fmt.Printf("Filename: '%s', Size: '%d'\n%s", path, l, chksums.Filter("MD5", "SHA1", "SHA2", "Blake2b"))
	}
	return nil
}
