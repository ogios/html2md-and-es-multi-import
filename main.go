package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/mattn/godown"
)

var BASE_DIR string

func GetFiles(path string) []fs.DirEntry {
	fmt.Println("Reading dir:", path)
	fs, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	return fs
}

func GetFile(path string) (*os.File, fs.FileInfo) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	ft, err := f.Stat()
	if err != nil {
		panic(err)
	}
	return f, ft
}

func CreateDir(path string) {
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}
}

func SaveFile(p string) *os.File {
	basename := path.Base(p)
	dir := path.Dir(p) + "/md/"
	base := strings.Split(basename, ".")[0]
	CreateDir(dir)
	f, err := os.OpenFile(dir+base+".md", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return f
}

func ParseOneFile(path string) {
	f, _ := GetFile(path)
	defer f.Close()
	c := &CustomWriter{}
	err := godown.Convert(c, f, nil)
	f.Close()
	if err != nil {
		panic(err)
	}
	s := SaveFile(path)
	defer s.Close()
	_, err = s.Write(c.Content)
	if err != nil {
		panic(err)
	}
}

func ParseFiles() {
	BASE_DIR := os.Args[1]
	if BASE_DIR[len(BASE_DIR)-1] == "/"[0] {
		BASE_DIR = BASE_DIR[:len(BASE_DIR)-1]
	}
	files := GetFiles(BASE_DIR)
	for _, file := range files {
		if !file.IsDir() {
			ParseOneFile(BASE_DIR + "/" + file.Name())
		}
	}
}

type CustomWriter struct {
	Content []byte
}

func (c *CustomWriter) Write(bs []byte) (int, error) {
	match, err := regexp.Match("!\\[\\]\\(data:image/.*?;base64,", bs)
	if err != nil {
		panic(err)
	}
	if !match {
		c.Content = append(c.Content, bs...)
	} else {
		fmt.Println("filter")
	}
	return len(bs), nil
}

func main() {
	// ParseFiles()
	InitCLI()
	// NewTableTest()
	// DeleteIndex()
	// CreateIndex()
	DBInit()
	Test()
}
