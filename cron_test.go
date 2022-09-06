package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_Cron_ToString(t *testing.T) {
	c := &Cron{
		Minute:     []int{0, 15, 30, 45},
		Hour:       []int{17, 18, 21},
		DayOfMonth: []int{25, 26, 27},
		Month:      []int{0, 1, 2, 3},
		DayOfWeek:  []int{7},
		Command:    "/usr/bin/find",
	}

	expected := `Minute: 0 15 30 45
Hour: 17 18 21
Day of Month: 25 26 27
Month: 0 1 2 3
Day of Week: 7
Command: /usr/bin/find`

	assert.Equal(t, expected, c.String())
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_Cron_IntSlice(t *testing.T) {
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Empty
	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t, IntSlice{}.String(), "")
	})

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Basic
	t.Run("Basic", func(t *testing.T) {
		assert.Equal(t, IntSlice{1, 2, 3}.String(), "1 2 3")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_Cron_PrintTable(t *testing.T) {
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Empty
	t.Run("Empty", func(t *testing.T) {
		c := &Cron{}
		out := CaptureOutput(c.PrintTable)

		expected := `minute        
hour          
day of month  
month         
day of week   
command       
`

		assert.Equal(t, expected, out)
	})

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Partially populated
	t.Run("Partially Populated", func(t *testing.T) {
		c := &Cron{Minute: IntSlice{1, 2, 3}, Month: IntSlice{0, 11}, Command: "/this/is/a/test"}
		out := CaptureOutput(c.PrintTable)

		expected := `minute        1 2 3
hour          
day of month  
month         0 11
day of week   
command       /this/is/a/test
`

		assert.Equal(t, expected, out)
	})

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Fully populated
	t.Run("Fully Populated", func(t *testing.T) {
		c := &Cron{
			Minute:     []int{0, 15, 30, 45},
			Hour:       []int{17, 18, 21},
			DayOfMonth: []int{25, 26, 27},
			Month:      []int{0, 1, 2, 3},
			DayOfWeek:  []int{7},
			Command:    "/this/is/a/test",
		}
		out := CaptureOutput(c.PrintTable)

		expected := `minute        0 15 30 45
hour          17 18 21
day of month  25 26 27
month         0 1 2 3
day of week   7
command       /this/is/a/test
`

		assert.Equal(t, expected, out)
	})
}
