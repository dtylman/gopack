package deb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestControl(t *testing.T) {
	info := Control{}
	info.Package = "package"
	info.Maintainer = "lala@hoho.com"
	assert.EqualValues(t, "Package: package\nMaintainer: lala@hoho.com\n", string(info.bytes()))
}
