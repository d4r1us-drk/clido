package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Constants for wrapping text and priority values.
const (
	PriorityHigh   = 1
	PriorityMedium = 2
	PriorityLow    = 3
	PriorityNone   = 4
)

// ParseIntOrError tries to parse a string as an integer and returns an error if the parsing fails.
func ParseIntOrError(value string) (int, error) {
	return strconv.Atoi(value)
}

// WrapText wraps a given text to a specified maximum line length.
func WrapText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}

	var result string
	words := strings.Fields(text) // Split text into words
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

// GetPriorityString returns the string representation of a Priority value.
func GetPriorityString(priority int) string {
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

// ColoredPastDue returns a colored string depending on the due date.
func ColoredPastDue(dueDate *time.Time, completed bool) string {
	if dueDate == nil {
		return color.GreenString("no")
	}

	// Ensure the current time is in the local time zone
	now := time.Now()
	localLocation := now.Location()

	// Grab dueDate and interpret it as local time
	dueDateAsLocalTime := time.Date(
		dueDate.Year(),
		dueDate.Month(),
		dueDate.Day(),
		dueDate.Hour(),
		dueDate.Minute(),
		dueDate.Second(),
		dueDate.Nanosecond(),
		localLocation, // Use local timezone for interpretation
	)

	if now.After(dueDateAsLocalTime) {
		if completed {
			return color.GreenString("yes")
		}
		return color.RedString("yes")
	}

	return color.GreenString("no")
}

// FormatDate formats a time.Time object into a human-readable string in the format "YYYY-MM-DD HH:MM".
func FormatDate(t *time.Time) string {
	if t == nil {
		return "None"
	}
	return t.Format("2006-01-02 15:04")
}

// ParseDueDate parses a date string in the format "2006-01-02 15:04" and returns a pointer to time.Time.
func ParseDueDate(dueDateStr string) (*time.Time, error) {
	date, err := time.Parse("2006-01-02 15:04", dueDateStr)
	if err != nil {
		return nil, err
	}
	return &date, nil
}
