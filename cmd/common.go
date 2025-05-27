package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
	"errors"
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

func list_users() ([]userInfo, error, int) {
	users := make([]userInfo, 1000)
	size := 0
	f, e := os.Open("/etc/passwd")
	if e != nil { 
		return users, errors.New("Cannot read the information about the user. Please run the command as the superuser."), size
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		u := strings.Split(sc.Text(), ":")
		if len(u) > 0 {
			users[size].userName = u[0]
			users[size].userID = u[2]
			users[size].groupID = u[3]
			users[size].description = u[4]
			users[size].shell = u[6]
			size += 1
		}
	}

	if err := sc.Err(); err != nil {
		return users, errors.New("Error reading the /etc/passwd file."), size
	}
	return users, nil, size
}

func getUsersByGroup (group string) (string, error) {
	groupFound := false
	f, e := os.Open("/etc/group")
	if e != nil { 
		return "", errors.New("Cannot get information about the group")
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		g := strings.Split(sc.Text(), ":")
		if len(g) > 0 && g[0] == group {
            groupFound = true
			users := strings.TrimSpace(g[3])
			return users, nil
		}
		if groupFound == true { break}
	}

	if err := sc.Err(); err != nil {
		return "", nil
	}
	return "", nil
}

func arrContains(source []string, search string) bool {
	if search == "" {return false}
	for _, v := range source {
		if v == "" {
			continue
		}
		if v == search {
			return true
		}
	}
	return false
}


func getLastLogin(username string) (time.Time, error) {
	e := "Cannot read last login information"
	cmd := exec.Command("last", username, "-t", "YYYY-MM-DD hh:mm:ss")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return time.Time{}, fmt.Errorf(e)
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return time.Time{}, fmt.Errorf("No login record found")
	}
	fields := strings.Fields(lines[0])
	if len(fields) < 5 {
		return time.Time{}, fmt.Errorf(e)
	}
	timeStr := strings.Join(fields[3:6], " ")
    // Define the layout to match the output of the last command
	layout := "Mon Jan 2 15:04:05 2006 MST"
	lastLoginTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf(e)
	}
    // Check if the last login time is the default "never logged in" time
    if lastLoginTime.Year() == 1970 && lastLoginTime.Month() == time.January && lastLoginTime.Day() == 1 {
        return time.Time{}, fmt.Errorf("User has never logged in")
    }
	return lastLoginTime, nil
}


func getCurrentUser() (string, error) {
	currentUser, e := user.Current()
	if e != nil { return "", errors.New("Cannot get the username")}
	return currentUser.Username, nil
}