package deb

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/blakesmith/ar"
)

type canonical struct {
	file      *os.File
	zip       *gzip.Writer
	md5       hash.Hash
	tarWriter *tar.Writer
}

func newCanonical() (*canonical, error) {
	c := new(canonical)
	var err error
	c.file, err = ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}
	c.zip = gzip.NewWriter(c.file)
	c.tarWriter = tar.NewWriter(c.zip)
	c.md5 = md5.New()
	return c, nil
}

func (c *canonical) addFile(name string, tarName string) error {
	fileInfo, err := os.Stat(name)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("%s is a directory, use AddFolder instead", name)
	}
	header, err := tar.FileInfoHeader(fileInfo, "")
	if tarName != "" {
		header.Name = tarName
	}
	if err != nil {
		return err
	}
	err = c.tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := io.TeeReader(file, c.md5)
	_, err = io.Copy(c.tarWriter, reader)
	if err != nil {
		return err
	}
	err = c.tarWriter.Flush()
	if err != nil {
		return err
	}
	err = c.zip.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (c *canonical) close() error {
	err := c.tarWriter.Flush()
	if err != nil {
		return err
	}
	err = c.tarWriter.Close()
	if err != nil {
		return err
	}
	err = c.zip.Close()
	if err != nil {
		return err
	}
	err = c.file.Close()
	if err != nil {
		return err
	}
	return nil
}

// cannot use io.copy because ar writes more data than wha'ts read
func ioCopy(writer io.Writer, reader io.Reader) error {
	buf := make([]byte, 64*1024)
	for {
		count, readErr := reader.Read(buf)
		if count > 0 {
			written, writeErr := writer.Write(buf[0:count])
			if writeErr != nil {
				return writeErr
			}
			if count > written {
				return io.ErrShortWrite
			}
		}
		if readErr == io.EOF {
			return nil
		}
		if readErr != nil {
			return readErr
		}
	}
}

func (c *canonical) write(writer *ar.Writer, name string) error {
	fileName := c.file.Name()
	defer os.Remove(fileName)

	err := c.close()
	if err != nil {
		return err
	}
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	in, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer in.Close()
	header := new(ar.Header)
	header.Name = name
	header.Size = fileInfo.Size()
	header.Mode = 0755
	header.ModTime = time.Now()
	err = writer.WriteHeader(header)
	if err != nil {
		return err
	}

	err = ioCopy(writer, in)
	if err != nil {
		return err
	}
	return nil
}
