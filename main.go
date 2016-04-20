package main

import (
	"fmt"
	"os"

	"github.com/dtylman/gopack/deb"
	"github.com/dtylman/gopack/rpm"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const (
	pkgName     = "helloworld"
	pkgVersion  = "1.0"
	pkgRevision = "1"
)

func sampleDeb() error {
	d, err := deb.New(pkgName, pkgVersion, pkgRevision, "amd64")
	if err != nil {
		return err
	}
	err = d.Data.AddFile("/bin/ls", "./usr/local/bin/helloworld")
	if err != nil {
		return err
	}
	d.Info.Maintainer = "Mickey Mouse <mickey@disney.com>"
	d.Info.Section = "base"
	d.Info.Homepage = "http://disney.org/"
	d.Info.Depends = "libc6 (>= 2.14), libgcrypt11 (>= 1.5.1), zlib1g (>= 1:1.1.4)"
	d.Info.Description = `Hello world
  Lorum ipsum
  Yada yada`
	return d.Create("")
}

func sampleRpm() error {
	r, err := rpm.New(pkgName, pkgVersion, pkgRevision, "x86_64")
	if err != nil {
		return err
	}
	r.Spec.Header[rpm.Summary] = "Hello world app"
	r.Spec.Header[rpm.Packager] = "Mickey Mouse <mickey@disney.com>"
	r.Spec.Header[rpm.URL] = "http://disney.org/"

	r.Spec.Depends("yum", "rpm", "mc")
	r.Spec.Description = `Hello world
  Lorum ipsum
  Yada yada`
	err = r.AddFile("/bin/ls", "/opt/danny")
	if err != nil {
		return err
	}
	defer r.Close()
	return r.Create()
}

func main() {
	check(sampleRpm())
	check(sampleDeb())
}
