package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	err := sampleRpm()
	assert.NoError(t, err)
	err = sampleDeb()
	assert.NoError(t, err)
}
