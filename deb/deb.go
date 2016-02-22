package deb

import (
	"io"
	"time"

	"fmt"

	"github.com/blakesmith/ar"
)

const (
	DebControl = "control.tar.gz"
	DebData    = "data.tar.gz"
	DebBinary  = "debian-binary"
)

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
	fmt.Println(d.Data.md5s.String())
	err := d.Control.AddBytes(d.Data.md5s.Bytes(), "md5sums")
	if err != nil {
		return err
	}
	ar := ar.NewWriter(d.output)
	err = ar.WriteGlobalHeader()
	if err != nil {
		return err
	}
	err = d.Data.write(ar, DebData)
	if err != nil {
		return err
	}
	err = d.Control.write(ar, DebControl)
	if err != nil {
		return err
	}
	err = d.addBinary(ar)
	if err != nil {
		return err
	}
	return nil
}

func (d *Deb) addBinary(writer *ar.Writer) error {
	body := []byte("2.0\n")
	header := new(ar.Header)
	header.Name = DebBinary
	header.Mode = 0644
	header.Size = int64(len(body))
	header.ModTime = time.Now()
	err := writer.WriteHeader(header)
	if err != nil {
		return err
	}
	_, err = writer.Write(body)
	return err
}
