package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Standard cron format
//
// ┌───────────── minute (0 - 59)
// │ ┌───────────── hour (0 - 23)
// │ │ ┌───────────── day of the month (1 - 31)
// │ │ │ ┌───────────── month (1 - 12)
// │ │ │ │ ┌───────────── day of the week (0 - 6)
// │ │ │ │ │
// │ │ │ │ │
// │ │ │ │ │
// * * * * * <command>
//
// Note: The following Special characters
//       * , - \ are supported
//
//  * == always
//  , == separate items (ex 0,1,2)
//  - == range (ex 0-15)
//  / == step by (ex */15)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Default slices for various parts of a cron expression
var (
	defaultMinuteSlice = []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
		30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
		50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
	}

	defaultHourSlice = defaultMinuteSlice[:24]

	defaultDomSlice = defaultMinuteSlice[1:32]

	defaultMonthSlice         = defaultMinuteSlice[1:13]
	defaultMonthSliceReplacer = strings.NewReplacer(
		"JAN", "1", "FEB", "2", "MAR", "3", "APR", "4", "MAY", "5", "JUN", "6",
		"JUL", "7", "AUG", "8", "SEP", "9", "OCT", "10", "NOV", "11", "DEC", "12")

	defaultDowSlice = defaultMinuteSlice[0:7]
	//defaultDowSliceStr = []string{
	//	"SUN", "MON", "TUE", "WED", "THR", "FRI", "SAT",
	//}
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ParseExpression parses a cron expression and
// builds a Cron stuct
func ParseExpression(exp string) (*Cron, error) {
	// Split and validate number of parts
	parts := strings.Split(exp, " ")
	if len(parts) < 6 {
		return nil, fmt.Errorf("not enough parts in the cron expression")
	}

	// Minute
	minute, err := parseSegment(parts[0], defaultMinuteSlice)
	if err != nil {
		return nil, fmt.Errorf("parsing error - minute - %s", err)
	}

	// Hour
	hour, err := parseSegment(parts[1], defaultHourSlice)
	if err != nil {
		return nil, fmt.Errorf("parsing error - hour - %s", err)
	}

	// Day of Month
	dayOfMonth, err := parseSegment(parts[2], defaultDomSlice)
	if err != nil {
		return nil, fmt.Errorf("parsing error - day of month - %s", err)
	}

	// Month
	monthReplaced := defaultMonthSliceReplacer.Replace(parts[3])
	month, err := parseSegment(monthReplaced, defaultMonthSlice)
	if err != nil {
		return nil, fmt.Errorf("parsing error - month - %s", err)
	}

	// Day of Week
	dayOfWeek, err := parseSegment(parts[4], defaultMonthSlice)
	if err != nil {
		return nil, fmt.Errorf("parsing error - day of week - %s", err)
	}

	return &Cron{
		Original:   exp,
		Minute:     minute,
		Hour:       hour,
		DayOfMonth: dayOfMonth,
		Month:      month,
		DayOfWeek:  dayOfWeek,
		Command:    strings.Join(parts[5:], " "),
	}, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// parseSegment parses an individual segment of an expression,
// such as the minute or hour
func parseSegment(expr string, inputSlice IntSlice) (IntSlice, error) {
	var result IntSlice

	// Is empty
	if expr == "" {
		return result, fmt.Errorf("empty")
	}

	// Wildcard
	if expr == "*" {
		result = make(IntSlice, len(inputSlice))
		copy(result, inputSlice)
		return result, nil
	}

	// Regexes (extremely generic on purpose)
	singleRegex := regexp.MustCompile(`^(\d+)$`)
	stepRegex := regexp.MustCompile(`^((?:(?:\d+-\d+)|\*)\/\d+)$`)
	rangeRegex := regexp.MustCompile(`^(\d+-\d+)$`)

	// Loop over each separate item
	for _, exp := range strings.Split(expr, ",") {
		if exp == "" {
			continue
		}

		// Single number
		match := singleRegex.MatchString(exp)
		if match {
			num, err := strconv.Atoi(exp)

			if err == nil && num >= inputSlice[0] && num <= inputSlice[len(inputSlice)-1] {
				result = append(result, num)
			} else {
				return result, fmt.Errorf("invalid")
			}

			continue
		}

		// Step
		match = stepRegex.MatchString(exp)
		if match {
			if out, err := explodeStep(exp, inputSlice); err != nil {
				return result, err
			} else {
				result = append(result, out...)
			}

			continue
		}

		// Range
		match = rangeRegex.MatchString(exp)
		if match {
			if out, err := explodeRange(exp, inputSlice); err != nil {
				return result, err
			} else {
				result = append(result, out...)
			}

			continue
		}
	}

	// Unique and sort
	result = UniqueIntSlice(result)
	sort.Ints(result)

	return result, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// explodeStep parses a step expression (ex */15) and explodes it
// into a slice
func explodeStep(stepExp string, inputSlice IntSlice) (IntSlice, error) {
	var (
		result, workingSlice IntSlice
		step                 int
	)

	// Is empty
	if stepExp == "" {
		return result, fmt.Errorf("step - empty")
	}

	// Validate step expression
	match, err := regexp.MatchString(`^(\d+-\d+|\*)\/\d+$`, stepExp)
	if err != nil || !match {
		return result, fmt.Errorf("step - invalid")
	}

	if stepExp[0] == '*' {
		workingSlice = inputSlice

		// Get the step
		fmt.Sscanf(stepExp, "*/%d", &step)
	} else {
		start, end := -1, -1

		// Get the start, end, step
		fmt.Sscanf(stepExp, "%d-%d/%d", &start, &end, &step)

		var err error
		workingSlice, err = explodeRange(fmt.Sprintf("%d-%d", start, end), inputSlice)
		if err != nil {
			return result, fmt.Errorf("step - %s", err)
		}
	}

	// Step is too big
	if step > workingSlice[len(workingSlice)-1] {
		return result, fmt.Errorf("step - step is too big")
	}

	// Build result
	for i := 0; i < len(workingSlice); i += step {
		result = append(result, workingSlice[i])
	}

	return result, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// explodeRange parses a range expression (1-15) and explodes it
// into a slice
func explodeRange(rangeExp string, inputSlice IntSlice) (result IntSlice, err error) {
	start, end := -1, -1

	// Is empty
	if rangeExp == "" {
		return result, fmt.Errorf("range - empty")
	}

	// Get the start and end
	fmt.Sscanf(rangeExp, "%d-%d", &start, &end)

	// Falls outside range
	if end < start || start < 0 || start < inputSlice[0] || end > inputSlice[len(inputSlice)-1] {
		return result, fmt.Errorf("range - invalid")
	}

	// Build the result
	result = make(IntSlice, end-start+1)
	for i := range result {
		result[i] = start + i
	}

	return result, nil
}
