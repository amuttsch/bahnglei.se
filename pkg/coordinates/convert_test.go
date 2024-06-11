package coordinates_test

import (
	"math"
	"testing"

	"github.com/amuttsch/bahnglei.se/pkg/coordinates"
)

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func TestXtoLon(t *testing.T) {
	actual := coordinates.X2lon(8504, 14)

	if almostEqual(actual, 6.855469) {
		t.Errorf("Expected 6.855469, got %f", actual)
	}
}

func TestYtoLat(t *testing.T) {
	actual := coordinates.Y2lat(5473, 14)

	if almostEqual(actual, 51.1578) {
		t.Errorf("Expected 51.1578, got %f", actual)
	}
}
