package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_Main_ArgsValidation tests the argument inputs
func Test_Main_ArgsValidation(t *testing.T) {
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Empty
	t.Run("Empty", func(t *testing.T) {
		res := argsValidation([]string{})
		assert.False(t, res)
	})

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Valid
	t.Run("Valid", func(t *testing.T) {
		res := argsValidation([]string{"*/15", "0", "1,15", "*", "1-5", "/usr/bin/find"})
		assert.True(t, res)
	})
}
