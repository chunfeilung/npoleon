package http

import (
	"strings"
	"testing"
)

func TestGetBuildId(t *testing.T) {
	// > Arrange
	httpClient := &Client{}

	// > Act
	res, _ := httpClient.Fetch("https://chuniversiteit.nl/chungfeilung/")

	// > Assert
	if !strings.Contains(string(res), "HET IS CHUN!!!") {
		t.Errorf("Remote response does not contain expected text")
	}
}
