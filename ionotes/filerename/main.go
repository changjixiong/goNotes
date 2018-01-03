package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	file, _ := os.OpenFile("rotate.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0644))

	for i := 0; i < 5; i++ {
		file.WriteString(fmt.Sprintf("%v line: %d\n", time.Now(), i))
		time.Sleep(time.Second * 2)
	}

	err := os.Rename("rotate.log", "rotate1.log") //重命名 C:\\log\\2013.log 文件为install.txt
	if err != nil {
		fmt.Println("file rename Error:", err)
	}

	fmt.Println("file.Name:", file.Name())

	for i := 5; i < 10; i++ {
		file.WriteString(fmt.Sprintf("%v line: %d\n", time.Now(), i))
		time.Sleep(time.Second * 2)
	}

	fileOld, err := os.OpenFile("rotate.log", os.O_WRONLY|os.O_APPEND, os.FileMode(0644))

	if nil != err {
		fmt.Println("fileOld open error:", err)
	} else {
		fmt.Println("fileOld.Name:", fileOld.Name())
	}

}
