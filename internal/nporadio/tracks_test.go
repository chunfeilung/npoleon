package nporadio

import (
	"github.com/google/uuid"
	"npoleon/internal/util"
	"testing"
)

func TestConvertResponse(t *testing.T) {
	response := Response{
		PageProps: PageProps{
			TrackPlays: []Play{
				{
					Id:     "12345678-90ab-cdef-1234-567890abcdef",
					Artist: "Youssou N'Dour",
					Track:  "Koeman gooit alle dingen om",
					Time:   "07:13",
				},
				{
					Id:     "abcdef12-3456-7890-abcd-ef1234567890",
					Artist: "Kaoma",
					Track:  "Waar is toch dat zebrahondje voor",
					Time:   "07:18",
				},
			},
			InitialValues: InitialValues{
				Date: "02-07-2024",
			},
		},
	}

	t.Run("Deserialize Response with multiple Play objects", func(t *testing.T) {
		// > Act
		res, _ := convertResponse(response)

		// > Assert
		if len(res) != 2 {
			t.Errorf("Expected 2 Tracks, found %d", len(res))
		}
	})

	t.Run("Deserialize single Play into Track object", func(t *testing.T) {
		// > Act
		res, _ := convertResponse(response)

		// > Assert
		expectedUuid, _ := uuid.Parse("12345678-90ab-cdef-1234-567890abcdef")
		expectedTime, _ := util.ParseTime("2024-07-02 07:13")
		expectedTrack := Track{
			Id:       expectedUuid,
			Artist:   "Youssou N'Dour",
			Title:    "Koeman gooit alle dingen om",
			PlayedAt: expectedTime.Time,
		}

		if !res[0].Equal(expectedTrack) {
			t.Errorf("Expected Track %v, found %v", expectedTrack, res[0])
		}
	})

	t.Run("Failure to deserialize an invalid Play object", func(t *testing.T) {
		// > Arrange
		response.PageProps.TrackPlays[0].Time = "42:00"

		// > Act
		_, err := convertResponse(response)

		// > Assert
		if err == nil {
			t.Errorf("Expected conversion to fail")
		}
	})
}

func TestTrack_PlayIdentifier(t *testing.T) {
	// > Arrange
	identifier, _ := uuid.Parse("9f2b7020-25b0-48d4-b895-d5f0c08857ca")
	playedAt, _ := util.ParseTime("2024-06-04 12:10:00")
	track := Track{
		Id:       identifier,
		Artist:   "Marko Schuitmakker",
		Title:    "Hengelbewaarder",
		PlayedAt: playedAt.Time,
	}

	// > Act
	res := track.PlayIdentifier()

	// > Assert
	if res != "12:10 9f2b7020-25b0-48d4-b895-d5f0c08857ca" {
		t.Errorf("Track play identifier is '%v'", res)
	}
}

func TestTrack_IsPlayedAt(t *testing.T) {
	playedAt, _ := util.ParseTime("2024-11-11 11:11:00")
	track := Track{
		Id:       uuid.New(),
		Artist:   "K4",
		Title:    "Klusjesdag",
		PlayedAt: playedAt.Time,
	}

	t.Run("One minute after track started playing", func(t *testing.T) {
		// > Arrange
		now, _ := util.ParseTime("2024-11-11 11:12:00")

		// > Act
		res := track.IsPlayedAt(now.Time)

		// > Assert
		if !res {
			t.Errorf("Track that started playing a minute ago should still be playing")
		}
	})

	t.Run("Four minutes after track started playing", func(t *testing.T) {
		// > Arrange
		now, _ := util.ParseTime("2024-11-11 11:15:00")

		// > Act
		res := track.IsPlayedAt(now.Time)

		// > Assert
		if res {
			t.Errorf("Track that started playing four minutes ago should no longer be playing")
		}
	})

	t.Run("One minute before track starts playing", func(t *testing.T) {
		// > Arrange
		now, _ := util.ParseTime("2024-11-11 11:10:00")

		// > Act
		res := track.IsPlayedAt(now.Time)

		// > Assert
		if res {
			t.Errorf("Track that starts playing in a minute should not be playing")
		}
	})
}
