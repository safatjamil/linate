package cmd

import (
	"io"
	"os"
	"strconv"
)

func toInt(raw string) int {
	if raw == "" {
		return 0
	}
	res, err := strconv.Atoi(raw)
	if err != nil {
		panic(err)
	}
	return res
}


func copyFile(src string, dst string) error {
	source, e:= os.Open(src)
	if e != nil {
		return e
	}
	defer source.Close()

	dest, e_ := os.Create(dst)
	if e_ != nil {
		return e_
	}
	defer dest.Close()
	_, er := io.Copy(dest, source)
	return er
}

func dirExists(dir string) bool {
	_, e := os.Stat(dir); 
	if e != nil {
		return false
	}
	return true
}

func fileExists(file string) bool {
	_, e := os.Stat(file); 
	if e != nil {
		return false
	}
	return true
}