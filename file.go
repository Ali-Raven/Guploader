package main

import "fmt"

type file struct {
	Name      string
	Path      string
	Size      int
	Extension string
}

func (f file) showingPath() {
	fmt.Printf("file path is : %s%s%s%s", f.Path, "/", f.Name, f.Extension)
}
