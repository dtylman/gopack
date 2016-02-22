package main

import (
	"os"

	"github.com/dtylman/gopack/deb"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	output, err := os.Create("lala.deb")
	check(err)
	defer output.Close()
	d, err := deb.New(output)
	check(err)
	d.Data.AddFile("deb/deb.go", "usr/local/dtylman/deb.go")
	err = d.Data.AddFile("main.go", "usr/local/dtylman/gopackers")
	check(err)
	err = d.Create()
	check(err)
}
