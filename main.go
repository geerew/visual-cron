package main

import (
	"log"
	"os"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func main() {
	setLogFlags()

	args := os.Args[1:]

	// Validate the args
	if !argsValidation(args) {
		log.Print("error - invalid input")
		os.Exit(1)
	}

	// Parse the expression
	cron, err := ParseExpression(args[0])
	if err != nil {
		log.Printf("error - %s", err.Error())
		os.Exit(1)
	}

	// Print as table
	cron.PrintTable()
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// setLogFlags sets logging flags
func setLogFlags() {
	// Remove timestamp from logging
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// argsValidation checks if the arguments are correct
//
// This currently just insures there is 1 arg but could
// be expanded to support more args
func argsValidation(args []string) bool {
	return len(args) >= 1
}
