package main

import (
	"fmt"
	"os"

	"github.com/dtylman/gopack/deb"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	fmt.Printf("hello world\n")
	output, err := os.Create("lala.deb")
	check(err)
	defer output.Close()
	d, err := deb.New(output)
	check(err)
	err = d.Data.AddFile("main.go", "usr/local/dtylman/gopackers")
	check(err)
	err = d.Create()
	check(err)
}
