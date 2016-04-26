package deb

import (
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestControl(t *testing.T) {
	info := Control{}
	info.Package = "package"
	info.Maintainer = "lala@hoho.com"
	assert.EqualValues(t, "Package: package\nMaintainer: lala@hoho.com\n", string(info.bytes()))
}

func TestEnsure(t *testing.T) {
	c, err := newCanonical()
	assert.NoError(t, err)
	defer c.close()
	defer os.Remove(c.file.Name())
	err = c.AddFile("control_test.go", "/var/lib/control/test/control.test")
	assert.NoError(t, err)
	assert.True(t, c.emptyFolders["/var/lib"])
	t.Log(c.emptyFolders)

}
