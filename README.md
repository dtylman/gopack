# gopack

[![Go](https://github.com/dtylman/gopack/actions/workflows/go.yml/badge.svg)](https://github.com/dtylman/gopack/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dtylman/gopack)](https://goreportcard.com/report/github.com/dtylman/gopack)

Easy create deb & rpm pacakges.


## Creating packages

Usage:
```sh
go get github.com/dtylman/gopack
```

### Using CLI

Create a configuration file, as follows:
```javascript
{
    "name": "testpackage",
    "version": "0",
    "revision": "1",
    "arch": "amd64", //use deb or rpm specific value
    "description": "test package for gopack", 
    "homepage": "https://github.com/dtylman/gopack/",
    "depends": "binutils",
    "section": "Utils",
    "maintainer": "dtylman@gmail.com",
    // add a list of folders to be copied into the deb/rpm,	
    "folders": {
        "." : "/lala", //will copy everything from the local folder to /lala/...
        "/usr/local/bobobo" : ""  // will just create /usr/local/bobobo on target
	},
    //list of files to be copied to target, soruce -> target
    "files": {      
        "main.go" : "/usr/local/main.go"  
	},
    //provide files to be used as pre-inst and post-inst scripts
    "scripts": {
        "pre_inst": "preinst.sh",
        "post_inst": "postinst.sh",
        "pre_uninst": "",
        "post_uninst": ""
    }
}
```

Save the above as `pkg.conf.json` and use the following command:

```bash
gopack -conf pkg.conf.json -deb -rpm -output /tmp
```

Usage:
```bash
Usage of gopack:
  -conf string
        config file name (default "pkg.config.json")
  -deb
        build deb package
  -output string
        output path
  -rpm
        build rpm pacakge
```

### From Sources
And then: 
```go
import	"github.com/dtylman/gopack/deb"
```

### Creating DEB:

```go

func sampleDeb() error {
	d, err := deb.New(pkgName, pkgVersion, pkgRevision, deb.AMD64)
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
	debFileName, err := d.Create("")
	fmt.Println("Created " + debFileName)
	return err
}
```

### Creating RPM:
 *Note:* must have `rpmbuild` installed

```go
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
 	rpmFileName, err := r.Create("")
 	fmt.Println("Created " + rpmFileName)
 	return err
 }
```


