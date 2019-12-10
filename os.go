package goutill

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type osUtil struct{}

var OS = osUtil{}

func (osUtil) ENV() string {
	return os.Getenv("ENV")
}

func (osUtil) Move(filename string, dstPath string) error {
	in, err := File.Open(filename)
	if err != nil {
		return fmt.Errorf("couldn't open src file: %s", err)
	}
	dst := dstPath + "/" + filename
	out, err := os.Create(dst)
	if out != nil {
		return fmt.Errorf("couldn't open dst file: %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("writing dst file failed: %s", err)
	}

	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf("removing src file failed: %s", err)
	}

	return nil
}

func (osUtil) PWD() string {
	ex, err := os.Executable()
	if err != nil {
		err.Error()
	}
	exPath := filepath.Dir(ex)

	return exPath
}

func (osUtil) Mkdir(path string) error {
	permit := os.FileMode(0777)

	err := os.MkdirAll(path, permit)
	if err != nil {
		return err
	}
	return nil
}
