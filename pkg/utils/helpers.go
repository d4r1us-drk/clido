package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// IsNumeric checks whether a given string is numeric.
// It attempts to convert the string into an integer and returns true if successful, otherwise false.
//
// Example:
//   IsNumeric("123") // returns true
//   IsNumeric("abc") // returns false
func IsNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// WrapText wraps a given text to a specified maximum line length.
// It breaks the text into words and ensures that no line exceeds the specified maxLength.
//
// If the text fits within maxLength, it is returned unchanged. Otherwise, the text is split into multiple lines.
//
// Example:
//   WrapText("This is a very long sentence that needs to be wrapped.", 20)
//   // returns:
//   "This is a very long\nsentence that needs\nto be wrapped."
func WrapText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}

	var result string
	words := strings.Fields(text)  // Split text into words
	line := ""

	for _, word := range words {
		// Check if adding the word would exceed the max length
		if len(line)+len(word)+1 > maxLength {
			if len(result) > 0 {
				result += "\n"
			}
			result += line
			line = word
		} else {
			if len(line) > 0 {
				line += " "
			}
			line += word
		}
	}

	// Add the remaining line
	if len(line) > 0 {
		if len(result) > 0 {
			result += "\n"
		}
		result += line
	}

	return result
}

// Priority represents the priority level of a task.
// Priority levels include High (1), Medium (2), and Low (3).
type Priority int

const (
	PriorityHigh   Priority = 1  // High priority
	PriorityMedium Priority = 2  // Medium priority
	PriorityLow    Priority = 3  // Low priority
)

// GetPriorityString returns the string representation of a Priority value.
// The possible values are "High", "Medium", "Low", or "None" (for undefined priority).
//
// Example:
//   GetPriorityString(PriorityHigh) // returns "High"
//   GetPriorityString(4)            // returns "None"
func GetPriorityString(priority Priority) string {
	switch priority {
	case PriorityHigh:
		return "High"
	case PriorityMedium:
		return "Medium"
	case PriorityLow:
		return "Low"
	default:
		return "None"
	}
}

// FormatDate formats a time.Time object into a human-readable string in the format "YYYY-MM-DD HH:MM".
// If the time is nil, it returns "None".
//
// Example:
//   FormatDate(time.Now())   // returns "2024-09-11 14:30"
//   FormatDate(nil)          // returns "None"
func FormatDate(t *time.Time) string {
	if t == nil {
		return "None"
	}
	return t.Format("2006-01-02 15:04")
}

// ColoredPastDue determines if a task is past due and returns a colored string indicating the result.
// If the task is past due and not completed, it returns "yes" in red. If it's completed or not past due, it returns "no" in green.
//
// It converts the given due date to the local time zone for comparison.
// If no due date is provided, it assumes the task is not past due.
//
// Example:
//   ColoredPastDue(&time.Now(), false) // returns "no" in green if task is not past due
//   ColoredPastDue(&pastTime, false)   // returns "yes" in red if task is past due and incomplete
//   ColoredPastDue(&pastTime, true)    // returns "yes" in green if task is completed
func ColoredPastDue(dueDate *time.Time, completed bool) string {
	if dueDate == nil {
		return color.GreenString("no")
	}

	// Ensure the current time is in the local time zone
	now := time.Now()
	localLocation := now.Location()

	// Convert due date to local time zone
	dueDateAsLocalTime := time.Date(
		dueDate.Year(),
		dueDate.Month(),
		dueDate.Day(),
		dueDate.Hour(),
		dueDate.Minute(),
		dueDate.Second(),
		dueDate.Nanosecond(),
		localLocation,  // Use local timezone
	)

	// Compare current time with the due date
	if now.After(dueDateAsLocalTime) {
		if completed {
			return color.GreenString("yes")
		}
		return color.RedString("yes")
	}

	return color.GreenString("no")
}
