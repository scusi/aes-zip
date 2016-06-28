package main

import (
    "github.com/alexmullins/zip"    // hacked zip lib, supports AES encrypted archives
    "io"
    "os"
    "path/filepath"
    "strings"
    "flag"
    "log"
    "fmt"
)

var zipFile string
var passwd string
var target string
var needBaseDir bool

var Usage = func() {
    // usage: compress [-p password] [--includeBaseDir=true] zipfile file1 [file2 dir1 dir2 file3 [...]]
    fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "%s [-p password] [--includeBaseDir=false] zipfile file [otherFile dir otherdir [...]]\n", os.Args[0])
    flag.PrintDefaults()
    fmt.Fprintln(os.Stderr, "")
    fmt.Fprintf(os.Stderr, "%s is written by scusi - https://github.com/scusi/")
}
func init() {
    flag.StringVar(&zipFile, "f", "", "zip file to create")
    flag.StringVar(&passwd, "p", "infected", "password to use")
    //flag.StringVar(&target, "d", "./", "file/dir to add to zip")
    flag.BoolVar(&needBaseDir, "includeBaseDir", true, "per default a base directory is included nto zip archive")
}

func main() {
    flag.Parse()

    filesToCompress := flag.Args()
    if len(filesToCompress) < 1 {
        log.Fatal("no file(s) to compress given. First argument - and following - must be files to add to archive")
    }

    err := zipit(zipFile, needBaseDir, filesToCompress)
    if err != nil {
        panic(err)
    }
}

//func zipit(source, target string, needBaseDir bool) error {
func zipit(target string, needBaseDir bool, sources []string) error {
    zipfile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer zipfile.Close()

    archive := zip.NewWriter(zipfile)
    defer archive.Close()

    for _, source := range sources {
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

    }
    return err
}

