package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_Cron_DefaultValues(t *testing.T) {

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Valid test cases
	testCases := []struct {
		name          string
		slice         IntSlice
		expectedSlice IntSlice
	}{
		{"Minute", defaultMinuteSlice, IntSlice{
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
			30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
			40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
			50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
		}},
		{"Hour", defaultHourSlice, IntSlice{
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22, 23,
		}},
		{"Day_of_Month", defaultDomSlice, IntSlice{
			1, 2, 3, 4, 5, 6, 7, 8, 9,
			10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
			30, 31,
		}},
		{"Month", defaultMonthSlice, IntSlice{
			1, 2, 3, 4, 5, 6, 7, 8, 9,
			10, 11, 12,
		}},
		{"Day_of_week", defaultDowSlice, IntSlice{
			0, 1, 2, 3, 4, 5, 6,
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.ElementsMatch(t, tc.expectedSlice, tc.slice)

		})
	}

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Valid replacer test cases
	monthReplacerTestCases := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{"Month", "JAN,FEB,MAR", "1,2,3"},
		{"Month", "JAN-DEC", "1-12"},
	}

	for _, tc := range monthReplacerTestCases {
		t.Run(tc.name, func(t *testing.T) {
			replaced := defaultMonthSliceReplacer.Replace(tc.input)

			assert.Equal(t, tc.expectedOutput, replaced)
		})
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_Cron_ParseExpression(t *testing.T) {
	errorTestCases := []struct {
		name          string
		inputString   string
		expectedError string
	}{
		// Minute
		{"Empty", "", "not enough parts in the cron expression"},
		{"Not_Enough_Part", "1 2 3 4", "not enough parts in the cron expression"},
		{"Invalid_Minute", "100 2 3 4 5 /command", "parsing error - minute - invalid"},
		{"Invalid_Hour", "1 100 3 4 5 /command", "parsing error - hour - invalid"},
		{"Invalid_DoM", "1 2 100 4 5 /command", "parsing error - day of month - invalid"},
		{"Invalid_Month", "1 2 3 100 5 /command", "parsing error - month - invalid"},
		{"Invalid_DoW", "1 2 3 4 100 /command", "parsing error - day of week - invalid"},
	}

	for _, tc := range errorTestCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := ParseExpression(tc.inputString)
			assert.NotNil(t, err)
			assert.EqualError(t, err, tc.expectedError)
			assert.Empty(t, res)
		})
	}

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Valid test cases
	validTestCases := []struct {
		name        string
		inputString string
		expected    *Cron
	}{
		{"1", "1 2 3 4 5 /command", &Cron{
			Original:   "1 2 3 4 5 /command",
			Minute:     IntSlice{1},
			Hour:       IntSlice{2},
			DayOfMonth: IntSlice{3},
			Month:      IntSlice{4},
			DayOfWeek:  IntSlice{5},
			Command:    "/command"}},
		{"2", "*/15 0 1,15 * 1-5 /usr/bin/find", &Cron{
			Original:   "*/15 0 1,15 * 1-5 /usr/bin/find",
			Minute:     IntSlice{0, 15, 30, 45},
			Hour:       IntSlice{0},
			DayOfMonth: IntSlice{1, 15},
			Month:      IntSlice{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			DayOfWeek:  IntSlice{1, 2, 3, 4, 5},
			Command:    "/usr/bin/find"}},
		{"3", "*/15 0 1,15 * 1-5 /usr/bin/find bob .", &Cron{
			Original:   "*/15 0 1,15 * 1-5 /usr/bin/find bob .",
			Minute:     IntSlice{0, 15, 30, 45},
			Hour:       IntSlice{0},
			DayOfMonth: IntSlice{1, 15},
			Month:      IntSlice{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			DayOfWeek:  IntSlice{1, 2, 3, 4, 5},
			Command:    "/usr/bin/find bob ."}},
	}

	for _, tc := range validTestCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := ParseExpression(tc.inputString)
			assert.Nil(t, err)
			assert.Equal(t, res, tc.expected)
		})
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_Cron_ParseSegment(t *testing.T) {
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Invalid test cases
	errorTestCases := []struct {
		name          string
		inputString   string
		inputSlice    IntSlice
		expectedError string
	}{
		// Minute
		{"Minute_Empty", "", defaultMinuteSlice, "empty"},
		{"Minute_Invalid_Number", "70", defaultMinuteSlice, "invalid"},
		{"Minute_Invalid_Range", "10-0", defaultMinuteSlice, "range - invalid"},
		{"Minute_Invalid_Step", "10-0/2", defaultMinuteSlice, "step - range - invalid"},
		// Hour
		{"Hour_Empty", "", defaultHourSlice, "empty"},
		{"Hour_Invalid_Number", "70", defaultHourSlice, "invalid"},
		{"Hour_Invalid_Range", "10-0", defaultHourSlice, "range - invalid"},
		{"Hour_Invalid_Step", "10-0/2", defaultHourSlice, "step - range - invalid"},
		// Day of Month
		{"DoM_Empty", "", defaultDomSlice, "empty"},
		{"DoM_Invalid_Number", "0", defaultDomSlice, "invalid"},
		{"DoM_Invalid_Range", "10-1", defaultDomSlice, "range - invalid"},
		{"DoM_Invalid_Step", "10-1/2", defaultDomSlice, "step - range - invalid"},
		// Month
		{"Month_Empty", "", defaultMonthSlice, "empty"},
		{"Month_Invalid_Number", "0", defaultMonthSlice, "invalid"},
		{"Month_Invalid_Range", "10-1", defaultMonthSlice, "range - invalid"},
		{"Month_Invalid_Step", "10-1/2", defaultMonthSlice, "step - range - invalid"},
		// Day of Week
		{"DoW_Empty", "", defaultDowSlice, "empty"},
		{"DoW_Invalid_Number", "8", defaultDowSlice, "invalid"},
		{"DoW_Invalid_Range", "10-1", defaultDowSlice, "range - invalid"},
		{"DoW_Invalid_Step", "10-1/2", defaultDowSlice, "step - range - invalid"},
	}

	for _, tc := range errorTestCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parseSegment(tc.inputString, tc.inputSlice)
			assert.NotNil(t, err)
			assert.EqualError(t, err, tc.expectedError)
			assert.Empty(t, res)
		})
	}

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Valid test cases
	validTestCases := []struct {
		name          string
		inputString   string
		inputSlice    IntSlice
		expectedSlice IntSlice
	}{
		// Minute
		{"Minute_Wildcard", "*", defaultMinuteSlice, IntSlice{
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
			30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
			40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
			50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
		}},
		{"Minute_Ranged_Step", "10-20/5", defaultMinuteSlice, IntSlice{10, 15, 20}},
		{"Minute_Range", "10-15", defaultMinuteSlice, IntSlice{10, 11, 12, 13, 14, 15}},
		{"Minute_Multiple", "1,10-15,20-26/2,59", defaultMinuteSlice, IntSlice{1, 10, 11, 12, 13, 14, 15, 20, 22, 24, 26, 59}},
		{"Minute_Trailing_Comma", "1,", defaultMinuteSlice, IntSlice{1}},
		// Hour
		{"Hour_Wildcard", "*", defaultHourSlice, IntSlice{
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
			14, 15, 16, 17, 18, 19, 20, 21, 22, 23,
		}},
		{"Hour_Ranged_Step", "10-23/5", defaultHourSlice, IntSlice{10, 15, 20}},
		{"Hour_Range", "10-15", defaultHourSlice, IntSlice{10, 11, 12, 13, 14, 15}},
		{"Hour_Multiple", "1,10-15,20-22/2,23", defaultHourSlice, IntSlice{1, 10, 11, 12, 13, 14, 15, 20, 22, 23}},
		{"Hour_Trailing_Comma", "0,", defaultHourSlice, IntSlice{0}},
		// Day of Month
		{"DoM_Wildcard", "*", defaultDomSlice, IntSlice{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
			14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
			25, 26, 27, 28, 29, 30, 31,
		}},
		{"DoM_Ranged_Step", "10-31/5", defaultDomSlice, IntSlice{10, 15, 20, 25, 30}},
		{"DoM_Range", "10-15", defaultDomSlice, IntSlice{10, 11, 12, 13, 14, 15}},
		{"DoM_Multiple", "1,10-15,20-22/2,30", defaultDomSlice, IntSlice{1, 10, 11, 12, 13, 14, 15, 20, 22, 30}},
		{"DoM_Trailing_Comma", "31,", defaultDomSlice, IntSlice{31}},
		// Month
		{"Month_Wildcard", "*", defaultMonthSlice, IntSlice{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12,
		}},
		{"Month_Ranged_Step", "5-12/3", defaultMonthSlice, IntSlice{5, 8, 11}},
		{"Month_Range", "9-12", defaultMonthSlice, IntSlice{9, 10, 11, 12}},
		{"Month_Multiple", "1,3-5,8-12/2", defaultMonthSlice, IntSlice{1, 3, 4, 5, 8, 10, 12}},
		{"Month_Trailing_Comma", "12,", defaultMonthSlice, IntSlice{12}},
		// Day of Week
		{"DoW_Wildcard", "*", defaultDowSlice, IntSlice{
			0, 1, 2, 3, 4, 5, 6,
		}},
		{"DoW_Ranged_Step", "0-6/2", defaultDowSlice, IntSlice{0, 2, 4, 6}},
		{"DoW_Range", "3-5", defaultDowSlice, IntSlice{3, 4, 5}},
		{"DoW_Multiple", "0,1-3,4-6/2", defaultDowSlice, IntSlice{0, 1, 2, 3, 4, 6}},
		{"DoW_Trailing_Comma", "0,", defaultDowSlice, IntSlice{0}},
	}

	for _, tc := range validTestCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parseSegment(tc.inputString, tc.inputSlice)
			assert.Nil(t, err)
			assert.ElementsMatch(t, res, tc.expectedSlice)
		})
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_Cron_ExplodeStep(t *testing.T) {
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Invalid test cases
	errorTestCases := []struct {
		name          string
		inputString   string
		inputSlice    IntSlice
		expectedError string
	}{
		{"Empty", "", defaultMinuteSlice, "step - empty"},
		{"Invalid Data", "test", defaultMinuteSlice, "step - invalid"},
		{"Step Too Big", "*/100", defaultMinuteSlice, "step - step is too big"},
		{"Invalid Range", "10-5/2", defaultMinuteSlice, "step - range - invalid"},
	}

	for _, tc := range errorTestCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := explodeStep(tc.inputString, tc.inputSlice)
			assert.NotNil(t, err)
			assert.EqualError(t, err, tc.expectedError)
			assert.Empty(t, res)
		})
	}

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Valid test cases
	validTestCases := []struct {
		name          string
		inputString   string
		inputSlice    IntSlice
		expectedSlice IntSlice
	}{
		{"Wildcard", "*/15", defaultMinuteSlice, IntSlice{0, 15, 30, 45}},
		{"Range", "10-25/5", defaultMinuteSlice, IntSlice{10, 15, 20, 25}},
	}

	for _, tc := range validTestCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := explodeStep(tc.inputString, tc.inputSlice)
			assert.Nil(t, err)
			assert.ElementsMatch(t, res, tc.expectedSlice)
		})
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_Cron_ExplodeRange(t *testing.T) {
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Invalid test cases
	errorTestCases := []struct {
		name          string
		inputString   string
		inputSlice    IntSlice
		expectedError string
	}{
		{"Empty", "", defaultMinuteSlice, "range - empty"},
		{"Invalid Data", "test", defaultMinuteSlice, "range - invalid"},
		{"Invalid Range", "10-5", defaultMinuteSlice, "range - invalid"},
	}

	for _, tc := range errorTestCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := explodeRange(tc.inputString, tc.inputSlice)
			assert.NotNil(t, err)
			assert.EqualError(t, err, tc.expectedError)
			assert.Empty(t, res)
		})
	}

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Valid test cases
	validTestCases := []struct {
		name          string
		inputString   string
		inputSlice    IntSlice
		expectedSlice IntSlice
	}{
		{"Range", "10-20", defaultMinuteSlice, IntSlice{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
	}

	for _, tc := range validTestCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := explodeRange(tc.inputString, tc.inputSlice)
			assert.Nil(t, err)
			assert.ElementsMatch(t, res, tc.expectedSlice)
		})
	}
}
