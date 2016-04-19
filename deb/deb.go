package deb

import (
	"errors"
	"time"

	"github.com/blakesmith/ar"
	"os"
	"path/filepath"
	"fmt"
)

const (
	debControl = "control.tar.gz"
	debData    = "data.tar.gz"
	debBinary  = "debian-binary"
)

//Deb represents a deb package
type Deb struct {
	Data    *canonical
	Control *canonical
	Info    Control
}

//New creates new deb writer
func New(name, version, revision, arch string) (*Deb, error) {
	deb := new(Deb)
	deb.Info.Package = name
	deb.Info.Version = version
	if revision!=""{
		deb.Info.Version+="-" + revision
	}
	deb.Info.Architecture =arch
	var err error
	deb.Data, err = newCanonical()
	if err != nil {
		return nil, err
	}
	deb.Control, err = newCanonical()
	if err != nil {
		return nil, err
	}
	return deb, nil
}

//Create creates the deb file
func (d *Deb) Create(folder string) error {
	if d.Info.Package == "" {
		return errors.New("Package name cannot be empty")
	}
	err := d.Control.AddEmptyFolder("./")
	if err != nil {
		return err
	}
	err = d.Control.AddBytes(d.Info.bytes(), "./control")
	if err != nil {
		return err
	}
	err = d.Control.AddBytes(d.Data.md5s.Bytes(), "./md5sums")
	if err != nil {
		return err
	}
	fileName:=filepath.Join(folder, fmt.Sprintf("%s_%s_%s.deb",d.Info.Package,d.Info.Version,d.Info.Architecture))
	debFile,err:=os.Create(fileName)
	if err != nil {
		return err
	}
	defer debFile.Close()
	ar := ar.NewWriter(debFile)
	err = ar.WriteGlobalHeader()
	if err != nil {
		return err
	}

	err = d.addBinary(ar)
	if err != nil {
		return err
	}
	err = d.Control.write(ar, debControl)
	if err != nil {
		return err
	}
	err = d.Data.write(ar, debData)
	if err != nil {
		return err
	}
	return nil
}

func (d *Deb) addBinary(writer *ar.Writer) error {
	body := []byte("2.0\n")
	header := new(ar.Header)
	header.Name = debBinary
	header.Mode = 0664
	header.Size = int64(len(body))
	header.ModTime = time.Now()
	err := writer.WriteHeader(header)
	if err != nil {
		return err
	}
	_, err = writer.Write(body)
	return err
}
