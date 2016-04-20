# gopack

Easy create deb & rpm pacakges from `go`

## Creating packages

Usage:
```sh
go get github.com/dtylman/gopack/deb
```
And then: 
```go
import	"github.com/dtylman/gopack/deb"
```

### Creating deb:

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

### Creating rpm:
 *Note:* must have rpmbuild installed
 func sampleRpm() error {
 	r, err := rpm.New(pkgName, pkgVersion, pkgRevision, rpm.AMD64)
 	if err != nil {
 		return err
 	}
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
 	defer r.Close()
 	rpmFileName, err := r.Create("")
 	fmt.Println("Created " + rpmFileName)
 	return err
 }
```go

```


