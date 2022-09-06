package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/gookit/goutil/arrutil"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Cron represents a single cron expression
type Cron struct {
	Original   string   `table:"-"`
	Minute     IntSlice `table:"minute"`
	Hour       IntSlice `table:"hour"`
	DayOfMonth IntSlice `table:"day of month"`
	Month      IntSlice `table:"month"`
	DayOfWeek  IntSlice `table:"day of week"`
	Command    string   `table:"command"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// IntSlice represents an int
type IntSlice []int

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// String returns a string representation of IntSlice
func (i IntSlice) String() string {
	// Ignore the error. In this case it will just
	// be an empty slice
	r, _ := arrutil.ToStrings(i)
	return strings.Join(r, " ")
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// String generates a stringified version of the struct
func (c Cron) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Minute: %s\n", c.Minute.String())
	fmt.Fprintf(&sb, "Hour: %s\n", c.Hour.String())
	fmt.Fprintf(&sb, "Day of Month: %s\n", c.DayOfMonth.String())
	fmt.Fprintf(&sb, "Month: %s\n", c.Month.String())
	fmt.Fprintf(&sb, "Day of Week: %s\n", c.DayOfWeek.String())
	fmt.Fprintf(&sb, "Command: %s", c.Command)

	return sb.String()
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// PrintTable outputs the struct in table format to stdout
func (c Cron) PrintTable() {
	tagging := "table"

	// String builder (for the tabwriter)
	var sb strings.Builder

	// Tabwriter
	w := new(tabwriter.Writer)
	w.Init(&sb, 0, 0, 2, ' ', 0)
	v := reflect.ValueOf(c)

	for i := 0; i < v.NumField(); i++ {
		// Field tag value
		tag := v.Type().Field(i).Tag.Get(tagging)

		// Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			continue
		}

		fmt.Fprintf(w, "%s\t%s\n", tag, v.Field(i).Interface())
	}

	// Flush
	w.Flush()

	log.Print(sb.String())
}
