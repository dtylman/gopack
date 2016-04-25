package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	err := sampleRpm()
	assert.NoError(t, err)
	err = sampleDeb()
	assert.NoError(t, err)
}
