package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
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
	source, e := os.Open(src)
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
	_, e := os.Stat(dir)
	if e != nil {
		return false
	}
	return true
}

func fileExists(file string) bool {
	_, e := os.Stat(file)
	if e != nil {
		return false
	}
	return true
}

func getFileOwner(fPath string) (string, error) {
	fInfo, e := os.Stat(fPath)
	if e != nil {
		return "", e
	}

	stat, ok := fInfo.Sys().(*syscall.Stat_t)
	if ok != true {
		return "", fmt.Errorf("unexpected file info type")
	}

	uid := int(stat.Uid)
	usr, err := user.LookupId(fmt.Sprint(uid))
	if err != nil {
		return "", err
	}
	return usr.Username, nil
}

func yesNoPrompt(label string, def bool) bool {
	choices := "Y/n"
	if def == false {
		choices = "y/N"
	}
	r := bufio.NewReader(os.Stdin)
	var s string
	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

func exitWithError(errorText string) {
	fmt.Printf(errorText)
	os.Exit(1)
}

func timeDifference(time1 time.Time, time2 time.Time) {
	os.Exit(0)
}
