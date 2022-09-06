package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Helper_UniqueIntSlice(t *testing.T) {
	res := UniqueIntSlice(IntSlice{1, 1, 2, 3, 3, 4, 5, 6})
	expected := IntSlice{1, 2, 3, 4, 5, 6}

	assert.ElementsMatch(t, res, expected)
}

func Test_Helper_CaptureLogs(t *testing.T) {
	res := CaptureOutput(func() {
		log.Print("test message")
	})

	assert.Equal(t, res, "test message\n")
}
