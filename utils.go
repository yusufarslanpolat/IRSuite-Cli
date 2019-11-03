package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path"
)

//var BUFFERSIZE = 1000000
//func File(src, dst string) error {
//	sourceFileStat, err := os.Stat(src)
//	if err != nil {
//		return err
//	}
//
//	if !sourceFileStat.Mode().IsRegular() {
//		return fmt.Errorf("%s is not a regular file.", src)
//	}
//
//	source, err := os.Open(src)
//	if err != nil {
//		return err
//	}
//	defer source.Close()
//
//	_, err = os.Stat(dst)
//	if err == nil {
//		return fmt.Errorf("File %s already exists.", dst)
//	}
//
//	destination, err := os.Create(dst)
//	if err != nil {
//		return err
//	}
//	defer destination.Close()
//
//	if err != nil {
//		panic(err)
//	}
//
//	buf := make([]byte, BUFFERSIZE)
//	for {
//		n, err := source.Read(buf)
//		if err != nil && err != io.EOF {
//			return err
//		}
//		if n == 0 {
//			break
//		}
//
//		if _, err := destination.Write(buf[:n]); err != nil {
//			return err
//		}
//	}
//	return err
//}
func File(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return errors.Wrap(err, "Copy file Failed.")
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
		return errors.Wrap(err, "Read Directory Failed.")
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
func WriteToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}
