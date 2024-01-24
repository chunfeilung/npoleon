package nporadio

import (
	"errors"
	"fmt"
)

type StationId string

const (
	NpoRadio1 StationId = "nporadio1"
	NpoRadio2 StationId = "nporadio2"
	NpoRadio3 StationId = "npo3fm"
)

func GetStationId(station string) (StationId, error) {
	switch station {
	case "nporadio1", "radio1":
		return NpoRadio1, nil
	case "nporadio2", "radio2":
		return NpoRadio2, nil
	case "3fm", "npo3fm", "nporadio3", "radio3":
		return NpoRadio3, nil
	}
	err := errors.New(fmt.Sprintf("invalid station '%s'", station))
	return "", err
}
