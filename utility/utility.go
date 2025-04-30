package utility

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	figure "github.com/common-nighthawk/go-figure"
)

// UniqueItems returns the array of strings containing unique items
func UniqueItems(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Creates art for console output
func CreateArt(msg string) {
	myFigure := figure.NewFigure(msg, "basic", true)
	myFigure.Print()
}

// Copies Files and Directories from source to destination
func File(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func Dir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err

	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = Dir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = File(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func IsFileSymlink(path string) bool {
	fileMetadata, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return fileMetadata.Mode()&os.ModeSymlink != 0
}

func SafeReadFile(path string) ([]byte, error) {
	if IsFileSymlink(path) {
		log.Debug("Symlink was detected! Not following symlink ", path)
		return []byte{}, nil
	}
	return ioutil.ReadFile(path)
}
