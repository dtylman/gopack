package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/dtylman/gopack/rpm"

	"github.com/dtylman/gopack/deb"

	"github.com/stretchr/testify/assert"
)

func sampleRpm() error {
	r, err := rpm.New("pkgName", "0", "1.1", rpm.AMD64)
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
	rpmFileName, err := r.Create("")
	fmt.Println("Created " + rpmFileName)
	return err
}

func sampleDeb() error {
	d, err := deb.New("pkgName", "0", "1.1", deb.AMD64)
	if err != nil {
		return err
	}
	err = d.Data.AddFile("/bin/ls", "/opt/danny/bin/ls")
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
	d.ConfFiles = "/etc/hoho.cfg"
	debFileName, err := d.Create("")
	fmt.Println("Created " + debFileName)
	return err
}

func TestCreate(t *testing.T) {
	err := sampleRpm()
	assert.NoError(t, err)
	err = sampleDeb()
	assert.NoError(t, err)
	defer func() {
		os.Remove("pkgName_0-1.1_amd64.deb")
		os.RemoveAll("x86_64")
	}()
	_, err = os.Stat("pkgName_0-1.1_amd64.deb")
	assert.NoError(t, err)
	_, err = os.Stat("x86_64/pkgName-0_1.1-1.x86_64.rpm")
	assert.NoError(t, err)
}
