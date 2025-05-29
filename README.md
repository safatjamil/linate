# Linate
Linate, a linux associate , is a CLI tool for Linux systems. It's a unified tool that helps you to do many things 
from a single place. Linux users run commands and may open many files to do a single task. With linate you will 
run simple commands to perform complex tasks. Right now linate comes with 3 commands and associated sub commands.
Hold tight, more commands are coming.
> [!NOTE]
> For bug reporting, feedback, and command suggestions please send an email to safaetxaamil@gmail.com .

## Commands
## 1) bk
**1.1) bk take**
<br/>Take backup of a file. The backup filename will be \<oldFilename>-\<year>\<month>\<day>-\<count>.<br />
The backup file will be created in the same directory as the original file.<br/>
**Flags**
> --dir&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;directory of the original file that needs to be backed up<br/>
> --file&nbsp;&nbsp;&nbsp;&nbsp;name of the file
![Alt text](img/bk_take.png)

**1.2) bk check**
<br />Check backup files from the newest to the oldest.<br/>
**Flags**
> --dir&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;directory of the original file that needs to be backed up<br/>
> --file&nbsp;&nbsp;&nbsp;&nbsp;name of the file
![Alt text](img/bk_check.png)

**1.3) bk delete**
<br />Delete backup files. The oldest one will be deleted first. A yes/no promt will be shown for confirmation<br />
**Flags**
> --dir&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;directory of the original file that needs to be backed up<br/>
> --file&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;name of the file<br/>
> --number&nbsp;&nbsp;number of backup to delete. The default is 1
![Alt text](img/bk_delete.png)
