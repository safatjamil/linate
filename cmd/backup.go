package cmd

import (
	"os"
	"fmt"
	_"io/ioutil"
	_"log"
	"time"
	"github.com/spf13/cobra"
)

func init() {
	backUpCmd.AddCommand(takeBackupCmd)
	takeBackupCmd.Flags().StringP("dir", "d", "", "Enter the directory of the file, absolute directory not relative.")
	takeBackupCmd.Flags().StringP("file", "f", "", "Enter the filename")
	takeBackupCmd.MarkFlagRequired("dir")
	takeBackupCmd.MarkFlagRequired("file")
	checkBackupCmd.Flags().StringP("dir", "d", "", "Enter the directory of the file, absolute directory not relative.")
	checkBackupCmd.Flags().StringP("file", "f", "", "Enter the filename")
	checkBackupCmd.MarkFlagRequired("dir")
	checkBackupCmd.MarkFlagRequired("file")
}

var backUpCmd = &cobra.Command{
	Use:   "bk",
	Short: "Take, check and delete backup files",
	Long:  `Take, check and delete backup files. Run linate bk --help for more options`,
}

var takeBackupCmd = &cobra.Command{
	Use:   "take",
	Short: "Take backup",
	Long:  `Take backup. Backup filename will be <filename>-<year><month><day>-<serialnumber> in the same directory`,
	Run:   take_backup,
}

var checkBackupCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the last backup files",
	Long:  `Check the last backup files`,
	Run:   check_backup,
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
		fmt.Printf("Directory '%s' does not exist or follows a strict permission. Please run as the superuser if the directory really exists.\n", dir)
		os.Exit(1)
	}
    
    if dir[len(dir)-1:] == "/"{
		filePath = dir + file
	} else {
		filePath = dir + "/" + file
		dir = dir + "/"
	}

    f := fileExists(filePath)
	if f == false {
		fmt.Printf("File '%s' does not exist or follows a strict permission. Please run as the superuser if the file really exists.\n", file)
		os.Exit(1)
	}
	for i:=1; i<100; i++ {
		fn = fmt.Sprintf("%s-%d%s%d-%d", file, year, month, day, i)
		_, e = os.Stat(dir + fn)
		if e!= nil {
            newFileName = fn
			break
		}
	}
	if newFileName == "" {
		fmt.Printf("It seems like there are already 99 backups.\n")
		os.Exit(1)
	}

	e = copyFile(dir+file, dir+newFileName)
	if e != nil {
		fmt.Printf("Can not create the backup file. Please run as the superuser if your user does not have permission to create a file in this directory.\n")
		os.Exit(1)
	}
    fmt.Printf("linate successfully created a backup file %s\n", newFileName)
}


func check_backup(cmd *cobra.Command, args []string) {
	var backups []string
	var e error

	dir, _ := cmd.Flags().GetString("dir")
	file, _ := cmd.Flags().GetString("file")

	d := dirExists(dir)
	if d == false {
		fmt.Printf("Directory '%s' does not exist or follows a strict permission. Please run as the superuser if the directory really exists.\n", dir)
		os.Exit(1)
	}
    
    
	
}
