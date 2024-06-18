package coordinates_test

import (
	"math"
	"testing"

	"github.com/amuttsch/bahnglei.se/pkg/coordinates"
)

func almostEqual(a, b, t float64) bool {
	return math.Abs(a-b) <= t
}

func TestXtoLon(t *testing.T) {
	actual := coordinates.X2lon(8504, 14)

	if !almostEqual(actual, 6.855469, 1e-6) {
		t.Errorf("Expected 6.855469, got %f", actual)
	}
}

func TestYtoLat(t *testing.T) {
	actual := coordinates.Y2lat(5473, 14)

	if !almostEqual(actual, 51.165567, 1e-6) {
		t.Errorf("Expected 51.1578, got %f", actual)
	}
}

func TestDistance(t *testing.T) {
	actual := coordinates.Distance(52.5144450, 13.3500838, 52.5162579, 13.3776784)

	if !almostEqual(actual, 1880, 1.0) {
		t.Errorf("Expected 1800m got %f", actual)
	}
}
