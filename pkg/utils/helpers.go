package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func IsNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func WrapText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}

	var result string
	words := strings.Fields(text)
	line := ""

	for _, word := range words {
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

	if len(line) > 0 {
		if len(result) > 0 {
			result += "\n"
		}
		result += line
	}

	return result
}

func GetPriorityString(priority int) string {
	switch priority {
	case 1:
		return "High"
	case 2:
		return "Medium"
	case 3:
		return "Low"
	default:
		return "None"
	}
}

func FormatDate(t *time.Time) string {
	if t == nil {
		return "None"
	}
	return t.Format("2006-01-02 15:04")
}

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
