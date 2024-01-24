package nporadio

import (
	"fmt"
	"testing"
)

var testDataGetStationId = []struct {
	input    string
	expected StationId
}{
	{"radio1", NpoRadio1},
	{"nporadio2", NpoRadio2},
	{"npo3fm", NpoRadio3},
	{"radio538", ""},
}

func TestGetStationId(t *testing.T) {
	for _, data := range testDataGetStationId {
		t.Run(fmt.Sprintf("input=%s", data.input), func(t *testing.T) {
			// > Act
			result, _ := GetStationId(data.input)

			// > Assert
			if string(result) != string(data.expected) {
				t.Errorf("Expected %v, got %v", data.expected, result)
			}
		})
	}
}
