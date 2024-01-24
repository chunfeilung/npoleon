package util

import (
	"fmt"
	"testing"
	"time"
)

var testDataParseTime = []struct {
	input    string
	expected string
}{
	{"2023-12-20 22:55", "2023-12-20 22:55:00"},
	{"10/1/2024 13:37", "2024-01-10 13:37:00"},
	{"10/01/2024 13:37", "2024-01-10 13:37:00"},
	{"10/01/2024 13:37:42", "2024-01-10 13:37:42"},
	{"12-02-2024 8:45:00", "2024-02-12 08:45:00"},
	{"12-02-2024", "2024-02-12 00:00:00"},
	{"12-2-2024", "2024-02-12 00:00:00"},
	{"14/3/2024", "2024-03-14 00:00:00"},
	{"12 februari 2024", ""},
}

func TestParseTime(t *testing.T) {
	for _, data := range testDataParseTime {
		t.Run(fmt.Sprintf("input=%s", data.input), func(t *testing.T) {
			// > Arrange
			expected, _ := time.ParseInLocation("2006-01-02 15:04:05", data.expected, loc)

			// > Act
			res, err := ParseTime(data.input)

			// > Assert
			if data.expected == "" && err == nil {
				t.Errorf("Parsing '%v' should have failed", data.input)
			}
			if !res.Time.Equal(expected) {
				t.Errorf("Expected '%v', got '%v'", data.expected, res)
			}
		})
	}
}

var testDataParseTimeFrom = []struct {
	input    string
	expected string
}{
	{"2024-01-10 11:55", "2024-01-10 11:55:00"},
	{"2024-01-10", "2024-01-10 00:00:00"},
	{"11:55", "2024-01-10 11:55:00"},
	{"22:55", "2024-01-09 22:55:00"},
}

func TestParseTimeFrom(t *testing.T) {
	now = func() time.Time {
		newNow, _ := time.ParseInLocation("2006-01-02 15:04:05", "2024-01-10 17:08:23", loc)
		return newNow
	}

	for _, data := range testDataParseTimeFrom {
		t.Run(fmt.Sprintf("input=%s", data.input), func(t *testing.T) {
			// > Arrange
			expected, _ := time.ParseInLocation("2006-01-02 15:04:05", data.expected, loc)

			// > Act
			res, err := ParseTimeFrom(data.input)

			// > Assert
			if data.expected == "" && err == nil {
				t.Errorf("Parsing '%v' should have failed", data.input)
			}
			if !res.Equal(expected) {
				t.Errorf("Expected '%v', got '%v'", data.expected, res)
			}
		})
	}
}

var testDataParseTimeUntil = []struct {
	input    string
	expected string
}{
	{"2024-01-10 11:55", "2024-01-10 11:55:59"},
	{"2024-01-10", "2024-01-10 23:59:59"},
	{"11:55", "2024-01-11 11:55:59"},
	{"22:55", "2024-01-10 22:55:59"},
}

func TestParseTimeUntil(t *testing.T) {
	now = func() time.Time {
		newNow, _ := time.ParseInLocation("2006-01-02 15:04:05", "2024-01-10 17:08:23", loc)
		return newNow
	}

	for _, data := range testDataParseTimeUntil {
		t.Run(fmt.Sprintf("input=%s", data.input), func(t *testing.T) {
			// > Arrange
			expected, _ := time.ParseInLocation("2006-01-02 15:04:05", data.expected, loc)

			// > Act
			res, err := ParseTimeUntil(data.input)

			// > Assert
			if data.expected == "" && err == nil {
				t.Errorf("Parsing '%v' should have failed", data.input)
			}
			if !res.Equal(expected) {
				t.Errorf("Expected '%v', got '%v'", data.expected, res)
			}
		})
	}
}
