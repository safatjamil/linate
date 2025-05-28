package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func init() {
	backUpCmd.AddCommand(takeBackupCmd)
	backUpCmd.AddCommand(checkBackupCmd)
	backUpCmd.AddCommand(deleteBackupCmd)
	takeBackupCmd.Flags().StringP("dir", "d", "", "Enter the directory of the file, absolute directory not relative.")
	takeBackupCmd.Flags().StringP("file", "f", "", "Enter the filename")
	takeBackupCmd.MarkFlagRequired("dir")
	takeBackupCmd.MarkFlagRequired("file")
	checkBackupCmd.Flags().StringP("dir", "d", "", "Enter the directory of the file, absolute directory not relative.")
	checkBackupCmd.Flags().StringP("file", "f", "", "Enter the filename")
	checkBackupCmd.MarkFlagRequired("dir")
	checkBackupCmd.MarkFlagRequired("file")
	deleteBackupCmd.Flags().StringP("dir", "d", "", "Enter the directory of the file, absolute directory not relative.")
	deleteBackupCmd.Flags().StringP("file", "f", "", "Enter the filename")
	deleteBackupCmd.Flags().IntP("number", "n", 1, "How many backups you want to delete. The oldest one will be deleted first.")
	deleteBackupCmd.MarkFlagRequired("dir")
	deleteBackupCmd.MarkFlagRequired("file")
}

var backUpCmd = &cobra.Command{
	Use:   "bk",
	Short: "Take, check and delete backup files.",
	Long:  `Take, check and delete backup files. Run linate bk --help for more options.`,
}

var takeBackupCmd = &cobra.Command{
	Use:   "take",
	Short: "Take backup.",
	Long:  `Take backup. Backup filename will be <filename>-<year><month><day>-<serialnumber> in the same directory.`,
	Run:   take_backup,
}

var checkBackupCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the last backup files.",
	Long:  `Check the last backup files.`,
	Run:   check_backup,
}

var deleteBackupCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the backup files.",
	Long:  `Delete the backup files. The oldest one will be deleted first.`,
	Run:   delete_backup,
}

type FileInfo struct {
	Name    string
	Size    string
	ModTime string
	Owner   string
}


func take_backup(cmd *cobra.Command, args []string) {
	var filePath string
	year, month, day := time.Now().Date()
	newFileName := ""
	var fn string
	var e error
	dir, _ := cmd.Flags().GetString("dir")
	file, _ := cmd.Flags().GetString("file")

	d := dirExists(dir)
	if d == false {
		exitWithError(fmt.Sprintf("Directory '%s' does not exist or follows a strict permission. Please run as the superuser if the directory really exists.\n", dir))
	}

	// Add forward slash after dir
	if dir[len(dir)-1:] == "/" {
		filePath = dir + file
	} else {
		filePath = dir + "/" + file
		dir = dir + "/"
	}

	f := fileExists(filePath)
	if f == false {
		exitWithError(fmt.Sprintf("File '%s' does not exist or follows a strict permission. Please run as the superuser if the file really exists.\n", file))
	}

	// Choose a filename
	for i := 1; i < 100; i++ {
		fn = fmt.Sprintf("%s-%d%s%d-%d", file, year, month, day, i)
		_, e = os.Stat(dir + fn)
		if e != nil {
			newFileName = fn
			break
		}
	}
	if newFileName == "" {
		exitWithError("It seems like there are already 99 backups.\n")
	}
    
	// Copy the old file to the backup file
	e = copyFile(dir+file, dir+newFileName)
	if e != nil {
		exitWithError(fmt.Sprintf("%sCan not create the backup file. Please run as the superuser if your user does not have permission to create a file in this directory.%s\n", colors["red"], colors["reset"]))
	}
	fmt.Printf("%slinate successfully created a backup file '%s'%s\n", colors["green"], newFileName, colors["reset"])
}


func check_backup(cmd *cobra.Command, args []string) {
	var e error
	dir, _ := cmd.Flags().GetString("dir")
	fileName, _ := cmd.Flags().GetString("file")
	d := dirExists(dir)
	if d == false {
		exitWithError(fmt.Sprintf("Directory '%s' does not exist or follows a strict permission. Please run as the superuser if the directory really exists.\n", dir))
	}

	// Add forward slash after dir
	if dir[len(dir)-1:] != "/" {
		dir = dir + "/"
	}
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		fmt.Println(e)
		return
	}

	// Sort by date modified
	sort.Slice(files, func(i, j int) bool {
		return files[j].ModTime().Before(files[i].ModTime())
	})

	var backups = make([]FileInfo, len(files))
	counter := 0
	var fn []string
	var tm time.Time
	for _, file := range files {
		// If file is a directory ignore
		if file.IsDir() == true {
			continue
		}
		// Check the filename
		fn = strings.Split(fmt.Sprintf("%s", file.Name()), "-")
		if len(fn) == 3 && fn[0] == fileName {
			tm = file.ModTime()
			backups[counter].Name = file.Name()
			backups[counter].Size = fmt.Sprintf("%v byte", file.Size())
			backups[counter].ModTime = fmt.Sprintf("%v-%v-%v | %v:%v", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute())
			ow, e := getFileOwner(dir + file.Name())
			if e == nil {
				backups[counter].Owner = ow
			} else {
				backups[counter].Owner = "ERROR"
			}
			counter += 1
		}
	}
    
	// Show 100 backups at most
	viewLength := 100
	if counter < viewLength {
		viewLength = counter
	}

	if viewLength == 0 {
		fmt.Printf("No backup found\n")
		os.Exit(0)
	}
	fmt.Printf("Total number of backups:%s %d%s\n", colors["yellow"], counter, colors["reset"])
    
	// Create the table view
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("File Name", "Size", "Date | Time", "Owner")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for i := 0; i < viewLength; i++ {
		tbl.AddRow(backups[i].Name, backups[i].Size, backups[i].ModTime, backups[i].Owner)
	}
	tbl.Print()
}


func delete_backup(cmd *cobra.Command, args []string) {
	var e error
	dir, _ := cmd.Flags().GetString("dir")
	fileName, _ := cmd.Flags().GetString("file")
	number, _ := cmd.Flags().GetInt("number")
	d := dirExists(dir)
	if d == false {
		exitWithError(fmt.Sprintf("Directory '%s' does not exist or follows a strict permission. Please run as the superuser if the directory really exists.\n", dir))
	}
	if dir[len(dir)-1:] != "/" {
		dir = dir + "/"
	}

	files, e := ioutil.ReadDir(dir)
	if e != nil {
		fmt.Println(e)
		return
	}

	// Sort by date modified
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	var backups = make([]FileInfo, len(files))
	counter := 0
	var fn []string
	var tm time.Time
	for _, file := range files {
		if file.IsDir() == true {
			continue
		}
		// Choose the files, oldest first
		fn = strings.Split(fmt.Sprintf("%s", file.Name()), "-")
		if len(fn) == 3 && fn[0] == fileName {
			tm = file.ModTime()
			backups[counter].Name = file.Name()
			backups[counter].Size = fmt.Sprintf("%v byte", file.Size())
			backups[counter].ModTime = fmt.Sprintf("%v-%v-%v | %v:%v", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute())
			ow, e := getFileOwner(dir + file.Name())
			if e == nil {
				backups[counter].Owner = ow
			} else {
				backups[counter].Owner = "ERROR"
			}
			counter += 1
		}
	}

	toDelete := number
	if number > counter {
		toDelete = counter
	}
	fmt.Printf("%sFollowing %d backup file(s) will be deleted%s\n\n", colors["red"], toDelete, colors["reset"])

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("File Name", "Size", "Date | Time", "Owner")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for i := 0; i < toDelete; i++ {
		tbl.AddRow(backups[i].Name, backups[i].Size, backups[i].ModTime, backups[i].Owner)
	}
	tbl.Print()
	fmt.Printf("\n")
	// Show a yes/no prompt
	ok := yesNoPrompt("Do you want to delete?", false)
	if ok == true {
		for i := 0; i < toDelete; i++ {
			e = os.Remove(dir + backups[i].Name)
			if e == nil {
				fmt.Printf("%sFile %s has been deleted successfully%s\n", colors["green"], backups[i].Name, colors["reset"])
			} else {
				fmt.Printf("%sFile %s could not be deleted. Check file permission or run as the super user.%s\n", colors["red"], backups[i].Name, colors["reset"])
			}
		}
	}
}
