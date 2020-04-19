package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func main() {

	fileOldName := "rotate.log"
	fileRename := "rotate1.log"
	file, _ := os.OpenFile(fileOldName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0644))

	fmt.Println("open file", fileOldName)
	for i := 0; i < 5; i++ {
		file.WriteString(fmt.Sprintf("%v line: %d\n", time.Now(), i))
	}

	var statOldFile syscall.Stat_t
	if err := syscall.Stat(fileOldName, &statOldFile); err != nil {
		panic(err)
	}

	fmt.Println(fileOldName, "statOldFile.Ino:", statOldFile.Ino)

	err := os.Rename(fileOldName, fileRename)
	if err != nil {
		fmt.Println("file rename Error:", err)
	}

	fmt.Println("rename ", fileOldName, "->", fileRename)

	var statRenamedFile syscall.Stat_t
	if err := syscall.Stat(fileRename, &statRenamedFile); err != nil {
		panic(err)
	}

	fmt.Println("fileRename", "statRenamedFile.Ino:", statRenamedFile.Ino)

	fmt.Println("file.Name:", file.Name())

	for i := 5; i < 10; i++ {
		file.WriteString(fmt.Sprintf("%v line: %d\n", time.Now(), i))
	}

	fileOld, err := os.OpenFile(fileOldName, os.O_WRONLY|os.O_APPEND, os.FileMode(0644))

	if nil != err {
		fmt.Println(fileOldName, " open error:", err)
	} else {
		fmt.Println("fileOld.Name:", fileOld.Name())
	}

}
