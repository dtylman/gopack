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
	d, err := deb.New(pkgName, pkgVersion, pkgRevision, deb.AMD64)
	if err != nil {
		return err
	}
	d.PreInst = `echo hello world!`
	d.Info.Maintainer = "Mickey Mouse <mickey@disney.com>"
	d.Info.Section = "base"
	d.Info.Homepage = "http://disney.org/"
	d.Info.Depends = "libc6 (>= 2.14), libgcrypt11 (>= 1.5.1), zlib1g (>= 1:1.1.4)"
	d.Info.Description = `Hello world
  Lorum ipsum
  Yada yada`
	err = d.Data.AddFile("/bin/ls", "/opt/danny/bin/ls")
	if err != nil {
		return err
	}
	err = d.Data.AddEmptyFolder("/var/log/lala")
	if err != nil {
		return err
	}
	err = d.AddFolder("/tmp/", "/my/remote/")
	if err != nil {
		return err
	}
	debFileName, err := d.Create("")
	fmt.Println("Created " + debFileName)
	return err
}

func sampleRpm() error {
	r, err := rpm.New(pkgName, pkgVersion, pkgRevision, rpm.AMD64)
	if err != nil {
		return err
	}
	defer r.Close()
	r.Spec.Pre = `echo hello world!`
	r.Spec.Header[rpm.Summary] = "Hello world app"
	r.Spec.Header[rpm.Packager] = "Mickey Mouse <mickey@disney.com>"
	r.Spec.Header[rpm.URL] = "http://disney.org/"

	r.Spec.Depends("yum", "rpm", "mc")
	r.Spec.Description = `Hello world
  Lorum ipsum
  Yada yada`
	err = r.AddFile("/bin/ls", "/opt/danny/bin/ls")
	if err != nil {
		return err
	}
	err = r.AddEmptyFolder("/var/log/lala")
	if err != nil {
		return err
	}
	err = r.AddFolder("/tmp/", "/my/remote/")
	if err != nil {
		return err
	}
	rpmFileName, err := r.Create("")
	fmt.Println("Created " + rpmFileName)
	return err
}

func main() {
	check(sampleRpm())
	check(sampleDeb())
}
