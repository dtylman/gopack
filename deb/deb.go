package deb

import (
	"errors"
	"io"
	"time"

	"github.com/blakesmith/ar"
)

const (
	debControl = "control.tar.gz"
	debData    = "data.tar.gz"
	debBinary  = "debian-binary"
)

//Deb represents a deb package
type Deb struct {
	output  io.Writer
	Data    *canonical
	Control *canonical
	Info    Control
}

//New creates new deb writer
func New(writer io.Writer) (*Deb, error) {
	deb := new(Deb)
	deb.output = writer
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
func (d *Deb) Create() error {
	if d.Info.Package == "" {
		return errors.New("Package name cannot be empty")
	}
	err := d.Control.AddFolder("./")
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

	ar := ar.NewWriter(d.output)
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
