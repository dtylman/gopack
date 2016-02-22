# gopack
Native implementation of deb in go

## Creating deb package

Usage:
```go
import	"github.com/dtylman/gopack/deb"
```

Creating a deb:

```go
func sampleDeb() error {
	output, err := os.Create("helloworld_1.0-1.deb")
	if err != nil {
		return err
	}
	defer output.Close()
	d, err := deb.New(output)
	if err != nil {
		return err
	}
	err = d.Data.AddFile("/bin/ls", "./usr/local/bin/helloworld")
	if err != nil {
		return err
	}
	d.Info.Package = "helloworld"
	d.Info.Version = "1.0-1"
	d.Info.Architecture = "amd64"
	d.Info.Maintainer = "Mickey Mouse <mickey@disney.com>"
	d.Info.Section = "base"
	d.Info.Homepage = "http://disney.org/"
	d.Info.Depends = "libc6 (>= 2.14), libgcrypt11 (>= 1.5.1), zlib1g (>= 1:1.1.4)"
	d.Info.Description = `Hello world
  Lorum ipsum
  Yada yada`
	return d.Create()
}
```

