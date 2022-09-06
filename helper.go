package main

import (
	"bytes"
	"log"
	"os"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UniqueIntSlice removes duplicates from a int slice
func UniqueIntSlice(in IntSlice) IntSlice {
	keys := make(map[int]bool)
	list := IntSlice{}

	for _, entry := range in {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func CaptureOutput(f func()) string {
	setLogFlags()

	// Create a byte buffer a redirect log
	var buf bytes.Buffer
	log.SetOutput(&buf)

	f()

	// Reset logging
	log.SetOutput(os.Stdout)

	return buf.String()
}
