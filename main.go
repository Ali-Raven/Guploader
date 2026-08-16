package main

import (
	"os"
)

func main() {
	var cPath string
	cPath, err := os.Getwd()
	if err != nil {
		panic("Err : path have problem !")
	}

	myfile := file{
		Name:      "myfile",
		Path:      cPath,
		Size:      2232322,
		Extension: ".txt",
	}

	// myfile.showingPath()

	Server(myfile)
}
